package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

func (s *Store) ListUsers(ctx context.Context) ([]*User, error) {
	query := `
		SELECT id, name, email, is_admin, color
		FROM users
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.IsAdmin, &user.Color)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

func (s *Store) FindUserByID(ctx context.Context, id int) (*User, error) {
	query := `
		SELECT id, name, email, is_admin, color
		FROM users
		WHERE id = ? LIMIT 1
	`

	var user User
	err := s.db.QueryRowContext(ctx, query, id).
		Scan(&user.ID, &user.Name, &user.Email, &user.IsAdmin, &user.Color)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w; id: %d", ErrUserNotFound, id)
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Store) FindUserByName(ctx context.Context, name string) (*User, error) {
	query := `
		SELECT id, name, email, password_hash, is_admin, color
		FROM users
		WHERE name = ? LIMIT 1
	`

	var user User
	err := s.db.QueryRowContext(ctx, query, name).
		Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.IsAdmin, &user.Color)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w; name: %s", ErrUserNotFound, name)
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Store) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, name, email, password_hash, is_admin, color
		FROM users
		WHERE email = ? LIMIT 1
	`

	var user User
	err := s.db.QueryRowContext(ctx, query, email).
		Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.IsAdmin, &user.Color)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w; email: %s", ErrUserNotFound, email)
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Store) AddUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users(name, email, password_hash, is_admin, color)
		VALUES(?, ?, ?, ?, ?)
	`

	_, err := s.db.ExecContext(ctx, query, user.Name, user.Email, user.PasswordHash, user.IsAdmin, user.Color)
	if err != nil {
		return err
	}

	// TODO: Rows affected check if necessary (applies to other delete/updates)

	return nil
}

func (s *Store) UpdateUser(ctx context.Context, user *User) error {
	query := `
    	UPDATE users
		SET name = ?, email = ?, is_admin = ?, color = ?
		WHERE id = ?
    `

	_, err := s.db.ExecContext(ctx, query, user.Name, user.Email, user.IsAdmin, user.Color, user.ID)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) DeleteUserByID(ctx context.Context, id int) error {
	query := `
		DELETE FROM users
		WHERE id = ?
	`

	_, err := s.db.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) UserCount() (int, error) {
	query := `
		SELECT COUNT(*)
		FROM users
	`

	var count int
	err := s.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
