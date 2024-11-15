package httpserver

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Yashupadhyaya/order-service/config"
	"github.com/Yashupadhyaya/order-service/database"
	"github.com/Yashupadhyaya/order-service/events"
	"github.com/Yashupadhyaya/order-service/models"
	"github.com/gorilla/mux"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config) *Server {
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/orders", createOrderHandler).Methods("POST")
	router.HandleFunc("/orders/{orderId}", getOrderHandler).Methods("GET")

	httpServer := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: router,
	}

	return &Server{httpServer: httpServer}
}

func (s *Server) Start() error {
	// Initialize database
	// if err := database.InitDatabase(os.Getenv("DATABASE_URL")); err != nil {
	// 	log.Fatalf("could not connect to database: %v", err)
	// }

	// Initialize event producer
	//	events.InitEventProducer([]string{os.Getenv("EVENT_STORE_URL")})

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

func createOrderHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("createOrderHandler")
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(order)
	if err := database.CreateOrder(&order); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// if err := events.ProduceOrderCreatedEvent(&order); err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	events.HandleOrderCreatedEvent(order.ID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func getOrderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderId := vars["orderId"]

	order, err := database.GetOrderById(orderId)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "order not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
