package jwtx

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	jwtLib "github.com/golang-jwt/jwt/v5"
)

type Jwt interface {
	Generate(claims Claims) (string, error)
	Parse(tokenString string, claims jwtLib.Claims) (jwtLib.Claims, error)
	ExtractTokenFromHeader(c *fiber.Ctx) (string, error)
	ValidateToken(tokenString string, claims jwtLib.Claims) error
}

type jwtx struct {
	secret string
}

func NewJwt(secret string) Jwt {
	return &jwtx{
		secret: secret,
	}
}

func (j *jwtx) Generate(claims Claims) (string, error) {
	token := jwtLib.NewWithClaims(jwtLib.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", fmt.Errorf("[jwtx] : %w", err)
	}

	return tokenString, nil
}

func (j *jwtx) ExtractTokenFromHeader(c *fiber.Ctx) (string, error) {
	authHeader := c.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("[jwtx] : invalid token type")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	return token, nil
}

func (j *jwtx) Parse(tokenString string, claims jwtLib.Claims) (jwtLib.Claims, error) {
	_, err := jwtLib.ParseWithClaims(tokenString, claims, func(token *jwtLib.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtLib.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("[jwtx] : unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secret), nil
	})

	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (j *jwtx) ValidateToken(tokenString string, claims jwtLib.Claims) error {
	_, err := j.Parse(tokenString, claims)
	if err != nil {
		return fmt.Errorf("[jwtx] : %w", err)
	}

	return nil
}

func (j *jwtx) CreateStandardClaims(id string, expireTime time.Duration) StandardClaims {
	return StandardClaims{
		ID: id,
		RegisteredClaims: jwtLib.RegisteredClaims{
			ExpiresAt: jwtLib.NewNumericDate(time.Now().Add(expireTime)),
			IssuedAt:  jwtLib.NewNumericDate(time.Now()),
		},
	}
}
