package getservers

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/domain/entity"
	"github.com/wisaitas/cloud-forest-assignment/pkg/httpx"
	"github.com/wisaitas/cloud-forest-assignment/pkg/validatorx"
)

type Handler struct {
	service   Service
	validator validatorx.Validator
}

func newHandler(
	service Service,
	validator validatorx.Validator,
) *Handler {
	return &Handler{
		service:   service,
		validator: validator,
	}
}

// Handler lists all servers for the authenticated user.
// @Summary		List servers
// @Description	รายการ servers ทั้งหมดของ user ที่ login (ต้องส่ง JWT ใน cookie หรือ header)
// @Tags		servers
// @Accept		json
// @Produce		json
// @Security	BearerAuth
// @Success	200	{object}	Response	"OK - returns servers in data"
// @Failure	401	"Unauthorized"
// @Failure	500	"Internal Server Error"
// @Router		/servers [get]
func (h *Handler) Handler(c *fiber.Ctx) error {
	userContext, ok := c.Locals("userContext").(entity.UserContext)
	if !ok {
		return httpx.NewErrorResponse[any](c, http.StatusUnauthorized, errors.New("unauthorized"))
	}
	req := Request{UserID: userContext.UserID}
	return h.service.Service(c, &req)
}
