package main

import (
	"log"

	"github.com/Yashupadhyaya/event-gateway-service/config"
	"github.com/Yashupadhyaya/event-gateway-service/events"
	"github.com/Yashupadhyaya/event-gateway-service/httpserver"
	"github.com/Yashupadhyaya/event-gateway-service/store/kafka"
	"github.com/Yashupadhyaya/event-gateway-service/store/nats"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize Event Store client
	var eventStore events.EventStore
	switch cfg.EventStore {
	case "kafka":
		eventStore, err = kafka.NewKafkaStore(cfg.KafkaBrokers, cfg.KafkaTopic)
		if err != nil {
			log.Fatalf("Failed to create Kafka event store: %v", err)
		}
	case "nats":
		eventStore, err = nats.NewNATSStore(cfg.NATSURL, cfg.NATSClusterID)
		if err != nil {
			log.Fatalf("Failed to create NATS event store: %v", err)
		}
	default:
		log.Fatalf("Invalid event store configuration: %s", cfg.EventStore)
	}
	defer eventStore.Close()

	// Create event handler
	eventHandler := events.NewEventHandler(eventStore)

	// Create HTTP server
	server := httpserver.NewServer(cfg.HTTPAddress, eventHandler)

	// Start HTTP server
	log.Printf("Starting Event Gateway Service on %s", cfg.HTTPAddress)
	err = server.Start()
	if err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
