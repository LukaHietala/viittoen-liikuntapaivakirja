package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	ErrChallengeNotFound = errors.New("challenge not found")
)

func (s *Store) ListChallenges(ctx context.Context) ([]*Challenge, error) {
	query := `
		SELECT id, title, description, unit, goal_points, start_date, end_date, user_id, collection_id, created_at
		FROM challenges
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var challenges []*Challenge
	for rows.Next() {
		c := &Challenge{}
		err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.Unit, &c.GoalPoints, &c.StartDate, &c.EndDate, &c.UserID, &c.CollectionID, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		challenges = append(challenges, c)
	}

	for _, challenge := range challenges {
		performances, err := s.ListChallengePerformances(ctx, challenge.ID)
		if err != nil {
			return challenges, err
		}
		challenge.Performances = performances
	}

	return challenges, rows.Err()
}

func (s *Store) ListActiveChallenges(ctx context.Context) ([]*Challenge, error) {
	query := `
		SELECT id, title, description, unit, goal_points, start_date, end_date, user_id, collection_id, created_at
		FROM challenges
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var challenges []*Challenge
	for rows.Next() {
		c := &Challenge{}
		err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.Unit, &c.GoalPoints, &c.StartDate, &c.EndDate, &c.UserID, &c.CollectionID, &c.CreatedAt)
		if err != nil {
			return nil, err
		}

		var startDate, endDate time.Time
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

		if timeIsBetween(time.Now(), startDate, endDate) {
			challenges = append(challenges, c)
		}
	}

	for _, challenge := range challenges {
		performances, err := s.ListChallengePerformances(ctx, challenge.ID)
		if err != nil {
			return challenges, err
		}
		challenge.Performances = performances
	}

	return challenges, rows.Err()
}

func (s *Store) ListCollectionChallenges(ctx context.Context, id int) ([]*Challenge, error) {
	query := `
		SELECT id, title, description, unit, goal_points, start_date, end_date, user_id, collection_id, created_at
		FROM challenges
		WHERE collection_id  = ?
	`

	rows, err := s.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	challenges := make([]*Challenge, 0)
	for rows.Next() {
		c := &Challenge{}
		err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.Unit, &c.GoalPoints, &c.StartDate, &c.EndDate, &c.UserID, &c.CollectionID, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		challenges = append(challenges, c)
	}

	for _, challenge := range challenges {
		performances, err := s.ListChallengePerformances(ctx, challenge.ID)
		if err != nil {
			return challenges, err
		}
		challenge.Performances = performances
	}

	return challenges, rows.Err()
}

func (s *Store) FindChallengeByID(ctx context.Context, id int) (*Challenge, error) {
	query := `
		SELECT id, title, description, unit, goal_points, start_date, end_date, user_id, collection_id, created_at
		FROM challenges
		WHERE id = ? LIMIT 1
	`

	var challenge Challenge
	err := s.db.QueryRowContext(ctx, query, id).
		Scan(&challenge.ID, &challenge.Title, &challenge.Description, &challenge.Unit, &challenge.GoalPoints, &challenge.StartDate, &challenge.EndDate, &challenge.UserID, &challenge.CollectionID, &challenge.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w; id: %d", ErrChallengeNotFound, id)
	}
	if err != nil {
		return nil, err
	}

	performances, err := s.ListChallengePerformances(ctx, challenge.ID)
	if err != nil {
		return &challenge, err
	}
	challenge.Performances = performances

	return &challenge, nil
}

func (s *Store) AddChallenge(ctx context.Context, challenge *Challenge) error {
	query := `
		INSERT INTO challenges(title, description, unit, goal_points, start_date, end_date, user_id, collection_id)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.ExecContext(ctx, query, challenge.Title, challenge.Description, challenge.Unit, challenge.GoalPoints, challenge.StartDate, challenge.EndDate, challenge.UserID, challenge.CollectionID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdateChallenge(ctx context.Context, challenge *Challenge) error {
	query := `
		UPDATE challenges
		SET title = ?, description = ?, unit = ?, goal_points = ?, start_date = ?, end_date = ?
		WHERE id = ?
	`

	_, err := s.db.ExecContext(ctx, query, challenge.Title, challenge.Description, challenge.Unit, challenge.GoalPoints, challenge.StartDate, challenge.EndDate, challenge.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) DeleteChallengeByID(ctx context.Context, id int) error {
	query := `
		DELETE FROM challenges
		WHERE id = ?
	`

	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
