package events

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Event struct {
	Type    string `json:"type"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Payload string `json:"payload"`
}

func HandleOrderCreatedEvent(orderID string) {
	// Assuming the event gateway microservice is running at eventGatewayURL
	eventGatewayURL := "http://localhost:8080/events"

	event := Event{
		Type:    "order_created",
		ID:      orderID,
		Name:    "Order Created",
		Payload: orderID,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		// Handle error
		return
	}

	resp, err := http.Post(eventGatewayURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		// Handle error
		return
	}
	defer resp.Body.Close()

	// Handle response
}
