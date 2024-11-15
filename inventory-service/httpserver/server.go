package httpserver

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Yashupadhyaya/inventory-service/config"
	"github.com/Yashupadhyaya/inventory-service/database"
	"github.com/gorilla/mux"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config) *Server {
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/products/{productId}/inventory", getInventoryHandler).Methods("GET")

	httpServer := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: router,
	}

	return &Server{httpServer: httpServer}
}

func (s *Server) Start() error {
	// Initialize database
	if err := database.InitDatabase(os.Getenv("DATABASE_URL")); err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	// Start the server in a new goroutine
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("could not listen on %s: %v", s.httpServer.Addr, err)
		}
	}()
	log.Printf("Server is ready to handle requests at %s", s.httpServer.Addr)

	// Wait for interrupt signal to gracefully shutdown the server
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	log.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.httpServer.Shutdown(ctx)
}

func getInventoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["productId"]

	inventory, err := database.GetInventoryByProductID(productID)
	if err != nil {
		if err.Error() == "product not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inventory)
}
