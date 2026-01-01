package api

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"

	"github.com/jmoiron/sqlx"
)

func NewSqliteDB(ctx context.Context, trace bool) (*sqlx.DB, error) {
	dbPath := os.Getenv("SQLITE_DB_PATH")
	if dbPath == "" {
		return nil, fmt.Errorf("SQLITE_DB_PATH environment variable not set")
	}

	connStr := fmt.Sprintf("file:%s?cache=shared&mode=rwc", dbPath)

	// Open connection using stdlib sql driver directly
	sqlDB, err := sql.Open("sqlite", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	// Create sqlx database wrapper around stdlib connection
	db := sqlx.NewDb(sqlDB, "sqlite")

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	// Verify connection
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping sqlite database: %w", err)
	}

	// Enable foreign keys
	if _, err := db.ExecContext(ctx, "PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	return db, nil
}
