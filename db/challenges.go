package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-chi/jwtauth/v5"
)

func GetChallenge(ctx context.Context, id int) (*Challenge, error) {
	c := new(Challenge)
	err := DB.QueryRowContext(ctx,
		"SELECT * FROM challenges WHERE id = ? LIMIT 1", id).
		Scan(&c.ID, &c.Title, &c.GoalPoints, &c.StartDate, &c.EndDate)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func GetAllChallenges(ctx context.Context, active bool) ([]*Challenge, error) {
	rows, err := DB.QueryContext(ctx, "SELECT * FROM challenges")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	challenges := make([]*Challenge, 0)
	for rows.Next() {
		c := new(Challenge)
		err := rows.Scan(&c.ID, &c.Title, &c.GoalPoints,
			&c.StartDate, &c.EndDate)
		if err != nil {
			return nil, err
		}
		if active {
			var startDate, endDate time.Time
			var err error
			startDate, err = time.Parse(time.DateOnly, c.StartDate)
			if err != nil {
				fmt.Printf("unable to parse date")
				continue
			}
			endDate, err = time.Parse(time.DateOnly, c.EndDate)
			if err != nil {
				fmt.Printf("unable to parse date")
				continue
			}
			if TimeIsBetween(time.Now(), startDate, endDate) {
				challenges = append(challenges, c)
			} else {
				continue
			}
		} else {
			challenges = append(challenges, c)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	for _, c := range challenges {
		pRows, err := DB.QueryContext(ctx, `
			SELECT * FROM performances
			WHERE challenge_id = ?`, c.ID)
		if err != nil {
			return nil, err
		}

		var pot int

		for pRows.Next() {
			p := new(Performance)
			err = pRows.Scan(&p.ID, &p.Points, &p.UserID,
				&p.ChallengeID, &p.CreatedAt)
			if err != nil {
				pRows.Close()
				return nil, err
			}
			pot += p.Points

			c.Performances = append(c.Performances, p)
		}

		if err := pRows.Err(); err != nil {
			pRows.Close()
			return nil, err
		}

		for _, p := range c.Performances {
			u := new(User)
			err := DB.QueryRowContext(ctx, `
				SELECT id, name FROM users WHERE id = ?
			`, p.UserID).Scan(&u.ID, &u.Name)
			if err == nil {
				p.User = u
			}
		}

		pRows.Close()
		c.Pot = pot
	}
	return challenges, nil
}

func GetLatestChallenge(ctx context.Context) (*Challenge, error) {
	c := new(Challenge)
	err := DB.QueryRowContext(ctx,
		"SELECT * FROM challenges WHERE date() > end_date LIMIT 1").
		Scan(&c.ID, &c.Title, &c.GoalPoints, &c.StartDate, &c.EndDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	var pRows *sql.Rows
	pRows, err = DB.QueryContext(ctx, `
			SELECT * FROM performances
			WHERE challenge_id = ?`, c.ID)
	if err != nil {
		return nil, err
	}

	var pot int

	for pRows.Next() {
		p := new(Performance)
		err = pRows.Scan(&p.ID, &p.Points, &p.UserID,
			&p.ChallengeID, &p.CreatedAt)
		if err != nil {
			pRows.Close()
			return nil, err
		}
		pot += p.Points

		c.Performances = append(c.Performances, p)
	}

	if err := pRows.Err(); err != nil {
		pRows.Close()
		return nil, err
	}

	for _, p := range c.Performances {
		u := new(User)
		err := DB.QueryRowContext(ctx, `
				SELECT id, name FROM users WHERE id = ?
			`, p.UserID).Scan(&u.ID, &u.Name)
		if err == nil {
			p.User = u
		}
	}

	pRows.Close()
	c.Pot = pot

	return c, nil
}

func AddChallenge(c *Challenge) error {
	if !IsValidDate(c.StartDate) || !IsValidDate(c.EndDate) {
		return fmt.Errorf("invalid dates")
	}
	_, err := DB.Exec(`
		INSERT INTO challenges(title, goal_points, start_date, end_date)
		VALUES(?, ?, ?, ?)`, c.Title, c.GoalPoints, c.StartDate, c.EndDate)

	if err != nil {
		return err
	}

	return nil
}

func DeleteChallenge(id int) error {
	_, err := DB.Exec(`
		DELETE FROM challenges
		WHERE id = ?`, id)

	if err != nil {
		return err
	}

	return nil
}

func UpdateChallenge(c *Challenge) error {
	if !IsValidDate(c.StartDate) || !IsValidDate(c.EndDate) {
		return fmt.Errorf("invalid dates")
	}
	_, err := DB.Exec(`
		UPDATE challenges
		SET title = ?, goal_points = ?, start_date = ?, end_date = ?
		WHERE id = ?`, c.Title, c.GoalPoints, c.StartDate, c.EndDate, c.ID)

	if err != nil {
		return err
	}

	return nil
}

func AddPerformance(ctx context.Context, p *Performance) error {
	_, claims, _ := jwtauth.FromContext(ctx)
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return fmt.Errorf("invalid user id in session")
	}

	_, err := DB.Exec(`INSERT INTO performances(points, user_id, challenge_id)
					VALUES (?,?,?)`, p.Points, int(userIDFloat), p.ChallengeID)
	if err != nil {
		return err
	}

	return nil
}