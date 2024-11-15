package config

import (
	"os"
	"strings"
)

type Config struct {
	EventStore   string
	KafkaBrokers []string
	KafkaTopic   string
	HTTPAddress  string
}

// should we be calling /events microservice for updates
func LoadConfig() (*Config, error) {
	eventStore := os.Getenv("EVENT_STORE")
	kafkaBrokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	httpAddress := os.Getenv("HTTP_ADDRESS")

	return &Config{
		EventStore:   eventStore,
		KafkaBrokers: kafkaBrokers,
		KafkaTopic:   kafkaTopic,
		HTTPAddress:  httpAddress,
	}, nil
}
