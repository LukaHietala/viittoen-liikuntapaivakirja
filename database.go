package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Challenge struct {
	ID         int
	Title      string
	GoalPoints int
	IsActive   bool
}

type User struct {
	ID           int
	Name         string
	PasswordHash string
	IsAdmin      bool
}

type Performance struct {
	ID     int
	Points int
}

func connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			password_hash TEXT NOT NULL,
			is_admin BOOLEAN NOT NULL DEFAULT FALSE
		);

		CREATE TABLE IF NOT EXISTS challenges (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			goal_points INTEGER NOT NULL,
			is_active BOOLEAN NOT NULL DEFAULT FALSE
		);

		CREATE TABLE IF NOT EXISTS performances (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			points INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			challenge_id INTEGER NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(challenge_id) REFERENCES challenges(id)
		);

		INSERT INTO users VALUES(NULL, "Jaakko", "1234", TRUE);
		INSERT INTO challenges VALUES(NULL, "Syö paljon leipää", 3, TRUE);
		INSERT INTO challenges VALUES(NULL, "Käy suihkussa", 1, FALSE);
		INSERT INTO performances VALUES(NULL, 2, 1, 1);
		INSERT INTO performances VALUES(NULL, 1, 1, 1);
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func getAllChallenges() ([]*Challenge, error) {
	rows, err := db.Query(`SELECT * FROM challenges`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	challenges := make([]*Challenge, 0)
	for rows.Next() {
		c := new(Challenge)
		err := rows.Scan(&c.ID, &c.Title, &c.GoalPoints, &c.IsActive)
		if err != nil {
			return nil, err
		}
		challenges = append(challenges, c)
	}
	return challenges, nil
}

func getAllUsers() ([]*User, error) {
	rows, err := db.Query(`SELECT * FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*User, 0)
	for rows.Next() {
		u := new(User)
		err := rows.Scan(&u.ID, &u.Name, &u.PasswordHash, &u.IsAdmin)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func getChallengePot(c *Challenge) (int, error) {
	var pot int
	err := db.QueryRow(`SELECT sum(points) FROM performances
					 	WHERE challenge_id = ?`, c.ID).Scan(&pot)
	if err != nil {
		return 0, err
	}

	return pot, nil
}

func addPerformance(points int, u *User, c *Challenge) error {
	_, err := db.Exec(`INSERT INTO performances(points, user_id, challenge_id)
					VALUES (?,?,?)`, points, u.ID, c.ID)
	if err != nil {
		return err
	}

	return nil
}

func addChallenge(title string, goal int, isActive bool) error {
	_, err := db.Exec(`
		INSERT INTO challenges(title, goal_points, is_active)
		VALUES(?, ?, ?)`, title, goal, isActive)

	if err != nil {
		return err
	}

	return nil
}
