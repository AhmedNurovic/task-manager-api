package model

type Task struct {
	ID     uint   `json:"id" db:"id"`
	UserID uint   `json:"user_id" db:"user_id"`
	Title  string `json:"title" db:"title"`
	Status string `json:"status" db:"status"`
}