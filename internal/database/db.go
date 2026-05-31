package database

import (
	"database/sql"
	_ "embed"

	_ "modernc.org/sqlite"
)

//go:embed migrations/001_initial.sql
var initialSchema string

func New(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if _, err := db.Exec(initialSchema); err != nil {
		return nil, err
	}

	return db, nil
}
