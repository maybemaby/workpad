package utils

import (
	"context"
	"database/sql"

	"github.com/maybemaby/workpad/migrations"
)

func SetupSqliteDb(db *sql.DB) error {

	err := migrations.RunMigrations(context.Background(), db)

	if err != nil {
		return err
	}

	return nil
}
