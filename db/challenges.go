package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

func (s *Store) GetChallengeByID(ctx context.Context, id int) (*Challenge, error) {
	c := new(Challenge)
	err := s.db.QueryRowContext(ctx,
		"SELECT * FROM challenges WHERE id = ? LIMIT 1", id).
		Scan(&c.ID, &c.Title, &c.Description, &c.Unit, &c.GoalPoints, &c.StartDate, &c.EndDate, &c.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}
	return c, nil
}

func (s *Store) ListChallenges(ctx context.Context, active bool) ([]*Challenge, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT * FROM challenges")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	challenges := make([]*Challenge, 0)
	for rows.Next() {
		c := new(Challenge)
		err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.Unit, &c.GoalPoints,
			&c.StartDate, &c.EndDate, &c.UserID)
		if err != nil {
			return nil, err
		}
		if active {
			var startDate, endDate time.Time
			var err error
			startDate, err = time.Parse(time.DateOnly, c.StartDate)
			if err != nil {
				fmt.Printf("unable to parse date, skipping")
				continue
			}
			endDate, err = time.Parse(time.DateOnly, c.EndDate)
			if err != nil {
				fmt.Printf("unable to parse date, skipping")
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
		u := new(User)
		err := s.db.QueryRowContext(ctx, `
				SELECT id, name, color FROM users WHERE id = ?
			`, c.UserID).Scan(&u.ID, &u.Name, &u.Color)
		if err == nil {
			c.User = u
		}
	}

	for _, c := range challenges {
		pRows, err := s.db.QueryContext(ctx, `
			SELECT * FROM performances
			WHERE challenge_id = ?`, c.ID)
		if err != nil {
			return nil, err
		}

		var pot int

		performances := make([]*Performance, 0)
		for pRows.Next() {
			p := new(Performance)
			err = pRows.Scan(&p.ID, &p.Points, &p.UserID,
				&p.ChallengeID, &p.CreatedAt)
			if err != nil {
				pRows.Close()
				return nil, err
			}
			pot += p.Points

			performances = append(performances, p)
		}

		if err := pRows.Err(); err != nil {
			pRows.Close()
			return nil, err
		}

		c.Performances = performances

		for _, p := range c.Performances {
			u := new(User)
			// Leaving user performances out is intentional
			err := s.db.QueryRowContext(ctx, `
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

func (s *Store) AddChallenge(ctx context.Context, c Challenge) error {
	_, claims, _ := jwtauth.FromContext(ctx)
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return fmt.Errorf("invalid user id in session")
	}
	if !IsValidDate(c.StartDate) || !IsValidDate(c.EndDate) {
		return fmt.Errorf("invalid dates")
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO challenges(title, description, unit, goal_points, start_date, end_date, user_id)
		VALUES(?, ?, ?, ?, ?, ?, ?)`, c.Title, c.Description, c.Unit, c.GoalPoints, c.StartDate, c.EndDate, int(userIDFloat))

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) DeleteChallenge(ctx context.Context, id int) error {
	token, claims, _ := jwtauth.FromContext(ctx)

	if token == nil || jwt.Validate(token) != nil {
		return fmt.Errorf("unable to validate user token")
	}
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return fmt.Errorf("unable to convert user_id to float")
	}

	u := new(User)
	err := s.db.QueryRowContext(ctx, "SELECT id, is_admin FROM users WHERE id = ? LIMIT 1", int(userIDFloat)).Scan(&u.ID, &u.IsAdmin)
	if err != nil {
		return err
	}

	ch := new(Challenge)
	err = s.db.QueryRowContext(ctx, "SELECT id, user_id FROM challenges WHERE id = ? LIMIT 1", id).Scan(&ch.ID, &ch.UserID)
	if err != nil {
		return err
	}

	if int(userIDFloat) != ch.UserID && !u.IsAdmin {
		return fmt.Errorf("no permisson to update")
	}

	_, err = s.db.ExecContext(ctx, `
		DELETE FROM challenges
		WHERE id = ?`, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdateChallenge(ctx context.Context, id int, c Challenge) error {
	token, claims, _ := jwtauth.FromContext(ctx)

	if token == nil || jwt.Validate(token) != nil {
		return fmt.Errorf("unable to validate user token")
	}
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return fmt.Errorf("unable to convert user_id to float")
	}

	u := new(User)
	err := s.db.QueryRowContext(ctx, "SELECT id, is_admin FROM users WHERE id = ? LIMIT 1", int(userIDFloat)).Scan(&u.ID, &u.IsAdmin)
	if err != nil {
		return err
	}

	ch := new(Challenge)
	err = s.db.QueryRowContext(ctx, "SELECT id, user_id FROM challenges WHERE id = ? LIMIT 1", id).Scan(&ch.ID, &ch.UserID)
	if err != nil {
		return err
	}

	if int(userIDFloat) != ch.UserID && !u.IsAdmin {
		return fmt.Errorf("no permisson to update")
	}

	if !IsValidDate(c.StartDate) || !IsValidDate(c.EndDate) {
		return fmt.Errorf("invalid dates")
	}
	_, err = s.db.ExecContext(ctx, `
		UPDATE challenges
		SET title = ?, description = ?, unit = ?, goal_points = ?, start_date = ?, end_date = ?
		WHERE id = ?`, c.Title, c.Description, c.Unit, c.GoalPoints, c.StartDate, c.EndDate, id)

	if err != nil {
		return err
	}

	return nil
}
