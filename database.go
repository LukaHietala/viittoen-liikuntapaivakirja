package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-chi/jwtauth/v5"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type Challenge struct {
	ID           int            `json:"id"`
	Title        string         `json:"title"`
	GoalPoints   int            `json:"goal_points"`
	IsActive     bool           `json:"is_active"`
	Pot          int            `json:"pot"`
	Performances []*Performance `json:"performances"`
}

type User struct {
	ID           int            `json:"id"`
	Name         string         `json:"name"`
	PasswordHash string         `json:"password_hash"`
	IsAdmin      bool           `json:"is_admin"`
	Performances []*Performance `json:"performances"`
}

type Performance struct {
	ID          int `json:"id"`
	Points      int `json:"points"`
	UserID      int `json:"user_id"`
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

		INSERT INTO users VALUES(NULL, "Jaakko", "$2a$14$dhSvJi8wLpc0iAB5LW91Le4GKK/w9i7IKyZ6tgE7L8xnW4b2S2/lG", TRUE);
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

func getAllChallenges(ctx context.Context) ([]*Challenge, error) {
	tx, _ := db.BeginTx(ctx, nil)
	rows, err := tx.QueryContext(ctx, `SELECT * FROM challenges`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	challenges := make([]*Challenge, 0)
	for rows.Next() {
		c := new(Challenge)
		err := rows.Scan(&c.ID, &c.Title, &c.GoalPoints, &c.IsActive)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("no challenges")
			}
			return nil, err
		}
		challenges = append(challenges, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	for _, c := range challenges {
		rows, err = tx.QueryContext(ctx, `
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
				if err == sql.ErrNoRows {
					return nil, nil
				}
				return nil, err
			}
			c.Performances = append(c.Performances, p)
		}

		c.Pot = pot

		if err != nil {
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

func addPerformance(ctx context.Context, p *Performance) error {
	_, claims, _ := jwtauth.FromContext(ctx)
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return fmt.Errorf("invalid user id in session")
	}

	_, err := db.Exec(`INSERT INTO performances(points, user_id, challenge_id)
					VALUES (?,?,?)`, p.Points, int(userIDFloat), p.ChallengeID)
	if err != nil {
		return err
	}

	return nil
}

func getActiveChallenge() (*Challenge, error) {
	c := new(Challenge)
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	err = tx.QueryRow(`SELECT * FROM challenges WHERE is_active = 1`).Scan(&c.ID, &c.Title, &c.GoalPoints, &c.IsActive)

	if err != nil {
		return nil, err
	}

	err = nil
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
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
		c.Performances = append(c.Performances, p)
	}

	c.Pot = pot

	if err = rows.Err(); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func addChallenge(c *Challenge) error {
	_, err := db.Exec(`
		INSERT INTO challenges(title, goal_points, is_active)
		VALUES(?, ?, ?)`, c.Title, c.GoalPoints, c.IsActive)

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

func updateChallenge(c *Challenge) error {
	_, err := db.Exec(`
		UPDATE challenges
		SET title = ?, is_active = ?, goal_points = ?
		WHERE id = ?`, c.Title, c.IsActive, c.GoalPoints, c.ID)

	if err != nil {
		return err
	}

	return nil
}

func getUser(id int) (*User, error) {
	u := new(User)
	err := db.QueryRow("SELECT * FROM users WHERE id = ?", id).Scan(&u.ID, &u.Name, &u.PasswordHash, &u.IsAdmin)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found based on id: %d", id)
		}
		return nil, err
	}
	return u, nil
}

func deleteUser(id int) error {
	_, err := db.Exec(`
		DELETE FROM users
		WHERE id = ?`, id)

	if err != nil {
		return err
	}

	return nil
}

// TODO: move to auth package
func VerifyHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func verifyUser(name, password string) (int, error) {
	u := new(User)
	row := db.QueryRow("SELECT * FROM users WHERE name = ?", name)
	err := row.Scan(&u.ID, &u.Name, &u.PasswordHash, &u.IsAdmin)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("no user found with the name of: %s", name)
		}
		return 0, err
	}

	if VerifyHash(password, u.PasswordHash) {
		return u.ID, nil
	} else {
		return 0, fmt.Errorf("wrong password")
	}
}