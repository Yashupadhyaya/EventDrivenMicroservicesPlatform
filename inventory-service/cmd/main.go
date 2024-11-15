package main

import (
	"log"

	"github.com/Garvit-Jethwani/inventory-service/config"
	"github.com/Garvit-Jethwani/inventory-service/events"
	"github.com/Garvit-Jethwani/inventory-service/httpserver"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// Start event consumer
	go events.StartConsumer(cfg)

	// Start HTTP server
	server := httpserver.NewServer(cfg)
	if err := server.Start(); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
