package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"

	dto "github.com/Vovarama1992/go-ai-messenger/ai-service/internal/dto"
)

var MessageChan = make(chan dto.AiBindingInitPayload, 100)

func StartKafkaConsumer(ctx context.Context, r *kafka.Reader) {
	go func() {
		defer r.Close()
		log.Println("📥 Kafka consumer started (binding.init)")

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

				var payload dto.AiBindingInitPayload
				if err := json.Unmarshal(m.Value, &payload); err != nil {
					log.Printf("❌ Invalid JSON: %v", err)
					continue
				}

				MessageChan <- payload
			}
		}
	}()
}
