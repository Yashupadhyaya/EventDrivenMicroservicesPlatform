
package config

import (
    "time"

    "github.com/kelseyhightower/envconfig"
)

type Config struct {
    HTTPPort       string        `envconfig:"HTTP_PORT" default:"8081"`
    GRPCPort       string        `envconfig:"GRPC_PORT" default:"8082"`
    DatabaseURL    string        `envconfig:"DATABASE_URL" required:"true"`
    JWTSecret      string        `envconfig:"JWT_SECRET" required:"true"`
    ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"10s"`
}

func LoadConfig() (*Config, error) {
    var cfg Config
    err := envconfig.Process("", &cfg)
    return &cfg, err
}
