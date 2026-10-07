package db

import (
	"errors"
	"net/mail"
	"time"
)

type Challenge struct {
	ID           int            `json:"id"`
	Title        string         `json:"title"`
	Description  string         `json:"description"`
	Unit         string         `json:"unit"`
	GoalPoints   int            `json:"goal_points"`
	StartDate    string         `json:"start_date"`
	EndDate      string         `json:"end_date"`
	Pot          int            `json:"pot"`
	UserID       int            `json:"user_id"`
	Performances []*Performance `json:"performances"`
	User         *User          `json:"user"`
}

type User struct {
	ID            int            `json:"id"`
	Name          string         `json:"name"`
	Email         string         `json:"email"`
	PasswordHash  string         `json:"password_hash,omitempty"`
	PasswordPlain string         `json:"password_plain,omitempty"`
	IsAdmin       bool           `json:"is_admin,omitempty"`
	Color         string         `json:"color"`
	Performances  []*Performance `json:"performances"`
}

type Performance struct {
	ID          int        `json:"id"`
	Points      int        `json:"points"`
	CreatedAt   string     `json:"created_at"`
	UserID      int        `json:"user_id"`
	ChallengeID int        `json:"challenge_id"`
	User        *User      `json:"user"`
	Challenge   *Challenge `json:"challenge"`
}

func (c *Challenge) Validate() error {
	if c == nil {
		return errors.New("Virheellinen haaste")
	}

	if c.Title == "" {
		return errors.New("Haasteen nimi puuttuu")
	}

	if len(c.Title) < 1 || len(c.Title) > 255 {
		return errors.New("Haasteen nimi voi olla 1-255 merkkiä pitkä")
	}

	if c.Description == "" {
		return errors.New("Haasteen kuvaus puuttuu")
	}

	if len(c.Description) < 1 || len(c.Description) > 2000 {
		return errors.New("Haasteen kuvaus voi olla 1-2000 merkkiä pitkä")
	}

	if c.Unit == "" {
		return errors.New("Haasteen yksikkö puuttuu")
	}

	if len(c.Unit) < 1 || len(c.Unit) > 255 {
		return errors.New("Haasteen yksikkö voi olla 1-255 merkkiä pitkä")
	}

	if c.GoalPoints < 1 {
		return errors.New("Haasteen tavoitepisteet tulee olla yli 1")
	}

	start, err := time.Parse(time.DateOnly, c.StartDate)
	if err != nil {
		return errors.New("Aloituspäivä on väärässä formaatissa")
	}

	end, err := time.Parse(time.DateOnly, c.EndDate)
	if err != nil {
		return errors.New("Lopetuspäivä on väärässä formaatissa")
	}

	if !end.After(start) {
		return errors.New("Aloituspäivän tulee olla aina ennen lopetuspäivää")
	}

	return nil
}

func (u *User) Validate() error {
	if u == nil {
		return errors.New("Virheellinen käyttäjä")
	}

	if u.Name == "" {
		return errors.New("Käyttäjän nimi puuttuu")
	}

	if len(u.Name) < 1 || len(u.Name) > 255 {
		return errors.New("Käyttäjän nimi voi olla 1-255 merkkiä pitkä")
	}

	if u.Email == "" {
		return errors.New("Käyttäjän sähköposti puuttuu")
	}

	if _, err := mail.ParseAddress(u.Email); err != nil {
		return errors.New("Sähköposti ei ole oikeassa muodossa")
	}

	return nil
}
