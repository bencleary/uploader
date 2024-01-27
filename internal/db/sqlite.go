package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	db *sql.DB
}

// NewSQLiteDatabase creates a new instance of DB using an SQLite database file.
//
// dbPath: the path to the SQLite database file.
// Returns a pointer to the DB instance and an error if any.
func NewSQLiteDatabase(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	return &DB{db: db}, nil
}

func (d *DB) Close() error {
	return d.db.Close()
}

func (d *DB) CreateTable() error {
	_, err := d.db.Exec(`
		CREATE TABLE IF NOT EXISTS uploads (
			id INTEGER PRIMARY KEY,
			uuid TEXT,
			owner_id INTEGER,
			file_name TEXT,
			file_size INTEGER,
			extension TEXT,
			mime_type TEXT
		)
	`)
	return err
}
