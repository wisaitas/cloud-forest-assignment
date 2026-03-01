package getservers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/domain/repositorysql"
	"github.com/wisaitas/cloud-forest-assignment/pkg/httpx"
)

type Service interface {
	Service(c *fiber.Ctx, request *Request) error
}

type service struct {
	serverRepo repositorysql.ServerRepository
}

func newService(
	serverRepo repositorysql.ServerRepository,
) Service {
	return &service{
		serverRepo: serverRepo,
	}
}

func (s *service) Service(c *fiber.Ctx, req *Request) error {
	list, err := s.serverRepo.ListByUserID(c.Context(), req.UserID)
	if err != nil {
		return httpx.NewErrorResponse[any](c, http.StatusInternalServerError, err)
	}
	items := make([]ServerItem, 0, len(list))
	for _, sv := range list {
		items = append(items, ServerItem{
			ID:                       sv.ID,
			InfrastructureResourceID: sv.InfrastructureResourceID,
			SKU:                      sv.SKU,
			PowerStatus:              sv.PowerStatus,
		})
	}
	resp := Response{Servers: items}
	return httpx.NewSuccessResponse(c, &resp, http.StatusOK, nil, nil)
}
