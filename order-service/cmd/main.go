package main

import (
	"fmt"
	"log"

	"github.com/Yashupadhyaya/order-service/config"
	"github.com/Yashupadhyaya/order-service/database"
	"github.com/Yashupadhyaya/order-service/httpserver"
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
