package models

type User struct {
	ID       string
	Name     string
	Email    string
	Password string // Need to hash this!
}
