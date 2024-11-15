package events

import "github.com/Garvit-Jethwani/event-gateway-service/models"

type EventStore interface {
	Publish(event models.Event) (string, error)
	Get(eventID string) (models.Event, error)
	List(offset, limit int) ([]models.Event, error)
	GetEventStatus(eventID string) (models.EventStatus, error)
	Close()
}
