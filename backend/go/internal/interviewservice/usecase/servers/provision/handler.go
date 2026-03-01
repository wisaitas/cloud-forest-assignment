package provision

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

// Handler provisions a new server. Validates SKU with infra service then creates server.
// @Summary		Provision server
// @Description	สร้าง server ใหม่ ต้อง validate SKU กับ infra service (ต้องส่ง JWT ใน cookie หรือ header)
// @Tags		servers
// @Accept		json
// @Produce		json
// @Security	BearerAuth
// @Param		request	body		Request	true	"Provision request (sku)"
// @Success	200		{object}	Response	"OK - success and id is Infrastructure Resource ID"
// @Failure	400	"Bad Request - invalid or unsupported sku"
// @Failure	401	"Unauthorized"
// @Failure	502	"Bad Gateway - infra service error"
// @Failure	500	"Internal Server Error"
// @Router		/servers [post]
func (h *Handler) Handler(c *fiber.Ctx) error {
	userContext, ok := c.Locals("userContext").(entity.UserContext)
	if !ok {
		return httpx.NewErrorResponse[any](c, http.StatusUnauthorized, errors.New("unauthorized"))
	}

	req := Request{}
	if err := c.BodyParser(&req); err != nil {
		return httpx.NewErrorResponse[any](c, http.StatusBadRequest, errors.New("invalid request body"))
	}

	if err := h.validator.ValidateStruct(&req); err != nil {
		return httpx.NewErrorResponse[any](c, http.StatusBadRequest, err)
	}

	return h.service.Service(c, userContext, &req)
}
