package projects

import "context"

// ProjectStore defines the interface for project data operations
type ProjectStore interface {
	// Create inserts a new project with the given name and returns the created project
	Create(ctx context.Context, name string) (*Project, error)

	// CreateMultiple inserts multiple projects and returns all created projects
	// Returns a slice of created projects or an error if any insertion fails
	CreateMultiple(ctx context.Context, names []string) ([]Project, error)

	// GetByName retrieves a project by its name
	GetByName(ctx context.Context, name string) (*Project, error)

	// GetAll retrieves all projects ordered by creation date (newest first)
	// If namePrefix is not empty, filters projects by name prefix (case-insensitive)
	GetAll(ctx context.Context, namePrefix string) ([]Project, error)

	DeleteByName(ctx context.Context, name string) error
}
