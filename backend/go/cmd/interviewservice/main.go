package main

import (
	"github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/initial"

	// Import for swag (parse handler comments)
	_ "github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/usecase/auth/login"
	_ "github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/usecase/auth/register"
	_ "github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/usecase/servers/getservers"
	_ "github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/usecase/servers/power"
	_ "github.com/wisaitas/cloud-forest-assignment/internal/interviewservice/usecase/servers/provision"
)

// @title			Interview Service API
// @version		1.0
// @description	API สำหรับ Interview Service (Auth: register, login)
// @termsOfService	http://swagger.io/terms/

// @contact.name	API Support
// @contact.url	http://www.swagger.io/support
// @contact.email	support@swagger.io

// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:8080
// @BasePath	/api/v1
//
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
//
//go:generate sh -c "cd ../.. && go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/interviewservice/main.go -d . -o internal/interviewservice/docs --parseDependency --parseInternal"

func main() {
	app := initial.New()
	defer app.Shutdown()
	app.Run()
}
