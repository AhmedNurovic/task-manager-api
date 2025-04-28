package model

type User struct {
	ID       uint   `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	Password string `json:"-" db:"password"`
}