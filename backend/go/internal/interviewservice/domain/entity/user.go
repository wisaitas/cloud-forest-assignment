package entity

import "github.com/golang-jwt/jwt/v5"

type User struct {
	ID       string
	Email    string
	Password string
}

type UserContext struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

func (u *UserContext) GetID() string {
	return u.ID
}

func (u *UserContext) GetEmail() string {
	return u.Email
}
