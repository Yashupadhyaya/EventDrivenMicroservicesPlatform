package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/Garvit-Jethwani/database-service/config"
	"github.com/Garvit-Jethwani/database-service/database"
	"github.com/Garvit-Jethwani/database-service/grpcserver"
	"github.com/Garvit-Jethwani/database-service/proto"
	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// Initialize PostgreSQL if configured
	if postgresURL := os.Getenv("POSTGRES_URL"); postgresURL != "" {
		if err := database.InitPostgres(postgresURL); err != nil {
			log.Fatalf("could not connect to PostgreSQL: %v", err)
		}
	}

	// Initialize MySQL if configured
	if mysqlURL := os.Getenv("MYSQL_URL"); mysqlURL != "" {
		if err := database.InitMySQL(mysqlURL); err != nil {
			log.Fatalf("could not connect to MySQL: %v", err)
		}
	}

	// Get the initialized DB
	db := database.GetDB()
	if db == nil {
		log.Fatalf("no database initialized")
	}

	// Create a gRPC server instance
	grpcServer := grpc.NewServer()

	// Register database service
	databaseServer := &grpcserver.Server{DB: db}
	proto.RegisterDatabaseServiceServer(grpcServer, databaseServer)

	// Start the gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Server is ready to handle requests at %s", lis.Addr().String())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
