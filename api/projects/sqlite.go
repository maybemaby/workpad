package projects

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/jmoiron/sqlx"
)

// ErrProjectNameConflict is returned when a project with the same name already exists
var ErrProjectNameConflict = errors.New("project name already exists")

// SqliteStore implements the Store interface using SQLite
type SqliteStore struct {
	db *sqlx.DB
}

// NewSqliteStore creates a new SQLite store
func NewSqliteStore(db *sqlx.DB) *SqliteStore {
	return &SqliteStore{db: db}
}

// Create inserts a new project into the database
func (s *SqliteStore) Create(ctx context.Context, name string) (*Project, error) {
	if name == "" {
		return nil, fmt.Errorf("project name cannot be empty")
	}

	cleanedName := strings.TrimSpace(name)

	query := `INSERT INTO projects (name) VALUES (?) RETURNING id, name, created_at`

	var project Project
	err := s.db.QueryRowContext(ctx, query, cleanedName).Scan(&project.ID, &project.Name, &project.CreatedAt)
	if err != nil {
		// Check for UNIQUE constraint violation
		if isSQLiteConstraintError(err) {
			return nil, ErrProjectNameConflict
		}
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	return &project, nil
}

// CreateMultiple inserts multiple projects into the database
// Uses a transaction to ensure all-or-nothing semantics
func (s *SqliteStore) CreateMultiple(ctx context.Context, names []string) ([]Project, error) {
	if len(names) == 0 {
		return []Project{}, nil
	}

	// Validate all names before inserting
	if slices.Contains(names, "") {
		return nil, fmt.Errorf("project name cannot be empty")
	}

	cleanedNames := make([]string, len(names))
	
	for i, name := range names {
		cleanedNames[i] = strings.TrimSpace(name)
	}

	// Start a transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	var projects []Project

	// Insert each project
	for _, name := range cleanedNames {
		query := `INSERT INTO projects (name) VALUES (?) RETURNING id, name, created_at`

		var project Project
		err := tx.QueryRowContext(ctx, query, name).Scan(&project.ID, &project.Name, &project.CreatedAt)
		if err != nil {
			// Check for UNIQUE constraint violation
			if isSQLiteConstraintError(err) {
				return nil, fmt.Errorf("duplicate project name: %s", name)
			}
			return nil, fmt.Errorf("failed to create project: %w", err)
		}

		projects = append(projects, project)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return projects, nil
}

// GetByID retrieves a project by its ID
func (s *SqliteStore) GetByID(ctx context.Context, id int) (*Project, error) {
	query := `SELECT id, name, created_at FROM projects WHERE id = ?`

	var project Project
	err := s.db.QueryRowContext(ctx, query, id).Scan(&project.ID, &project.Name, &project.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("project not found")
		}
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return &project, nil
}

// GetAll retrieves all projects ordered by creation date (newest first)
// If namePrefix is not empty, filters projects by name prefix (case-insensitive)
func (s *SqliteStore) GetAll(ctx context.Context, namePrefix string) ([]Project, error) {
	var query string
	var args []any

	if namePrefix != "" {
		// Filter by name prefix (case-insensitive)
		query = `SELECT id, name, created_at FROM projects WHERE LOWER(name) LIKE LOWER(?) ORDER BY created_at DESC`
		args = []any{namePrefix + "%"}
	} else {
		// Get all projects
		query = `SELECT id, name, created_at FROM projects ORDER BY created_at DESC`
	}

	var projects []Project
	err := s.db.SelectContext(ctx, &projects, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}

	// Return empty slice instead of nil for consistency
	if projects == nil {
		projects = make([]Project, 0)
	}

	return projects, nil
}

// isSQLiteConstraintError checks if the error is a SQLite constraint violation
// SQLite returns error messages like "UNIQUE constraint failed: projects.name"
func isSQLiteConstraintError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	// Check for UNIQUE constraint violation in error message
	return strings.Contains(errMsg, "UNIQUE constraint failed") ||
		strings.Contains(errMsg, "unique constraint failed")
}
