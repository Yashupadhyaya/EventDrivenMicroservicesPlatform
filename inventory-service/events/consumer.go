package events

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Yashupadhyaya/inventory-service/config"
	"github.com/Yashupadhyaya/inventory-service/database"
	"github.com/segmentio/kafka-go"
)

const (
	topic   = "order_events"
	groupID = "github.com/Yashupadhyaya/inventory-service"
)

func StartConsumer(cfg *config.Config) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{cfg.EventStoreURL},
		GroupID: groupID,
		Topic:   topic,
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("could not read message: ", err)
			continue
		}

		var event map[string]interface{}
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Println("could not unmarshal message: ", err)
			continue
		}

		if event["event"] == "order_created" {
			handleOrderCreatedEvent(event["data"].(map[string]interface{}))
		}
	}
}

func handleOrderCreatedEvent(data map[string]interface{}) {
	for _, item := range data["items"].([]interface{}) {
		itemMap := item.(map[string]interface{})
		productID := itemMap["product_id"].(string)
		quantity := int(itemMap["quantity"].(float64))

		inventory, err := database.GetInventoryByProductID(productID)
		if err != nil {
			log.Println("error getting inventory: ", err)
			continue
		}

		newQuantity := inventory.Quantity - quantity
		if err := database.UpdateInventory(productID, newQuantity); err != nil {
			log.Println("error updating inventory: ", err)
		}
	}
}
