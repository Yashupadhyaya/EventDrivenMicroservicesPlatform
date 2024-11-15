
package config

import (
    "time"

    "github.com/kelseyhightower/envconfig"
)

type Config struct {
    HTTPPort       string        `envconfig:"HTTP_PORT" default:"8083"`
    DatabaseURL    string        `envconfig:"DATABASE_URL" required:"true"`
    EventStoreURL  string        `envconfig:"EVENT_STORE_URL" required:"true"`
    ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"10s"`
}

func LoadConfig() (*Config, error) {
    var cfg Config
    err := envconfig.Process("", &cfg)
    return &cfg, err
}
