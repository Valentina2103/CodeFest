package models

// User is a struct that defines the user model for login
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

