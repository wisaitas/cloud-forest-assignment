package redisx

type Config struct {
	Host     string `env:"HOST"`
	Port     string `env:"PORT"`
	Password string `env:"PASSWORD"`
	DB       int    `env:"DB"`
}
