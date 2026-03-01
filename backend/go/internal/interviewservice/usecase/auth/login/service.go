package login

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/wisaitas/cloud-forest-assignment/internal/interviewservice"
	"github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/domain/entity"
	"github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/domain/repositorysql"

	"github.com/wisaitas/cloud-forest-assignment/pkg/bcryptx"
	"github.com/wisaitas/cloud-forest-assignment/pkg/httpx"
	"github.com/wisaitas/cloud-forest-assignment/pkg/jwtx"
)

type Service interface {
	Service(c *fiber.Ctx, request *Request) error
}

type service struct {
	userRepository repositorysql.UserRepository
	bcrypt         bcryptx.Bcrypt
	jwt            jwtx.Jwt
}

func newService(
	userRepository repositorysql.UserRepository,
	bcrypt bcryptx.Bcrypt,
	jwt jwtx.Jwt,
) Service {
	return &service{
		userRepository: userRepository,
		bcrypt:         bcrypt,
		jwt:            jwt,
	}
}

func (s *service) Service(c *fiber.Ctx, req *Request) error {
	user, err := s.userRepository.GetByEmail(c.Context(), req.Email)
	if err != nil {
		if errors.Is(err, repositorysql.ErrNotFound) {
			return httpx.NewErrorResponse[any](c, http.StatusNotFound, err)
		}
		return httpx.NewErrorResponse[any](c, http.StatusInternalServerError, err)
	}

	if err := s.bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return httpx.NewErrorResponse[any](c, http.StatusUnauthorized, err)
	}

	accessClaims := entity.UserContext{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(interviewservice.Config.Jwt.AccessTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	accessToken, err := s.jwt.Generate(&accessClaims)
	if err != nil {
		return httpx.NewErrorResponse[any](c, http.StatusInternalServerError, err)
	}

	refreshClaims := entity.UserContext{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(interviewservice.Config.Jwt.RefreshTTL)),
		},
	}
	refreshToken, err := s.jwt.Generate(&refreshClaims)
	if err != nil {
		return httpx.NewErrorResponse[any](c, http.StatusInternalServerError, err)
	}

	sessionData := map[string]string{"user_id": user.ID, "email": user.Email}
	sessionDataJSON, err := json.Marshal(sessionData)
	if err != nil {
		return httpx.NewErrorResponse[any](c, http.StatusInternalServerError, err)
	}

	interviewservice.MockInMemmory[fmt.Sprintf("access_token:%s", user.ID)] = string(sessionDataJSON)
	interviewservice.MockInMemmory[fmt.Sprintf("refresh_token:%s", user.ID)] = string(sessionDataJSON)

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HTTPOnly: true,
		SameSite: "Lax",
		MaxAge:   int(interviewservice.Config.Jwt.AccessTTL.Seconds()),
	})

	resp := Response{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return httpx.NewSuccessResponse(c, &resp, http.StatusOK, nil, nil)
}
