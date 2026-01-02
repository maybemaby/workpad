package projects

import "time"

type Project struct {
	Name      string    `json:"name" required:"true"`
	CreatedAt time.Time `json:"created_at" db:"created_at" required:"true"`
}

type CreateProjectRequest struct {
	Name string `json:"name" required:"true"`
}

// CreateMultipleProjectsRequest is the request body for batch project creation
type CreateMultipleProjectsRequest struct {
	Projects []string `json:"projects" example:"[Project A, Project B]" required:"true"`
}
