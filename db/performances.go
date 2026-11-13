package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrPerformanceNotFound = errors.New("performance not found")
)

func (s *Store) FindPerformanceByID(ctx context.Context, id int) (*Performance, error) {
	query := `
		SELECT id, points, created_at, user_id, challenge_id
		FROM performances
		WHERE id = ? LIMIT 1
	`

	var performance Performance
	err := s.db.QueryRowContext(ctx, query, id).
		Scan(&performance.ID, &performance.Points, &performance.CreatedAt, &performance.UserID, &performance.ChallengeID)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w; id: %d", ErrPerformanceNotFound, id)
	}
	if err != nil {
		return nil, err
	}

	return &performance, nil
}

func (s *Store) AddPerformance(ctx context.Context, performance *Performance) error {
	query := `
		INSERT INTO performances(points, user_id, challenge_id)
		VALUES (?,?,?)
	`

	_, err := s.db.ExecContext(ctx, query, performance.Points, performance.UserID, performance.ChallengeID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) DeletePerformanceByID(ctx context.Context, id int) error {
	query := `
		DELETE FROM performances
		WHERE id = ?
	`
	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) ListPerformances(ctx context.Context) ([]*Performance, error) {
	query := `
		SELECT id, points, user_id, challenge_id, created_at
		FROM performances
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var performances []*Performance
	for rows.Next() {
		p := &Performance{}
		err := rows.Scan(&p.ID, &p.Points, &p.UserID, &p.ChallengeID, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		performances = append(performances, p)
	}

	return performances, rows.Err()
}

func (s *Store) ListChallengePerformances(ctx context.Context, id int) ([]*Performance, error) {
	query := `
		SELECT id, points, user_id, challenge_id, created_at
		FROM performances
		WHERE challenge_id = ?
	`

	rows, err := s.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var performances []*Performance
	for rows.Next() {
		p := &Performance{}
		err := rows.Scan(&p.ID, &p.Points, &p.UserID, &p.ChallengeID, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		performances = append(performances, p)
	}

	return performances, rows.Err()
}

func (s *Store) ListUserPerformances(ctx context.Context, id int) ([]*Performance, error) {
	query := `
		SELECT id, points, user_id, challenge_id, created_at
		FROM performances
		WHERE user_id = ?
	`

	rows, err := s.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var performances []*Performance
	for rows.Next() {
		p := &Performance{}
		err := rows.Scan(&p.ID, &p.Points, &p.UserID, &p.ChallengeID, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		performances = append(performances, p)
	}

	return performances, rows.Err()
}
