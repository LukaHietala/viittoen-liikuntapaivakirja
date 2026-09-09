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
			description TEXT NOT NULL,
			unit TEXT NOT NULL,
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

		INSERT INTO users VALUES(NULL, "Jaakko", "jaakko@cat-v.org", "$2a$12$EA6dIz7gLEN4Ziiy60meleFgxxRscjAkfVLkqCuCFnhfnFX/AMike", FALSE);
		INSERT INTO users VALUES(NULL, "Tero", "tero@cat-v.org", "$2a$12$EA6dIz7gLEN4Ziiy60meleFgxxRscjAkfVLkqCuCFnhfnFX/AMike", FALSE);
		INSERT INTO users VALUES(NULL, "admin", "admin@admin.com", "$2a$12$oNY0On.kHD0PFA4yg/HwaOxyAKjcrDacFaf3dKtxLzH3XULDPVCtu", TRUE);
		INSERT INTO challenges VALUES(NULL, "Käy salilla", "Kerätään sali tunteja.", "tunteja", 10, "2026-08-01", "2026-09-25");
		INSERT INTO challenges VALUES(NULL, "Lenkkihaaste", "Joka päivä kävele 1 km", "km", 25, "2026-06-01", "2026-07-25");
		INSERT INTO performances(points, user_id, challenge_id) VALUES (2, 1, 1);
		INSERT INTO performances(points, user_id, challenge_id) VALUES (1, 1, 1);
		INSERT INTO performances(points, user_id, challenge_id) VALUES (1, 2, 1);
		INSERT INTO performances(points, user_id, challenge_id) VALUES (4, 2, 2);
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}