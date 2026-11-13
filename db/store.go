package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var schema = `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		password_hash TEXT NOT NULL,
		is_admin BOOLEAN NOT NULL DEFAULT FALSE,
		color TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS collections (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		goal_points INTEGER NOT NULL,
		start_date TEXT NOT NULL,
		end_date TEXT NOT NULL,
		user_id	INTEGER NOT NULL,
		created_at TEXT NOT NULL DEFAULT (datetime(CURRENT_TIMESTAMP, 'localtime'))
	);

	CREATE TABLE IF NOT EXISTS challenges (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		unit TEXT NOT NULL,
		goal_points INTEGER NOT NULL,
		start_date TEXT NOT NULL,
		end_date TEXT NOT NULL,
		user_id	INTEGER NOT NULL,
		collection_id INTEGER NOT NULL,
		created_at TEXT NOT NULL DEFAULT (datetime(CURRENT_TIMESTAMP, 'localtime')),
		FOREIGN KEY(collection_id) REFERENCES collections(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS performances (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		points INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		challenge_id INTEGER NOT NULL,
		created_at TEXT NOT NULL DEFAULT (datetime(CURRENT_TIMESTAMP, 'localtime')),
		FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY(challenge_id) REFERENCES challenges(id) ON DELETE CASCADE
	);
`

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func Connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	if _, err = db.Exec(schema); err != nil {
		return nil, err
	}

	return db, nil
}
