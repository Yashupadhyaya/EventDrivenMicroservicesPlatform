package events

// const (
// 	topic     = "order_events"
// 	partition = 0
// )

// var writer *kafka.Writer

// func InitEventProducer(brokers []string) {
// 	writer = &kafka.Writer{
// 		Addr:     kafka.TCP(brokers...),
// 		Topic:    topic,
// 		Balancer: &kafka.LeastBytes{},
// 	}
// }

// func ProduceOrderCreatedEvent(order *models.Order) error {
// 	event := map[string]interface{}{
// 		"event": "order_created",
// 		"data":  order,
// 	}
// 	payload, err := json.Marshal(event)
// 	if err != nil {
// 		return err
// 	}

// 	err = writer.WriteMessages(nil, kafka.Message{
// 		Value: payload,
// 	})
// 	if err != nil {
// 		log.Printf("could not write message: %v", err)
// 		return err
// 	}

// 	log.Printf("order created event produced: %v", order)
// 	return nil
// }
