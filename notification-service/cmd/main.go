package main

import (
	"log"

	"github.com/Garvit-Jethwani/notification-service/config"
	"github.com/Garvit-Jethwani/notification-service/events"
	"github.com/Garvit-Jethwani/notification-service/notification"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	consumer, err := events.NewKafkaConsumer(cfg.KafkaBrokers, "notification-service-group", cfg.KafkaTopic)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()

	// Create notification service
	notificationService := notification.NewNotificationService(consumer)

	// Start the notification service
	err = notificationService.Start()
	if err != nil {
		log.Fatalf("Failed to start notification service: %v", err)
	}
}
