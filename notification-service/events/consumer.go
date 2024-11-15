package events

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/Garvit-Jethwani/notification-service/models"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// EventConsumer is an interface for consuming events from an event store.
type EventConsumer interface {
	Consume() (<-chan models.Event, error)
	Close() error
}

type KafkaConsumer struct {
	consumer *kafka.Consumer
	topic    string
}

func NewKafkaConsumer(brokers []string, groupID, topic string) (*KafkaConsumer, error) {
	brokersList := strings.Join(brokers, ",")
	configMap := &kafka.ConfigMap{
		"bootstrap.servers": brokersList,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	}

	c, err := kafka.NewConsumer(configMap)
	if err != nil {
		return nil, err
	}

	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		// Close the consumer if subscription fails
		c.Close()
		return nil, err
	}

	return &KafkaConsumer{
		consumer: c,
		topic:    topic,
	}, nil
}
func (k *KafkaConsumer) Consume() (<-chan models.Event, error) {
	eventChan := make(chan models.Event)

	// Check if k.consumer is nil before starting the goroutine
	if k.consumer == nil {
		return nil, fmt.Errorf("kafka consumer is not initialized")
	}

	go func() {
		for {
			ev, err := k.consumer.ReadMessage(-1)
			if err != nil {
				// Handle error
				log.Printf("Error reading message from Kafka: %v", err)
				continue
			}

			var event models.Event
			err = json.Unmarshal(ev.Value, &event)
			if err != nil {
				log.Printf("Error unmarshalling event: %v", err)
				continue
			}

			eventChan <- event
		}
	}()

	return eventChan, nil
}

func (k *KafkaConsumer) Close() error {
	if k.consumer == nil {
		return nil // Return without error if consumer is nil
	}
	return k.consumer.Close()
}
