package db

import "database/sql"

var DB *sql.DB

func Connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:?_journal_mode=WAL")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL,
			password_hash TEXT NOT NULL,
			is_admin BOOLEAN NOT NULL DEFAULT FALSE
		);

		CREATE TABLE IF NOT EXISTS challenges (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			goal_points INTEGER NOT NULL,
			start_date TEXT NOT NULL,
			end_date TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS performances (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			points INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			challenge_id INTEGER NOT NULL,
			created_at TEXT NOT NULL DEFAULT (datetime(CURRENT_TIMESTAMP, 'localtime')),
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(challenge_id) REFERENCES challenges(id)
		);

		INSERT INTO users VALUES(NULL, "Jaakko", "jaakko@cat-v.org", "$2a$14$dhSvJi8wLpc0iAB5LW91Le4GKK/w9i7IKyZ6tgE7L8xnW4b2S2/lG", TRUE);
		INSERT INTO users VALUES(NULL, "Tero", "tero@cat-v.org", "$2a$14$dhSvJi8wLpc0iAB5LW91Le4GKK/w9i7IKyZ6tgE7L8xnW4b2S2/lG", TRUE);
		INSERT INTO users VALUES(NULL, "Luka", "luka.hietala08@gmail.com", "$2a$14$dhSvJi8wLpc0iAB5LW91Le4GKK/w9i7IKyZ6tgE7L8xnW4b2S2/lG", TRUE);
		INSERT INTO challenges VALUES(NULL, "Syö paljon leipää", 3, "2026-09-01", "2026-09-10");
		INSERT INTO challenges VALUES(NULL, "Käy suihkussa", 1, "2026-08-01", "2026-09-25");
		INSERT INTO challenges VALUES(NULL, "Sammuta Lukan koti", 15, "2026-06-01", "2026-07-25");
		INSERT INTO performances(points, user_id, challenge_id) VALUES (2, 1, 1);
		INSERT INTO performances(points, user_id, challenge_id) VALUES (5, 2, 3);
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}
