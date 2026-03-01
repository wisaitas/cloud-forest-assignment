package sqlx

import (
	"time"

	"gorm.io/gorm"
)

type Config struct {
	Host            string        `env:"HOST"`
	Port            string        `env:"PORT"`
	User            string        `env:"USER"`
	Password        string        `env:"PASSWORD"`
	DBName          string        `env:"DB_NAME"`
	SSLMode         string        `env:"SSL_MODE"`
	MaxIdleConns    int           `env:"MAX_IDLE_CONNS"`
	MaxOpenConns    int           `env:"MAX_OPEN_CONNS"`
	ConnMaxLifetime time.Duration `env:"CONN_MAX_LIFETIME"`
	gorm.Config
}
