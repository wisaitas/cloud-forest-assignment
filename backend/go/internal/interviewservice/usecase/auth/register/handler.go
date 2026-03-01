package register

import (
	"github.com/gofiber/fiber/v2"
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

// Handler registers a new user.
// @Summary		Register new user
// @Description	สร้างบัญชีผู้ใช้ใหม่
// @Tags		auth
// @Accept		json
// @Produce		json
// @Param		request	body		Request	true	"Register request"
// @Success	201		{object}	nil				"Created"
// @Failure	400	"Bad Request - validation error"
// @Failure	409	"Conflict - email already exists"
// @Failure	500	"Internal Server Error"
// @Router		/auth/register [post]
func (h *Handler) Handler(c *fiber.Ctx) error {
	req := Request{}
	if err := c.BodyParser(&req); err != nil {
		return httpx.NewErrorResponse[any](c, fiber.StatusBadRequest, err)
	}

	if err := h.validator.ValidateStruct(&req); err != nil {
		return httpx.NewErrorResponse[any](c, fiber.StatusBadRequest, err)
	}

	return h.service.Service(c, &req)
}
