package api

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"

	"github.com/maybemaby/workpad/api/auth"
	"github.com/maybemaby/workpad/api/notes"
	"github.com/maybemaby/workpad/api/projects"
	"github.com/oaswrap/spec-ui/config"
	"github.com/oaswrap/spec/adapter/httpopenapi"
	"github.com/oaswrap/spec/option"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func httpSpanName(operation string, r *http.Request) string {
	return fmt.Sprintf("HTTP %s %s", r.Method, r.URL.Path)
}

func (s *Server) MountRoutesOapi() {

	mux := http.NewServeMux()

	authHandler := &AuthHandler{
		jwtManager: s.jwtManager,
		pool:       s.pool,
	}

	rootMw := RootMiddleware(s.logger, MiddlewareConfig{
		CorsOrigin: "http://localhost:5173",
	})

	authMw := rootMw.Append(auth.RequireAccessToken(s.jwtManager))

	r := httpopenapi.NewGenerator(mux,
		option.WithTitle("workpad"),
		option.WithVersion("0.1.0"),
		option.WithSecurity("bearerAuth", option.SecurityHTTPBearer("Bearer")),
		option.WithSwaggerUI(config.SwaggerUI{
			UIConfig: map[string]string{
				"persistAuthorization": "true",
			},
		}),
		option.WithExternalDocs("/docs/openapi.json"),
		option.WithDisableDocs(s.prod),
	)

	if !s.prod {

		s.logger.Info("Swagger docs enabled at /docs")

		r.HandleFunc("/docs/openapi.json", func(w http.ResponseWriter, req *http.Request) {
			spec, err := r.MarshalJSON()

			if err != nil {
				s.logger.Error("Error generating OpenAPI spec", "error", err)
				http.Error(w, "Error generating OpenAPI spec", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(spec)
		})
	}

	authRoute := r.Group("/auth").With(option.GroupTags("auth"))

	authRoute.Handle("GET /me", authMw.ThenFunc(authHandler.GetAuthMe)).With(
		option.Response(200, new(MeResponse)),
		option.Response(401, "Unauthorized"),
	)

	authRoute.Handle("POST /signup", rootMw.ThenFunc(authHandler.SignupJWT)).With(
		option.Request(new(PassSignupBody)),
		ResponsesWithDefault(map[int]any{
			200: new(LoginJwtResponse),
		}),
	)

	authRoute.Handle("POST /login", rootMw.ThenFunc(authHandler.LoginJWT)).With(
		option.Request(new(PassLoginBody)),
		Responses(map[int]any{
			401: new(AuthErrorResponse),
			200: new(LoginJwtResponse),
		}),
	)

	// Projects routes
	projectsStore := projects.NewSqliteStore(s.sqliteDB)
	projectsHandler := projects.NewHandler(projectsStore)

	r.Handle("POST /projects", rootMw.ThenFunc(projectsHandler.CreateProject)).With(
		option.Request(new(projects.CreateProjectRequest)),
		option.Response(201, new(projects.Project)),
		option.Tags("Projects"),
	)

	r.Handle("GET /projects", rootMw.ThenFunc(projectsHandler.ListProjects)).With(
		option.Request(new(projects.ListProjectsRequest)),
		option.Response(200, new([]projects.Project)),
		option.Tags("Projects"),
	)

	r.Handle("GET /projects/{name}", rootMw.ThenFunc(projectsHandler.GetProject)).With(
		option.Request(new(projects.GetProjectRequest)),
		option.Response(200, new(projects.Project)),
		option.Response(404, "Not Found"),
		option.Tags("Projects"),
	)

	r.Handle("POST /projects/batch", rootMw.ThenFunc(projectsHandler.CreateMultipleProjects)).With(
		option.Request(new(projects.CreateMultipleProjectsRequest)),
		option.Response(201, new([]projects.Project)),
		option.Tags("Projects"),
	)

	r.Handle("DELETE /projects/{name}", rootMw.ThenFunc(projectsHandler.DeleteProject)).With(
		option.Request(new(projects.GetProjectRequest)),
		option.Response(204, nil),
		option.Tags("Projects"),
	)

	// Notes routes
	noteStore := notes.NewNoteService(s.sqliteDB)
	notesHandler := notes.NewNoteHandler(noteStore)

	r.Handle("GET /notes/by-date", rootMw.ThenFunc(notesHandler.GetNoteByDate)).With(
		option.Request(new(notes.GetNoteByDateRequest)),
		option.Response(200, new(notes.Note)),
		option.Response(404, "Not Found"),
		option.Tags("Notes"),
	)

	r.Handle("POST /notes", rootMw.ThenFunc(notesHandler.CreateNote)).With(
		option.Request(new(notes.CreateNoteRequest)),
		option.Response(201, new(notes.Note)),
		option.Tags("Notes"),
	)

	r.Handle("GET /notes/for-month", rootMw.ThenFunc(notesHandler.GetMonthNotes)).With(
		option.Request(new(notes.GetMonthNotesRequest)),
		option.Response(200, new([]int)),
		option.Tags("Notes"),
	)

	r.Handle("PUT /notes/excerpts", rootMw.ThenFunc(notesHandler.UpdateNoteExcerpts)).With(
		option.Request(new(notes.UpdateNoteExcerptRequest)),
		option.Response(204, nil),
		option.Tags("Notes"),
	)

	r.Handle("GET /notes/excerpts/{project}", rootMw.ThenFunc(notesHandler.GetExcerptsForProject)).With(
		option.Request(new(notes.GetExcerptsForProjectRequest)),
		option.Response(200, new([]notes.NoteExcerpt)),
		option.Tags("Notes"),
	)

	// For CORS preflight requests
	r.Handle("/", rootMw.ThenFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			http.NotFound(w, r)
		},
	))

	srv := &http.Server{
		Addr:    ":" + s.port,
		Handler: otelhttp.NewHandler(mux, "server", otelhttp.WithSpanNameFormatter(httpSpanName)),
	}

	s.srv = srv
}

func MountSpa(mux *http.ServeMux, pattern string, filesys fs.FS) {
	fileServer := http.FileServer(http.FS(filesys))

	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Try to open the requested file
		file, err := filesys.Open(path)
		if err == nil {
			// File exists, close it and let the file server handle it
			file.Close()
			fileServer.ServeHTTP(w, r)
			return
		}

		// File doesn't exist, check if it's a not found error
		if os.IsNotExist(err) {
			// Serve index.html as fallback for SPA routing
			indexData, err := fs.ReadFile(filesys, "index.html")
			if err != nil {
				http.Error(w, "Index file not found", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write(indexData)
			return
		}

		// Other error occurred
		http.Error(w, "Error accessing file", http.StatusInternalServerError)
	})
}
