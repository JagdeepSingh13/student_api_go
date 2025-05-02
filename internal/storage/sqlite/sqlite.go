package sqlite

import (
	"database/sql"

	"github.com/JagdeepSingh13/student_api_go/internal/config"
	// imported but not used
	// _ "github.com/mattn/go-sqlite3"

	_ "modernc.org/sqlite"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	// did not work with sqlite3
	// db, err := sql.Open("sqlite3", cfg.StoragePath)
	db, err := sql.Open("sqlite", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT,
	age INTEGER
	)`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil
}
