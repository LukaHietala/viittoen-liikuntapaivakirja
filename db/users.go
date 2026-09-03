package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v3/jwt"
)

func GetAllUsers() ([]*User, error) {
	rows, err := DB.Query(`SELECT * FROM users`)
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
		uRows, err := DB.Query(`
			SELECT * FROM performances
			WHERE user_id = ?`, u.ID)
		if err != nil {
			return nil, err
		}

		for uRows.Next() {
			p := new(Performance)
			err = uRows.Scan(&p.ID, &p.Points, &p.UserID, &p.ChallengeID, &p.CreatedAt)
			if err != nil {
				uRows.Close()
				return nil, err
			}
			u.Performances = append(u.Performances, p)
		}

		if err = uRows.Err(); err != nil {
			uRows.Close()
			return nil, err
		}

		uRows.Close()
	}

	return users, nil
}

func GetUser(id int) (*User, error) {
	u := new(User)
	err := DB.QueryRow("SELECT * FROM users WHERE id = ? LIMIT 1", id).Scan(&u.ID, &u.Name, &u.PasswordHash, &u.IsAdmin)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user found based on id: %d", id)
		}
		return nil, err
	}
	return u, nil
}

func AddUser(u *User) error {
	hash, err := HashPassword(u.PasswordPlain)
	if err != nil {
		return err
	}
	_, err = DB.Exec(`
		INSERT INTO users(name, password_hash, is_admin)
		VALUES(?, ?, ?)`, u.Name, hash, u.IsAdmin)

	if err != nil {
		return err
	}

	return nil
}

func UpdateUser(ctx context.Context, u *User) error {
	token, claims, _ := jwtauth.FromContext(ctx)

	if token == nil || jwt.Validate(token) != nil {
		return fmt.Errorf("unable to validate user token")
	}
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return fmt.Errorf("unable to convert user_id to float")
	}
	if u.ID == int(userIDFloat) && !u.IsAdmin {
		return fmt.Errorf("can't remove own admin privileges")
	}

	var err error
	var hash string
	if u.PasswordPlain != "" {
		hash, err = HashPassword(u.PasswordPlain)
		_, err = DB.ExecContext(ctx, `
		UPDATE users
		SET name = ?, password_hash = ?, is_admin = ?
		WHERE id = ?`, u.Name, hash, u.IsAdmin, u.ID)
	} else {
		_, err = DB.ExecContext(ctx, `
		UPDATE users
		SET name = ?, is_admin = ?
		WHERE id = ?`, u.Name, u.IsAdmin, u.ID)
	}

	// TODO: forgot pass
	if err != nil {
		return err
	}

	return nil
}

func DeleteUser(ctx context.Context, id int) error {
	token, claims, _ := jwtauth.FromContext(ctx)

	if token == nil || jwt.Validate(token) != nil {
		return fmt.Errorf("unable to validate user token")
	}
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return fmt.Errorf("unable to convert user_id to float")
	}
	if id == int(userIDFloat) {
		return fmt.Errorf("can't remove self")
	}

	_, err := DB.ExecContext(ctx, `
		DELETE FROM users
		WHERE id = ?`, id)

	if err != nil {
		return err
	}

	return nil
}

func GetSelf(ctx context.Context) (*User, error) {
	token, claims, _ := jwtauth.FromContext(ctx)

	if token == nil || jwt.Validate(token) != nil {
		return nil, fmt.Errorf("unable to validate user token")
	}
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("unable to convert user_id to float")
	}

	u := new(User)
	err := DB.QueryRow("SELECT id, name, is_admin FROM users WHERE id = ? LIMIT 1", int(userIDFloat)).Scan(&u.ID, &u.Name, &u.IsAdmin)
	if err != nil {
		return nil, err
	}
	rows, err := DB.Query(`
		SELECT * FROM performances
		WHERE user_id = ?`, u.ID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		p := new(Performance)
		err = rows.Scan(&p.ID, &p.Points, &p.UserID, &p.ChallengeID)
		if err != nil {
			rows.Close()
			return nil, err
		}
		u.Performances = append(u.Performances, p)
	}

	if err = rows.Err(); err != nil {
		rows.Close()
		return nil, err
	}

	rows.Close()

	return u, nil
}

func GetSession(id int) (*User, error) {
	u := new(User)
	err := DB.QueryRow("SELECT id, name, is_admin FROM users WHERE id = ? LIMIT 1", id).Scan(&u.ID, &u.Name, &u.IsAdmin)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func VerifyUser(name, password string) (int, error) {
	u := new(User)
	row := DB.QueryRow("SELECT * FROM users WHERE name = ?", name)
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
