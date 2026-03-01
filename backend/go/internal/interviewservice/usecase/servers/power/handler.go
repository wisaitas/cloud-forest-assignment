package power

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

// Handler changes power state of a server (on/off). Calls infra service.
// @Summary		Power on/off server
// @Description	เปิดหรือปิด server ต้องส่ง JWT ใน cookie หรือ header
// @Tags		servers
// @Accept		json
// @Produce		json
// @Security	BearerAuth
// @Param		server-id	path		string	true	"Server ID (API server id)"
// @Param		request		body		Request	true	"Power action: on or off"
// @Success	200		{object}	Response	"OK - success and current state"
// @Failure	400	"Bad Request - invalid body or server-id"
// @Failure	401	"Unauthorized"
// @Failure	404	"Not Found - server not found"
// @Failure	502	"Bad Gateway - infra service error"
// @Failure	500	"Internal Server Error"
// @Router		/servers/{server-id}/power [post]
func (h *Handler) Handler(c *fiber.Ctx) error {
	userContext, ok := c.Locals("userContext").(entity.UserContext)
	if !ok {
		return httpx.NewErrorResponse[any](c, http.StatusUnauthorized, errors.New("unauthorized"))
	}

	serverID := c.Params("id")
	if serverID == "" {
		return httpx.NewErrorResponse[any](c, http.StatusBadRequest, errors.New("server-id is required"))
	}

	req := Request{}
	if err := c.BodyParser(&req); err != nil {
		return httpx.NewErrorResponse[any](c, http.StatusBadRequest, errors.New("invalid request body"))
	}

	if err := h.validator.ValidateStruct(&req); err != nil {
		return httpx.NewErrorResponse[any](c, http.StatusBadRequest, err)
	}

	return h.service.Service(c, userContext, serverID, &req)
}
