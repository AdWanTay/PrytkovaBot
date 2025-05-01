package storage

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

func InitDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	query := `
	CREATE TABLE IF NOT EXISTS slots (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		time DATETIME NOT NULL,
		is_booked BOOLEAN NOT NULL DEFAULT 0,
		user_id INTEGER,
		user_name TEXT
	);
	`
	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}
	Db = db
	return db, nil
}
