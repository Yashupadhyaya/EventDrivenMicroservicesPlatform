package main

import (
	"log"

	"github.com/Yashupadhyaya/notification-service/config"
	"github.com/Yashupadhyaya/notification-service/events"
	"github.com/Yashupadhyaya/notification-service/notification"
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
