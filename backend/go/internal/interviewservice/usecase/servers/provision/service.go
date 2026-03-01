package provision

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
	Service(c *fiber.Ctx, userContext entity.UserContext, request *Request) error
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

func (s *service) Service(c *fiber.Ctx, userContext entity.UserContext, req *Request) error {
	valid, err := s.infraClient.IsValidSKU(c.Context(), req.SKU)
	if err != nil {
		return httpx.NewErrorResponse[any](c, http.StatusBadGateway, err)
	}
	if !valid {
		return httpx.NewErrorResponse[any](c, http.StatusBadRequest, errors.New("invalid sku"))
	}
	provisioned, err := s.infraClient.Provision(c.Context(), req.SKU)
	if err != nil {
		return httpx.NewErrorResponse[any](c, http.StatusBadGateway, err)
	}
	powerStatus := "on"
	if provisioned.Status == "stopped" {
		powerStatus = "off"
	}
	server := entity.Server{
		ID:                       uuid.New().String(),
		UserID:                   userContext.UserID,
		InfrastructureResourceID: provisioned.ID,
		SKU:                      req.SKU,
		PowerStatus:              powerStatus,
		CreatedAt:                time.Now().UTC(),
	}
	if err := s.serverRepo.Create(c.Context(), server); err != nil {
		return httpx.NewErrorResponse[any](c, http.StatusInternalServerError, err)
	}
	_ = s.activityLogRepo.Append(c.Context(), entity.ActivityLog{
		ID:           uuid.New().String(),
		UserID:       userContext.UserID,
		Action:       "provision_server",
		ResourceType: "server",
		ResourceID:   server.ID,
		Details:      req.SKU,
		CreatedAt:    time.Now().UTC(),
	})
	resp := Response{Success: true, ID: provisioned.ID}
	return httpx.NewSuccessResponse(c, &resp, http.StatusOK, nil, nil)
}
