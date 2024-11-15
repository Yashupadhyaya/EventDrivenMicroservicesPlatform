package models

import "time"

type Event struct {
	Type    string `json:"type"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Payload string `json:"payload"`
}

type EventStatus struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

const (
	StatusPending   = "PENDING"
	StatusProcessed = "PROCESSED"
	StatusFailed    = "FAILED"
)

type EventSummary struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Status    EventStatus `json:"status"`
	Timestamp time.Time   `json:"timestamp"`
}

type EventFilter struct {
	Type      string
	StartTime time.Time
	EndTime   time.Time
}
