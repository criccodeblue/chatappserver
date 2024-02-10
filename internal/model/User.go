package model

type User struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Handle string `json:"handle"` // make it to always have @ at beginning
}

type AuthUser struct {
	*User
	Password     string `json:"password"`
	PasswordHash []byte
}
