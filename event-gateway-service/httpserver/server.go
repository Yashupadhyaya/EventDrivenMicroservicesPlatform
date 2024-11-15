package httpserver

import (
	"net/http"

	"github.com/Garvit-Jethwani/event-gateway-service/events"
	"github.com/gorilla/mux"
)

type Server struct {
	addr    string
	handler *events.EventHandler
}

func NewServer(addr string, handler *events.EventHandler) *Server {
	return &Server{
		addr:    addr,
		handler: handler,
	}
}

func (s *Server) Start() error {
	r := mux.NewRouter()
	r.HandleFunc("/events", s.handler.IngestEvent).Methods("POST")
	r.HandleFunc("/events/{eventId}", s.handler.GetEvent).Methods("GET")
	r.HandleFunc("/events", s.handler.ListEvents).Methods("GET")
	r.HandleFunc("/events/{eventId}/status", s.handler.GetEventStatus).Methods("GET") // Add the status API

	http.Handle("/", r)
	return http.ListenAndServe(s.addr, nil)
}
