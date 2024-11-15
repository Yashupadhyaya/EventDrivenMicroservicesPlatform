package notification

import (
	"fmt"
	"log"

	"github.com/Garvit-Jethwani/notification-service/events"
	"github.com/Garvit-Jethwani/notification-service/models"
)

type NotificationService struct {
	consumer events.EventConsumer
}

func NewNotificationService(consumer events.EventConsumer) *NotificationService {
	return &NotificationService{
		consumer: consumer,
	}
}

func (n *NotificationService) Start() error {
	eventChan, err := n.consumer.Consume()
	if err != nil {
		return err
	}

	for event := range eventChan {
		err = n.processEvent(event)
		if err != nil {
			// Handle error
			log.Printf("Error processing event: %v", err)
		}
	}

	return nil
}

func (n *NotificationService) processEvent(event models.Event) error {
	// Example processing logic
	fmt.Println("Processing event:", event)
	switch event.Type {
	case "order_created":
		n.sendOrderCreatedNotification(event)
	case "order_updated":
		n.sendOrderUpdatedNotification(event)
	case "order_shipped":
		n.sendOrderShippedNotification(event)
	default:
		log.Printf("Unhandled event type: %s", event.Type)
	}
	return nil
}

func (n *NotificationService) sendOrderCreatedNotification(event models.Event) {
	// Implement the logic to send 'order created' notification
	log.Printf("Sending order created notification for event: %v", event)
}

func (n *NotificationService) sendOrderUpdatedNotification(event models.Event) {
	// Implement the logic to send 'order updated' notification
	log.Printf("Sending order updated notification for event: %v", event)
}

func (n *NotificationService) sendOrderShippedNotification(event models.Event) {
	// Implement the logic to send 'order shipped' notification
	log.Printf("Sending order shipped notification for event: %v", event)
}
