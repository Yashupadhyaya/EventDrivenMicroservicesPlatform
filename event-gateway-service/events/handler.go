package events

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"log"

	"github.com/Yashupadhyaya/event-gateway-service/models"
	"github.com/gorilla/mux"
)

type EventHandler struct {
	eventStore EventStore
}

func NewEventHandler(eventStore EventStore) *EventHandler {
	return &EventHandler{
		eventStore: eventStore,
	}
}

func (h *EventHandler) IngestEvent(w http.ResponseWriter, r *http.Request) {
	var event models.Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	fmt.Println("publishing the event", event)
	eventID, err := h.eventStore.Publish(event)
	if err != nil {
		http.Error(w, "Failed to publish event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(eventID))
}

func (h *EventHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID := vars["eventId"]

	event, err := h.eventStore.Get(eventID)
	if err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

func (h *EventHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	offsetParam := query.Get("offset")
	limitParam := query.Get("limit")

	offset, err := strconv.Atoi(offsetParam)
	if err != nil {
		offset = 0
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 10 // Set a default limit
	}

	events, err := h.eventStore.List(offset, limit)
	if err != nil {
		log.Printf("Failed to list events: %v", err)
		http.Error(w, "Failed to list events", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func (h *EventHandler) GetEventStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID := vars["eventId"]

	status, err := h.eventStore.GetEventStatus(eventID)
	if err != nil {
		http.Error(w, "Failed to get event status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
