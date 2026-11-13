package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	ErrCollectionNotFound = errors.New("collection not found")
)

func (s *Store) ListCollections(ctx context.Context) ([]*Collection, error) {
	query := `
		SELECT id, title, description, goal_points, start_date, end_date, user_id, created_at
		FROM collections
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collections []*Collection
	for rows.Next() {
		c := &Collection{}
		err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.GoalPoints, &c.StartDate, &c.EndDate, &c.UserID, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		collections = append(collections, c)
	}

	for _, collection := range collections {
		challenges, err := s.ListCollectionChallenges(ctx, collection.ID)
		if err != nil {
			return collections, err
		}
		collection.Challenges = challenges
	}

	return collections, rows.Err()
}

func (s *Store) ListActiveCollections(ctx context.Context) ([]*Collection, error) {
	query := `
		SELECT id, title, description, goal_points, start_date, end_date, user_id, created_at
		FROM collections
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collections []*Collection
	for rows.Next() {
		c := &Collection{}
		err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.GoalPoints, &c.StartDate, &c.EndDate, &c.UserID, &c.CreatedAt)
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
			collections = append(collections, c)
		}
	}

	for _, collection := range collections {
		challenges, err := s.ListCollectionChallenges(ctx, collection.ID)
		if err != nil {
			return collections, err
		}
		collection.Challenges = challenges
	}

	return collections, rows.Err()
}

func (s *Store) FindCollectionByID(ctx context.Context, id int) (*Collection, error) {
	query := `
		SELECT id, title, description, goal_points, start_date, end_date, user_id, created_at
		FROM collections 
		WHERE id = ? LIMIT 1
	`

	var collection Collection
	err := s.db.QueryRowContext(ctx, query, id).
		Scan(&collection.ID, &collection.Title, &collection.Description, &collection.GoalPoints, &collection.StartDate, &collection.EndDate, &collection.UserID, &collection.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w; id: %d", ErrCollectionNotFound, id)
	}
	if err != nil {
		return nil, err
	}

	challenges, err := s.ListCollectionChallenges(ctx, collection.ID)
	if err != nil {
		return &collection, err
	}
	collection.Challenges = challenges

	return &collection, nil
}

func (s *Store) AddCollection(ctx context.Context, collection *Collection) error {
	query := `
		INSERT INTO collections(title, description, goal_points, start_date, end_date, user_id)
		VALUES(?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.ExecContext(ctx, query, collection.Title, collection.Description, collection.GoalPoints, collection.StartDate, collection.EndDate, collection.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdateCollection(ctx context.Context, collection *Collection) error {
	query := `
		UPDATE collections
		SET title = ?, description = ?, goal_points = ?, start_date = ?, end_date = ?
		WHERE id = ?
	`

	_, err := s.db.ExecContext(ctx, query, collection.Title, collection.Description, collection.GoalPoints, collection.StartDate, collection.EndDate, collection.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) DeleteCollectionByID(ctx context.Context, id int) error {
	query := `
		DELETE FROM collections
		WHERE id = ?
	`

	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
