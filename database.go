package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Challenge struct {
	ID         int `json:"id"`
	Title      string `json:"title"`
	GoalPoints int `json:"goal_points"`
	IsActive   bool `json:"is_active"`
	Pot int `json:"pot"`
	Performances []*Performance `json:"performances"`
}

type User struct {
	ID           int `json:"id"`
	Name         string `json:"name"`
	PasswordHash string `json:"password_hash"`
	IsAdmin      bool `json:"is_admin"`
	Performances []*Performance `json:"performances"`
}

type Performance struct {
	ID     int `json:"id"`
	Points int `json:"points"`
	UserID int `json:"user_id"`
	ChallengeID int `json:"challenge_id"`
}

func connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}
	err = db.Ping()
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
	tx, _ := db.Begin()
	rows, err := tx.Query(`SELECT * FROM challenges`)
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

	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	for _, c := range challenges {
		rows, err = tx.Query(`
			SELECT * FROM performances
			WHERE challenge_id = ?`, c.ID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var pot int
		
		for rows.Next() {
			p := new(Performance)
			err = rows.Scan(&p.ID, &p.Points, &p.UserID, &p.ChallengeID)
			pot += p.Points
			if err != nil {
				return nil, err
			}
			c.Performances = append(c.Performances, p)
		}

		c.Pot = pot

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return challenges, nil
}

func getAllUsers() ([]*User, error) {
	tx, _ := db.Begin()
	rows, err := tx.Query(`SELECT * FROM users`)
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

	if err = rows.Err(); err != nil {
		return nil, err
	}

	for _, u := range users {
		rows, err = tx.Query(`
			SELECT * FROM performances
			WHERE user_id = ?`, u.ID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			p := new(Performance)
			err = rows.Scan(&p.ID, &p.Points, &p.UserID, &p.ChallengeID)
			if err != nil {
				return nil, err
			}
			u.Performances = append(u.Performances, p)
		}
		
		if err = rows.Err(); err != nil {
			return nil, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func addPerformance(points, user_id, challenge_id int) error {
	_, err := db.Exec(`INSERT INTO performances(points, user_id, challenge_id)
					VALUES (?,?,?)`, points, user_id, challenge_id)
	if err != nil {
		return err
	}

	return nil
}

func getActiveChallenge() (*Challenge, error) {
	c := new(Challenge)
	tx, err := db.Begin()
	err = tx.QueryRow(`SELECT * FROM challenges WHERE is_active = 1 LIMIT 1`).Scan(&c.ID, &c.Title, &c.GoalPoints, &c.IsActive)

	if err != nil {
		return nil, err
	}
	
	rows, err := tx.Query(`
		SELECT * FROM performances
		WHERE challenge_id = ?`, c.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pot int

	for rows.Next() {
		p := new(Performance)
		err = rows.Scan(&p.ID, &p.Points, &p.UserID, &p.ChallengeID)
		pot += p.Points
		if err != nil {
			return nil, err
		}
		c.Performances = append(c.Performances, p)
	}

	c.Pot = pot
	
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return c, nil
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

func deleteChallenge(id int) error {
	_, err := db.Exec(`
		DELETE FROM challenges
		WHERE id = ?`, id)

	if err != nil {
		return err
	}

	return nil
}