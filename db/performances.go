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

func (s *Store) GetPerformanceByID(ctx context.Context, id int) (*Performance, error) {
	p := new(Performance)
	err := s.db.QueryRowContext(ctx,
		"SELECT id, points, created_at, user_id, challenge_id FROM performances WHERE id = ? LIMIT 1", id).
		Scan(&p.ID, &p.Points, &p.CreatedAt, &p.UserID, &p.ChallengeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}
	return p, nil
}

func (s *Store) AddPerformance(ctx context.Context, p *Performance) error {
	_, claims, _ := jwtauth.FromContext(ctx)
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return fmt.Errorf("invalid user id in session")
	}

	_, err := s.db.ExecContext(ctx, `INSERT INTO performances(points, user_id, challenge_id)
					VALUES (?,?,?)`, p.Points, int(userIDFloat), p.ChallengeID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) DeletePerformance(ctx context.Context, id int) error {
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

	p := new(Performance)
	err = s.db.QueryRowContext(ctx, "SELECT id, user_id, challenge_id FROM performances WHERE id = ? LIMIT 1", id).Scan(&p.ID, &p.UserID, &p.ChallengeID)
	if err != nil {
		return err
	}

	c := new(Challenge)
	err = s.db.QueryRowContext(ctx, "SELECT end_date FROM challenges WHERE id = ? LIMIT 1", p.ChallengeID).Scan(&c.EndDate)

	endDate, err := time.Parse(time.DateOnly, c.EndDate)
	if err != nil {
		return err
	}

	if endDate.Before(time.Now()) && !u.IsAdmin {
		return fmt.Errorf("cannot delete performance if challenge has expired")
	}

	if int(userIDFloat) != p.UserID && !u.IsAdmin {
		return fmt.Errorf("no permisson to delete performance")
	}

	_, err = s.db.ExecContext(ctx, `
		DELETE FROM performances
		WHERE id = ?`, id)

	if err != nil {
		return err
	}

	return nil
}
