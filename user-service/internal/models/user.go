package models

type UserProfile struct {
	UserID int64  `json:"user_id"`
	Name   string `json:"name"`
	Phone  string `json:"phone"`
}
type Order struct {
	ID          int64
	UserID      int64
	ServiceID   int64
	Description string
	Area        float64
	Status      string
	CreatedAt   int64
}

type Favorite struct {
	UserID    int64
	ServiceID int64
}
