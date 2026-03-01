package power

import (
	"errors"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/caller/infraservice"
	"github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/domain/entity"
	"github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/domain/repositorysql"
	"github.com/wisaitas/cloud-forest-assignment/pkg/httpx"
)

type Service interface {
	Service(c *fiber.Ctx, userContext entity.UserContext, serverID string, request *Request) error
}

type service struct {
	serverRepo      repositorysql.ServerRepository
	activityLogRepo repositorysql.ActivityLogRepository
	infraClient     infraservice.InfraServiceCaller
}

func newService(
	serverRepo repositorysql.ServerRepository,
	activityLogRepo repositorysql.ActivityLogRepository,
	infraClient infraservice.InfraServiceCaller,
) Service {
	return &service{
		serverRepo:      serverRepo,
		activityLogRepo: activityLogRepo,
		infraClient:     infraClient,
	}
}

func (s *service) Service(c *fiber.Ctx, userContext entity.UserContext, serverID string, req *Request) error {
	server, err := s.serverRepo.GetByID(c.Context(), serverID)
	if err != nil {
		return httpx.NewErrorResponse[any](c, http.StatusNotFound, err)
	}

	if server.UserID != userContext.UserID {
		return httpx.NewErrorResponse[any](c, http.StatusNotFound, errors.New("server not found"))
	}
	result, err := s.infraClient.Power(c.Context(), server.InfrastructureResourceID, req.Action)
	if err != nil {
		return httpx.NewErrorResponse[any](c, http.StatusBadGateway, err)
	}
	if err := s.serverRepo.UpdatePowerStatus(c.Context(), serverID, result.State); err != nil {
		return httpx.NewErrorResponse[any](c, http.StatusInternalServerError, err)
	}
	_ = s.activityLogRepo.Append(c.Context(), entity.ActivityLog{
		ID:           uuid.New().String(),
		UserID:       userContext.UserID,
		Action:       "power_" + req.Action,
		ResourceType: "server",
		ResourceID:   serverID,
		Details:      result.State,
		CreatedAt:    time.Now().UTC(),
	})
	resp := Response{Success: true, State: result.State}
	return httpx.NewSuccessResponse(c, &resp, http.StatusOK, nil, nil)
}
