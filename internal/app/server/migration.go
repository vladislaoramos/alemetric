package server

import (
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

func runMigration(dbURL, migrationDir string) error {
	db, err := goose.OpenDBWithDriver("postgres", dbURL)
	if err != nil {
		return err
	}

	defer db.Close()

	if err := goose.Up(db, migrationDir); err != nil {
		return err
	}

	return nil
}
