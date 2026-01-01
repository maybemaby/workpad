package projects

import "time"

type Project struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CreateProjectRequest struct {
	Name string `json:"name"`
}

// CreateMultipleProjectsRequest is the request body for batch project creation
type CreateMultipleProjectsRequest struct {
	Projects []string `json:"projects" example:"[Project A, Project B]"`
}
