package initial

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/wisaitas/cloud-forest-assignment/internal/interviewservice"
	"github.com/wisaitas/cloud-forest-assignment/pkg/httpx"
)

func newMiddleware(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins:     interviewservice.Config.Service.AllowedOrigins,
		AllowCredentials: true,
	}))
	app.Use(httpx.NewLogger(interviewservice.Config.Service.Name))

}
