package api

import (
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/maybemaby/workpad/api/notes"
	"github.com/maybemaby/workpad/api/projects"
	"github.com/maybemaby/workpad/frontend"
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

	rootMw := RootMiddleware(s.logger, MiddlewareConfig{
		CorsOrigin: "http://localhost:5173",
	})

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

	apiRoute := r.Group("/api")

	// Projects routes
	projectsStore := projects.NewSqliteStore(s.sqliteDB)
	projectsHandler := projects.NewHandler(projectsStore)

	apiRoute.Handle("POST /projects", rootMw.ThenFunc(projectsHandler.CreateProject)).With(
		option.Request(new(projects.CreateProjectRequest)),
		option.Response(201, new(projects.Project)),
		option.Tags("Projects"),
	)

	apiRoute.Handle("GET /projects", rootMw.ThenFunc(projectsHandler.ListProjects)).With(
		option.Request(new(projects.ListProjectsRequest)),
		option.Response(200, new([]projects.Project)),
		option.Tags("Projects"),
	)

	apiRoute.Handle("GET /projects/{name}", rootMw.ThenFunc(projectsHandler.GetProject)).With(
		option.Request(new(projects.GetProjectRequest)),
		option.Response(200, new(projects.Project)),
		option.Response(404, "Not Found"),
		option.Tags("Projects"),
	)

	apiRoute.Handle("POST /projects/batch", rootMw.ThenFunc(projectsHandler.CreateMultipleProjects)).With(
		option.Request(new(projects.CreateMultipleProjectsRequest)),
		option.Response(201, new([]projects.Project)),
		option.Tags("Projects"),
	)

	apiRoute.Handle("DELETE /projects/{name}", rootMw.ThenFunc(projectsHandler.DeleteProject)).With(
		option.Request(new(projects.GetProjectRequest)),
		option.Response(204, nil),
		option.Tags("Projects"),
	)

	// Notes routes
	noteStore := notes.NewNoteService(s.sqliteDB)
	notesHandler := notes.NewNoteHandler(noteStore)

	apiRoute.Handle("GET /notes/by-date", rootMw.ThenFunc(notesHandler.GetNoteByDate)).With(
		option.Request(new(notes.GetNoteByDateRequest)),
		option.Response(200, new(notes.Note)),
		option.Response(404, "Not Found"),
		option.Tags("Notes"),
	)

	apiRoute.Handle("POST /notes", rootMw.ThenFunc(notesHandler.CreateNote)).With(
		option.Request(new(notes.CreateNoteRequest)),
		option.Response(201, new(notes.Note)),
		option.Tags("Notes"),
	)

	apiRoute.Handle("GET /notes/for-month", rootMw.ThenFunc(notesHandler.GetMonthNotes)).With(
		option.Request(new(notes.GetMonthNotesRequest)),
		option.Response(200, new([]int)),
		option.Tags("Notes"),
	)

	apiRoute.Handle("PUT /notes/excerpts", rootMw.ThenFunc(notesHandler.UpdateNoteExcerpts)).With(
		option.Request(new(notes.UpdateNoteExcerptRequest)),
		option.Response(204, nil),
		option.Tags("Notes"),
	)

	apiRoute.Handle("GET /notes/excerpts/{project}", rootMw.ThenFunc(notesHandler.GetExcerptsForProject)).With(
		option.Request(new(notes.GetExcerptsForProjectRequest)),
		option.Response(200, new([]notes.NoteExcerpt)),
		option.Tags("Notes"),
	)

	apiRoute.Handle("/", rootMw.ThenFunc(
		func(w http.ResponseWriter, r *http.Request) {
			slog.Default().Info("Handling CORS preflight")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			http.NotFound(w, r)
		},
	))

	// For CORS preflight requests
	r.Handle("/", rootMw.ThenFunc(
		func(w http.ResponseWriter, r *http.Request) {

			if r.Method == http.MethodGet {
				HandleSPA(frontend.Assets).ServeHTTP(w, r)
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
		slog.Default().Error("Error accessing file", "error", err)
		http.Error(w, "Error accessing file", http.StatusInternalServerError)
	})
}

func HandleSPA(filesys fs.FS) http.Handler {

	sub, err := fs.Sub(filesys, "build")

	if err != nil {
		panic(err)
	}

	fileServer := http.FileServer(http.FS(sub))
	prefix := "build"

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		immutablePath := strings.Contains(path, "/_app/immutable")

		file, err := filesys.Open(prefix + path)

		if err == nil {
			file.Close()

			if immutablePath {
				w.Header().Add("Cache-Control", "max-age=3600")
			}

			fileServer.ServeHTTP(w, r)
			return
		}

		if os.IsNotExist(err) {
			indexData, err := fs.ReadFile(filesys, prefix+"/index.html")
			if err != nil {
				http.Error(w, "Index file not found", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write(indexData)
			return
		}

		http.Error(w, "Error accessing file", http.StatusInternalServerError)
	})
}
