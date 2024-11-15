package kafka

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/Garvit-Jethwani/event-gateway-service/models"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaStore struct {
	producer *kafka.Producer
	consumer *kafka.Consumer
	topic    string
}

func NewKafkaStore(brokers []string, topic string) (*KafkaStore, error) {
	brokerList := strings.Join(brokers, ",")

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokerList,
	})
	if err != nil {
		return nil, err
	}

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":    brokerList,
		"group.id":             "event-consumer-group",
		"auto.offset.reset":    "earliest",
		"enable.auto.commit":   false, // Disable automatic offset commit
		"session.timeout.ms":   6000,
		"enable.partition.eof": true,
	})
	if err != nil {
		return nil, err
	}

	err = consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return nil, err
	}

	return &KafkaStore{
		producer: producer,
		consumer: consumer,
		topic:    topic,
	}, nil
}

func (k *KafkaStore) Publish(event models.Event) (string, error) {
	value, err := json.Marshal(event)
	if err != nil {
		return "", err
	}

	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &k.topic, Partition: kafka.PartitionAny},
		Value:          value,
	}

	err = k.producer.Produce(msg, nil)
	if err != nil {
		return "", err
	}

	return string(msg.Key), nil
}

func (k *KafkaStore) Get(eventID string) (models.Event, error) {
	for {
		msg, err := k.consumer.ReadMessage(10 * time.Second)
		if err != nil {
			if err.(kafka.Error).Code() == kafka.ErrTimedOut {
				log.Printf("Error reading message: Timed out while reading message with ID %s", eventID)
				break
			}
			return models.Event{}, err
		}

		var event models.Event
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			log.Printf("Error unmarshalling event: %v", err)
			continue
		}

		if event.ID == eventID {
			// Commit the offset for this message after processing
			_, commitErr := k.consumer.CommitMessage(msg)
			if commitErr != nil {
				log.Printf("Failed to commit message: %v", commitErr)
				return event, commitErr
			}

			return event, nil
		}
	}
	return models.Event{}, errors.New("event not found")
}

func (k *KafkaStore) List(offset, limit int) ([]models.Event, error) {
	var events []models.Event
	var count int

	for {
		msg, err := k.consumer.ReadMessage(10 * time.Second)
		if err != nil {
			if err.(kafka.Error).Code() == kafka.ErrTimedOut {
				log.Printf("Error reading message: Local: Timed out")
				break
			}
			return nil, err
		}

		if count < offset {
			count++
			continue
		}

		var event models.Event
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			log.Printf("Error unmarshalling event: %v", err)
			continue
		}

		events = append(events, event)
		count++

		// Commit the offset for the message after processing
		_, commitErr := k.consumer.CommitMessage(msg)
		if commitErr != nil {
			log.Printf("Failed to commit message: %v", commitErr)
		}

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

func (k *KafkaStore) GetEventStatus(eventID string) (models.EventStatus, error) {
	for {
		msg, err := k.consumer.ReadMessage(10 * time.Second)
		if err != nil {
			if err.(kafka.Error).Code() == kafka.ErrTimedOut {
				log.Printf("Error reading message: Timed out while reading message with ID %s", eventID)
				break
			}
			return models.EventStatus{}, err
		}

		var event models.Event
		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			log.Printf("Error unmarshalling event: %v", err)
			continue
		}

		if event.ID == eventID {
			// Commit the offset for the message after processing
			_, commitErr := k.consumer.CommitMessage(msg)
			if commitErr != nil {
				log.Printf("Failed to commit message: %v", commitErr)
				return models.EventStatus{}, commitErr
			}

			return models.EventStatus{
				ID:     eventID,
				Status: models.StatusProcessed,
			}, nil
		}
	}
	return models.EventStatus{}, errors.New("event not found")
}

func (k *KafkaStore) Close() {
	k.producer.Close()
	k.consumer.Close()
}
