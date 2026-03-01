package login

import (
	"github.com/gofiber/fiber/v2"
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

// Handler authenticates user and returns tokens.
// @Summary		Login
// @Description	เข้าสู่ระบบด้วย email และ password ได้รับ access_token และ refresh_token
// @Tags		auth
// @Accept		json
// @Produce		json
// @Param		request	body		Request	true	"Login request"
// @Success	200		{object}	Response	"OK - returns access_token and refresh_token in data"
// @Failure	400	"Bad Request"
// @Failure	401	"Unauthorized - wrong password"
// @Failure	404	"Not Found - user not found"
// @Failure	500	"Internal Server Error"
// @Router		/auth/login [post]
func (h *Handler) Handler(c *fiber.Ctx) error {
	req := Request{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := h.validator.ValidateStruct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return h.service.Service(c, &req)
}
