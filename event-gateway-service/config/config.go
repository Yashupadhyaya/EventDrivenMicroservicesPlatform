package config

import (
	"errors"
	"os"
	"strings"
)

type Config struct {
	HTTPAddress   string
	EventStore    string
	KafkaBrokers  []string
	KafkaTopic    string
	NATSURL       string
	NATSClusterID string
}

func LoadConfig() (*Config, error) {
	httpAddress := os.Getenv("HTTP_ADDRESS")
	eventStore := os.Getenv("EVENT_STORE")
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	natsURL := os.Getenv("NATS_URL")
	natsClusterID := os.Getenv("NATS_CLUSTER_ID")

	// Validate and parse configuration values
	if httpAddress == "" || eventStore == "" {
		return nil, errors.New("required environment variables HTTP_ADDRESS or EVENT_STORE are missing")
	}

	return &Config{
		HTTPAddress:   httpAddress,
		EventStore:    eventStore,
		KafkaBrokers:  parseKafkaBrokers(kafkaBrokers),
		KafkaTopic:    kafkaTopic,
		NATSURL:       natsURL,
		NATSClusterID: natsClusterID,
	}, nil
}

// parseKafkaBrokers splits a comma-separated brokers list into a slice of strings
func parseKafkaBrokers(brokersList string) []string {
	if brokersList == "" {
		return nil
	}
	return strings.Split(brokersList, ",")
}
