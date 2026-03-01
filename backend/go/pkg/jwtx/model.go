package jwtx

import jwtLib "github.com/golang-jwt/jwt/v5"

type Claims interface {
	jwtLib.Claims
	GetID() string
}

type StandardClaims struct {
	jwtLib.RegisteredClaims
	ID string `json:"id"`
}

func (s StandardClaims) GetID() string {
	return s.ID
}
