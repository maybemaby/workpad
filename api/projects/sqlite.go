package projects

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strings"

	"github.com/jmoiron/sqlx"
)

// SqliteStore implements the Store interface using SQLite
type SqliteStore struct {
	db *sqlx.DB
}

// NewSqliteStore creates a new SQLite store
func NewSqliteStore(db *sqlx.DB) *SqliteStore {
	return &SqliteStore{db: db}
}

// Create inserts a new project or returns the existing one if name already exists
// Uses SQLite UPSERT syntax: INSERT ... ON CONFLICT ... DO NOTHING
// This is atomic and returns the project (new or existing)
func (s *SqliteStore) Create(ctx context.Context, name string) (*Project, error) {
	if name == "" {
		return nil, fmt.Errorf("project name cannot be empty")
	}

	cleanedName := strings.TrimSpace(name)

	// SQLite upsert: insert if not exists, do nothing if conflict on unique constraint
	// Then retrieve the (existing or newly created) project
	query := `INSERT INTO projects (name) VALUES (?) ON CONFLICT(name) DO NOTHING`

	_, err := s.db.ExecContext(ctx, query, cleanedName)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	// Retrieve the project (existing or newly created) by name
	getQuery := `SELECT id, name, created_at FROM projects WHERE name = ?`
	var project Project
	err = s.db.QueryRowContext(ctx, getQuery, cleanedName).Scan(&project.ID, &project.Name, &project.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve project: %w", err)
	}

	return &project, nil
}

// CreateMultiple inserts multiple projects using SQLite upsert syntax
// Uses a transaction to ensure atomicity
// Projects with duplicate names will be ignored (no change if already exist)
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

	// Insert each project using SQLite upsert syntax
	for _, name := range cleanedNames {
		query := `INSERT INTO projects (name) VALUES (?) ON CONFLICT(name) DO NOTHING`

		_, err := tx.ExecContext(ctx, query, name)
		if err != nil {
			return nil, fmt.Errorf("failed to create project: %w", err)
		}

		// Retrieve the project (existing or newly created) by name
		getQuery := `SELECT id, name, created_at FROM projects WHERE name = ?`
		var project Project
		err = tx.QueryRowContext(ctx, getQuery, name).Scan(&project.ID, &project.Name, &project.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve project: %w", err)
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
