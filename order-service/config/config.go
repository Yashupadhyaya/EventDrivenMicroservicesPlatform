package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	HTTPPort         string        `envconfig:"HTTP_PORT" default:"8081"`
	OrderServiceHost string        `envconfig:"DATABASE_SERVICE_HOST" default:"localhost"`
	OrderServicePort int           `envconfig:"DATABASE_SERVICE_PORT" default:"9093"`
	JWTSecret        string        `envconfig:"JWT_SECRET" required:"true"`
	ShutdownTimeout  time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"10s"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	return &cfg, err
}
