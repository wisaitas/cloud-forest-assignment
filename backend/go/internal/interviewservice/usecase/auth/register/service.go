package register

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/domain/repositorysql"
	"github.com/wisaitas/cloud-forest-assignment/pkg/httpx"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Service(c *fiber.Ctx, request *Request) error
}

type service struct {
	userRepository repositorysql.UserRepository
}

func newService(
	userRepository repositorysql.UserRepository,
) Service {
	return &service{
		userRepository: userRepository,
	}
}

func (s *service) Service(c *fiber.Ctx, request *Request) error {
	user := s.mapRequestToEntity(request)
	user.ID = uuid.New().String()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return httpx.NewErrorResponse[any](c, fiber.StatusInternalServerError, err)
	}

	user.Password = string(hashedPassword)

	if err := s.userRepository.CreateUser(c.Context(), *user); err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return httpx.NewErrorResponse[any](c, fiber.StatusConflict, err)
		}

		return httpx.NewErrorResponse[any](c, fiber.StatusInternalServerError, err)
	}

	return httpx.NewSuccessResponse[any](c, nil, fiber.StatusCreated, nil, nil)
}
