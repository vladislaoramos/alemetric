package server

import (
	"fmt"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

const migrationDir = "migrations"

func applyMigration(dbURL, _ string) error {
	db, err := goose.OpenDBWithDriver("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("error open db: %w", err)
	}

	defer db.Close()

	if err := goose.Up(db, migrationDir); err != nil {
		return fmt.Errorf("error up migration: %w", err)
	}

	return nil
}
