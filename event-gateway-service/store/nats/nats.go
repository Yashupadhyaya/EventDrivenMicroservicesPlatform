package nats

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/Yashupadhyaya/event-gateway-service/models"
	"github.com/nats-io/nats.go"
)

type NATSStore struct {
	conn *nats.Conn
}

func NewNATSStore(url, clusterID string) (*NATSStore, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	return &NATSStore{
		conn: conn,
	}, nil
}

func (n *NATSStore) Publish(event models.Event) (string, error) {
	payload, err := json.Marshal(event)
	if err != nil {
		return "", err
	}

	msg, err := n.conn.Request("events", payload, 10*time.Second)
	if err != nil {
		return "", err
	}

	return string(msg.Data), nil
}

func (n *NATSStore) Get(eventID string) (models.Event, error) {
	sub, err := n.conn.SubscribeSync("events")
	if err != nil {
		return models.Event{}, err
	}
	defer sub.Unsubscribe()

	for {
		msg, err := sub.NextMsg(10 * time.Second)
		if err != nil {
			if err == nats.ErrTimeout {
				log.Printf("Error reading message: Timed out while reading message with ID %s", eventID)
				break
			}
			return models.Event{}, err
		}

		var event models.Event
		err = json.Unmarshal(msg.Data, &event)
		if err != nil {
			log.Printf("Error unmarshalling event: %v", err)
			continue
		}

		if event.ID == eventID {
			return event, nil
		}
	}
	return models.Event{}, errors.New("event not found")
}

func (n *NATSStore) List(offset, limit int) ([]models.Event, error) {
	var events []models.Event
	var count int
	sub, err := n.conn.SubscribeSync("events")
	if err != nil {
		log.Printf("Error subscribing to sync: %v", err)
		return nil, err
	}
	defer sub.Unsubscribe()

	for {
		msg, err := sub.NextMsg(10 * time.Second)
		if err != nil {
			log.Printf("Error getting next message: %v", err)
			break
		}

		if count < offset {
			count++
			continue
		}

		var event models.Event
		err = json.Unmarshal(msg.Data, &event)
		if err != nil {
			log.Printf("Error unmarshalling event: %v", err)
			continue
		}

		events = append(events, event)
		count++

		if count >= offset+limit {
			break
		}
	}

	if len(events) == 0 {
		log.Println("No events found")
		return nil, errors.New("no events found")
	}

	return events, nil
}

func (n *NATSStore) GetEventStatus(eventID string) (models.EventStatus, error) {
	sub, err := n.conn.SubscribeSync("events")
	if err != nil {
		log.Printf("Error subscribing to sync: %v", err)
		return models.EventStatus{}, err
	}
	defer sub.Unsubscribe()

	for {
		msg, err := sub.NextMsg(10 * time.Second)
		if err != nil {
			log.Printf("Error getting next message: %v", err)
			return models.EventStatus{}, err
		}

		var event models.Event
		err = json.Unmarshal(msg.Data, &event)
		if err != nil {
			log.Printf("Error unmarshalling event: %v", err)
			continue
		}

		if event.ID == eventID {
			return models.EventStatus{
				ID:     eventID,
				Status: models.StatusProcessed,
			}, nil
		}
	}
	return models.EventStatus{}, errors.New("event not found")
}

func (n *NATSStore) Close() {
	n.conn.Close()
}
