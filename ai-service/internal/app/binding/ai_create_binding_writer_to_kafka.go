package app

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/ports"
	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/stream"
)

// ai_create_binding_writer_to_kafka
// Читает ThreadResult из канала и публикует его в Kafka
func RunAiCreateBindingWriterToKafka(ctx context.Context, producer ports.KafkaProducer, topic string) {
	go func() {
		log.Printf("📤 [ai_create_binding_writer_to_kafka] started (topic=%s)", topic)

		for {
			select {
			case <-ctx.Done():
				log.Println("🛑 [ai_create_binding_writer_to_kafka] stopped")
				return

			case result := <-stream.BindingResultChan:
				data, err := json.Marshal(result)
				if err != nil {
					log.Printf("❌ Failed to marshal ThreadResult: %v", err)
					continue
				}

				err = producer.Publish(ctx, "", data)
				if err != nil {
					log.Printf("❌ Failed to publish to Kafka: %v", err)
					continue
				}

				log.Printf("✅ Published threadID=%s to Kafka", result.ThreadID)
			}
		}
	}()
}
