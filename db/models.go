package db

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Collection struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	GoalPoints  int    `json:"goalPoints"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	CreatedAt   string `json:"createdAt"`
	UserID      int    `json:"userId"`

	Challenges []*Challenge `json:"challenges"`
}

type Challenge struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Unit         string `json:"unit"`
	StartDate    string `json:"startDate"`
	EndDate      string `json:"endDate"`
	CreatedAt    string `json:"createdAt"`
	UserID       int    `json:"userId"`
	CollectionID int    `json:"collectionId"`

	Performances []*Performance `json:"performances"`
}

type User struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"passwordHash,omitempty"`
	IsAdmin      bool   `json:"isAdmin"`
	Color        string `json:"color"`
}

func (u *User) Verify(plain string) bool {
	if u == nil {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(plain))
	return err == nil
}

func (u *User) SetPasswordHash(plain string) error {
	if u == nil {
		return errors.New("user is not valid")
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), 12)
	if err != nil {
		return err
	}
	u.PasswordHash = string(bytes)

	return nil
}

type Performance struct {
	ID          int    `json:"id"`
	Points      int    `json:"points"`
	CreatedAt   string `json:"createdAt"`
	UserID      int    `json:"userId"`
	ChallengeID int    `json:"challengeId"`
}
