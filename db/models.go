package db

type Challenge struct {
	ID           int            `json:"id"`
	Title        string         `json:"title"`
	GoalPoints   int            `json:"goal_points"`
	StartDate    string         `json:"start_date"`
	EndDate      string         `json:"end_date"`
	Pot          int            `json:"pot"`
	Performances []*Performance `json:"performances"`
}

type User struct {
	ID            int            `json:"id"`
	Name          string         `json:"name"`
	Email 		  string 		 `json:"email"`
	PasswordHash  string         `json:"password_hash,omitempty"`
	PasswordPlain string         `json:"password_plain,omitempty"`
	IsAdmin       bool           `json:"is_admin,omitempty"`
	Performances  []*Performance `json:"performances"`
}

type Performance struct {
	ID          int    `json:"id"`
	Points      int    `json:"points"`
	CreatedAt   string `json:"created_at"`
	UserID      int    `json:"user_id"`
	ChallengeID int    `json:"challenge_id"`
	User        *User  `json:"user"`
}