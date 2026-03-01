package initial

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v11"
	"github.com/wisaitas/cloud-forest-assignment/internal/interviewservice"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	fiberSwagger "github.com/swaggo/fiber-swagger"
	_ "github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/docs"
)

func init() {
	for _, path := range []string{".env", "../.env", "../../.env"} {
		if err := godotenv.Load(path); err == nil {
			break
		}
	}

	if err := env.Parse(&interviewservice.Config); err != nil {
		log.Fatalln(err)
	}
}

type App struct {
	FiberApp *fiber.App
	config   *config
}

func New() *App {
	config := newConfig()
	sdk := newSDK()
	repository := newRepository(sdk)
	useCase := newUseCase(config, repository, sdk)
	app := fiber.New()
	newMiddleware(app)

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	newRouter(app, useCase, sdk)

	return &App{
		FiberApp: app,
		config:   config,
	}
}

func (a *App) Run() {
	go func() {
		if err := a.FiberApp.Listen(":" + interviewservice.Config.Service.Port); err != nil {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func (a *App) Shutdown() {
	fmt.Println("Shutting down...")
}
