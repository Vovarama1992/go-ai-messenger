package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/dto"

	kafka "github.com/segmentio/kafka-go"
)

func StartKafkaConsumer(ctx context.Context, r *kafka.Reader, out chan<- dto.IncomingMessage) {
	go func() {
		defer r.Close()
		log.Println("📥 Kafka consumer started")

		for {
			select {
			case <-ctx.Done():
				log.Println("🛑 Kafka consumer shutdown")
				return
			default:
				m, err := r.ReadMessage(ctx)
				if err != nil {
					log.Printf("❌ Kafka read error: %v", err)
					continue
				}

				var msg dto.IncomingMessage
				if err := json.Unmarshal(m.Value, &msg); err != nil {
					log.Printf("❌ Invalid message JSON: %v", err)
					continue
				}

				out <- msg
			}
		}
	}()
}
