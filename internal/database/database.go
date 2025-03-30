package database

import (
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func New(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	schema := `
	CREATE TABLE IF NOT EXISTS jokes (
		id INTEGER PRIMARY KEY,
		joke type NOT NULL,
		nsfw bool NOT NULL,
		createdat date NOT NULL
	);`

	_, err = db.Exec(schema)
	if err != nil {
		panic(err.Error())
	}

	return db, nil
}
