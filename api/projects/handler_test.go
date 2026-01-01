package projects

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// mockStore is a mock implementation of ProjectStore for testing
type mockStore struct {
	createFunc         func(ctx context.Context, name string) (*Project, error)
	createMultipleFunc func(ctx context.Context, names []string) ([]Project, error)
	getByIDFunc        func(ctx context.Context, id int) (*Project, error)
	getAllFunc         func(ctx context.Context, namePrefix string) ([]Project, error)
}

func (m *mockStore) Create(ctx context.Context, name string) (*Project, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, name)
	}
	return nil, nil
}

func (m *mockStore) CreateMultiple(ctx context.Context, names []string) ([]Project, error) {
	if m.createMultipleFunc != nil {
		return m.createMultipleFunc(ctx, names)
	}
	return nil, nil
}

func (m *mockStore) GetByID(ctx context.Context, id int) (*Project, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockStore) GetAll(ctx context.Context, namePrefix string) ([]Project, error) {
	if m.getAllFunc != nil {
		return m.getAllFunc(ctx, namePrefix)
	}
	return nil, nil
}

// TestCreateProject_Success tests successful project creation
func TestCreateProject_Success(t *testing.T) {
	mock := &mockStore{
		createFunc: func(ctx context.Context, name string) (*Project, error) {
			return &Project{
				ID:        1,
				Name:      name,
				CreatedAt: time.Now(),
			}, nil
		},
	}

	handler := NewHandler(mock)
	req := httptest.NewRequest("POST", "/projects", strings.NewReader(`{"name":"Test Project"}`))
	w := httptest.NewRecorder()

	handler.CreateProject(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var result Project
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result.ID != 1 || result.Name != "Test Project" {
		t.Errorf("unexpected response body: %+v", result)
	}
}

// TestCreateProject_InvalidBody tests invalid request body
func TestCreateProject_InvalidBody(t *testing.T) {
	mock := &mockStore{}
	handler := NewHandler(mock)

	req := httptest.NewRequest("POST", "/projects", strings.NewReader(`{invalid json}`))
	w := httptest.NewRecorder()

	handler.CreateProject(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestCreateProject_NameConflict tests duplicate project name error
func TestCreateProject_NameConflict(t *testing.T) {
	mock := &mockStore{
		createFunc: func(ctx context.Context, name string) (*Project, error) {
			return nil, ErrProjectNameConflict
		},
	}

	handler := NewHandler(mock)
	req := httptest.NewRequest("POST", "/projects", strings.NewReader(`{"name":"Duplicate"}`))
	w := httptest.NewRecorder()

	handler.CreateProject(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("expected status %d, got %d", http.StatusConflict, w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "Project name already exists") {
		t.Errorf("expected error message about duplicate, got: %s", body)
	}
}

// TestCreateProject_DatabaseError tests database error handling
func TestCreateProject_DatabaseError(t *testing.T) {
	mock := &mockStore{
		createFunc: func(ctx context.Context, name string) (*Project, error) {
			return nil, errors.New("database error")
		},
	}

	handler := NewHandler(mock)
	req := httptest.NewRequest("POST", "/projects", strings.NewReader(`{"name":"Test"}`))
	w := httptest.NewRecorder()

	handler.CreateProject(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// TestGetProject_Success tests successful project retrieval
func TestGetProject_Success(t *testing.T) {
	mock := &mockStore{
		getByIDFunc: func(ctx context.Context, id int) (*Project, error) {
			return &Project{
				ID:        id,
				Name:      "Test Project",
				CreatedAt: time.Now(),
			}, nil
		},
	}

	handler := NewHandler(mock)
	req := httptest.NewRequest("GET", "/projects/1", nil)
	w := httptest.NewRecorder()

	// Simulate path value extraction
	req = req.WithContext(context.WithValue(req.Context(), "id", "1"))
	req.SetPathValue("id", "1")

	handler.GetProject(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result Project
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result.ID != 1 || result.Name != "Test Project" {
		t.Errorf("unexpected response body: %+v", result)
	}
}

// TestGetProject_NotFound tests project not found error
func TestGetProject_NotFound(t *testing.T) {
	mock := &mockStore{
		getByIDFunc: func(ctx context.Context, id int) (*Project, error) {
			return nil, errors.New("project not found")
		},
	}

	handler := NewHandler(mock)
	req := httptest.NewRequest("GET", "/projects/999", nil)
	req.SetPathValue("id", "999")
	w := httptest.NewRecorder()

	handler.GetProject(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

// TestGetProject_InvalidID tests invalid project ID
func TestGetProject_InvalidID(t *testing.T) {
	mock := &mockStore{}
	handler := NewHandler(mock)

	req := httptest.NewRequest("GET", "/projects/invalid", nil)
	req.SetPathValue("id", "invalid")
	w := httptest.NewRecorder()

	handler.GetProject(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestListProjects_Success tests successful listing of all projects
func TestListProjects_Success(t *testing.T) {
	projects := []Project{
		{ID: 1, Name: "Project 1", CreatedAt: time.Now()},
		{ID: 2, Name: "Project 2", CreatedAt: time.Now()},
	}

	mock := &mockStore{
		getAllFunc: func(ctx context.Context, namePrefix string) ([]Project, error) {
			return projects, nil
		},
	}

	handler := NewHandler(mock)
	req := httptest.NewRequest("GET", "/projects", nil)
	w := httptest.NewRecorder()

	handler.ListProjects(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result []Project
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 projects, got %d", len(result))
	}
}

// TestListProjects_Empty tests listing when no projects exist
func TestListProjects_Empty(t *testing.T) {
	mock := &mockStore{
		getAllFunc: func(ctx context.Context, namePrefix string) ([]Project, error) {
			return []Project{}, nil
		},
	}

	handler := NewHandler(mock)
	req := httptest.NewRequest("GET", "/projects", nil)
	w := httptest.NewRecorder()

	handler.ListProjects(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result []Project
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("expected 0 projects, got %d", len(result))
	}
}

// TestListProjects_DatabaseError tests database error handling
func TestListProjects_DatabaseError(t *testing.T) {
	mock := &mockStore{
		getAllFunc: func(ctx context.Context, namePrefix string) ([]Project, error) {
			return nil, errors.New("database error")
		},
	}

	handler := NewHandler(mock)
	req := httptest.NewRequest("GET", "/projects", nil)
	w := httptest.NewRecorder()

	handler.ListProjects(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

// TestCreateMultipleProjects_Success tests successful batch creation
func TestCreateMultipleProjects_Success(t *testing.T) {
	createdProjects := []Project{
		{ID: 1, Name: "Project 1", CreatedAt: time.Now()},
		{ID: 2, Name: "Project 2", CreatedAt: time.Now()},
		{ID: 3, Name: "Project 3", CreatedAt: time.Now()},
	}

	mock := &mockStore{
		createMultipleFunc: func(ctx context.Context, names []string) ([]Project, error) {
			return createdProjects, nil
		},
	}

	handler := NewHandler(mock)
	body := `{"projects":["Project 1","Project 2","Project 3"]}`
	req := httptest.NewRequest("POST", "/projects/batch", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateMultipleProjects(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var result []Project
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("expected 3 projects, got %d", len(result))
	}
}

// TestCreateMultipleProjects_EmptyList tests empty projects list
func TestCreateMultipleProjects_EmptyList(t *testing.T) {
	mock := &mockStore{}
	handler := NewHandler(mock)

	body := `{"projects":[]}`
	req := httptest.NewRequest("POST", "/projects/batch", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateMultipleProjects(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestCreateMultipleProjects_InvalidBody tests invalid request body
func TestCreateMultipleProjects_InvalidBody(t *testing.T) {
	mock := &mockStore{}
	handler := NewHandler(mock)

	req := httptest.NewRequest("POST", "/projects/batch", strings.NewReader(`{invalid}`))
	w := httptest.NewRecorder()

	handler.CreateMultipleProjects(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// TestCreateMultipleProjects_DuplicateName tests duplicate name handling in batch
func TestCreateMultipleProjects_DuplicateName(t *testing.T) {
	mock := &mockStore{
		createMultipleFunc: func(ctx context.Context, names []string) ([]Project, error) {
			return nil, errors.New("duplicate project name: Project 1")
		},
	}

	handler := NewHandler(mock)
	body := `{"projects":["Project 1","Project 1"]}`
	req := httptest.NewRequest("POST", "/projects/batch", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateMultipleProjects(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("expected status %d, got %d", http.StatusConflict, w.Code)
	}

	bodyStr := w.Body.String()
	if !strings.Contains(bodyStr, "duplicate project name") {
		t.Errorf("expected error message about duplicate, got: %s", bodyStr)
	}
}

// TestCreateMultipleProjects_DatabaseError tests database error handling in batch
func TestCreateMultipleProjects_DatabaseError(t *testing.T) {
	mock := &mockStore{
		createMultipleFunc: func(ctx context.Context, names []string) ([]Project, error) {
			return nil, errors.New("database error")
		},
	}

	handler := NewHandler(mock)
	body := `{"projects":["Project 1","Project 2"]}`
	req := httptest.NewRequest("POST", "/projects/batch", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateMultipleProjects(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
