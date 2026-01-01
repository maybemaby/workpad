package projects

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/maybemaby/workpad/api/utils"
)

// ProjectHandler handles HTTP requests for projects
type ProjectHandler struct {
	store ProjectStore
}

// NewHandler creates a new projects handler
func NewHandler(store ProjectStore) *ProjectHandler {
	return &ProjectHandler{store: store}
}

type GetProjectRequest struct {
	ID int `path:"id" example:"1"`
}

type ListProjectsRequest struct {
	Prefix string `query:"prefix" example:"Proj" required:"false"`
}

// CreateProject handles POST /projects
// Returns the created project or an existing project with the same name
func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req CreateProjectRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	project, err := h.store.Create(r.Context(), req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(project)
}

// GetProject handles GET /projects/{id}
func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	project, err := h.store.GetByID(r.Context(), id)
	if err != nil {
		if err.Error() == "project not found" {
			http.Error(w, "Project not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(project)
}

// ListProjects handles GET /projects
func (h *ProjectHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	namePrefix := r.URL.Query().Get("prefix")

	projects, err := h.store.GetAll(r.Context(), namePrefix)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, r, projects)
}

// CreateMultipleProjects handles POST /projects/batch
// Returns all projects after upserting the provided names
func (h *ProjectHandler) CreateMultipleProjects(w http.ResponseWriter, r *http.Request) {
	var req CreateMultipleProjectsRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Projects) == 0 {
		http.Error(w, "Projects list cannot be empty", http.StatusBadRequest)
		return
	}

	projects, err := h.store.CreateMultiple(r.Context(), req.Projects)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(projects)
}
