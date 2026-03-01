package initial

import (
	"github.com/gofiber/fiber/v2"
	appRouter "github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/router"
)

type router struct {
	authRouter    *appRouter.AuthRouter
	serversRouter *appRouter.ServersRouter
}

func newRouter(
	app *fiber.App,
	useCase *useCase,
	sdk *sdk,
) {
	apiRouter := app.Group("/api/v1")

	r := &router{
		authRouter:    appRouter.NewAuthRouter(apiRouter, useCase.authUseCase),
		serversRouter: appRouter.NewServersRouter(apiRouter, useCase.serversUseCase),
	}
	r.setup()
}

func (r *router) setup() {
	r.authRouter.Setup()
	r.serversRouter.Setup()
}
