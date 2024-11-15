package main

import (
	"fmt"
	"log"

	"github.com/Garvit-Jethwani/order-service/config"
	"github.com/Garvit-Jethwani/order-service/database"
	"github.com/Garvit-Jethwani/order-service/httpserver"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Use the host and port from the configuration
	addr := fmt.Sprintf("%s:%d", cfg.OrderServiceHost, cfg.OrderServicePort)

	if err := database.InitDatabase(addr); err != nil {
		log.Fatalf("Failed to initialize database client: %v", err)
	}

	server := httpserver.NewServer(cfg)
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
