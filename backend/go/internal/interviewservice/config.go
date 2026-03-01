package interviewservice

import (
	"time"
)

var Config struct {
	Service struct {
		Port           string `env:"PORT" envDefault:"8080"`
		Name           string `env:"NAME" envDefault:"interview-service"`
		AllowedOrigins string `env:"ALLOWED_ORIGINS" envDefault:"http://localhost:3000"`
	} `envPrefix:"SERVICE_"`

	InfraService struct {
		URL string `env:"URL" envDefault:"http://localhost:8081"`
	} `envPrefix:"INFRA_SERVICE_"`

	Jwt struct {
		AccessTTL  time.Duration `env:"ACCESS_TTL" envDefault:"1h"`
		RefreshTTL time.Duration `env:"REFRESH_TTL" envDefault:"24h"`
		Secret     string        `env:"SECRET" envDefault:"dev-secret-change-in-production"`
	} `envPrefix:"JWT_"`
}

var MockInMemmory = map[string]any{}
