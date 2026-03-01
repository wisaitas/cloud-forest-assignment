package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/wisaitas/cloud-forest-assignment/internal/interviewservice"
	"github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/domain/entity"
	"github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/usecase/servers"
	"github.com/wisaitas/cloud-forest-assignment/pkg/httpx"
)

type ServersRouter struct {
	group         fiber.Router
	serverUseCase *servers.UseCase
}

func NewServersRouter(group fiber.Router, serverUseCase *servers.UseCase) *ServersRouter {
	return &ServersRouter{
		group:         group,
		serverUseCase: serverUseCase,
	}
}

func (r *ServersRouter) Setup() {
	serversGroup := r.group.Group("/servers", middleware)
	serversGroup.Get("/", r.serverUseCase.List.Handler)
	serversGroup.Post("/", r.serverUseCase.Provision.Handler)
	serversGroup.Post("/:id/power", r.serverUseCase.Power.Handler)
}

func middleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return httpx.NewErrorResponse[any](c, http.StatusUnauthorized, errors.New("authorization header is required"))
	}

	// if !strings.HasPrefix(authHeader, "Bearer ") {
	// 	return httpx.NewErrorResponse[any](c, http.StatusUnauthorized, errors.New("invalid token type"))
	// }

	// token := strings.TrimPrefix(authHeader, "Bearer ")

	var tokenContext entity.UserContext
	_, err := jwt.ParseWithClaims(authHeader, &tokenContext, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(interviewservice.Config.Jwt.Secret), nil
	})
	if err != nil {
		return httpx.NewErrorResponse[any](c, http.StatusUnauthorized, err)
	}

	userContextJSON, ok := interviewservice.MockInMemmory[fmt.Sprintf("access_token:%s", tokenContext.UserID)].(string)
	if !ok {
		return httpx.NewErrorResponse[any](c, http.StatusUnauthorized, errors.New("access token not found"))
	}

	userContext := entity.UserContext{}
	err = json.Unmarshal([]byte(userContextJSON), &userContext)
	if err != nil {
		return httpx.NewErrorResponse[any](c, http.StatusInternalServerError, err)
	}

	c.Locals("userContext", userContext)
	return c.Next()
}
