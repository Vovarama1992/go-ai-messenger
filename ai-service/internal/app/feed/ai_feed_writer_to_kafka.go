package app

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/infra/kafka"
	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/stream"
)

func RunAiFeedWriterToKafka(ctx context.Context, producer *kafka.Producer) {
	go func() {
		log.Println("📤 [ai_feed_writer_to_kafka] started")

		for {
			select {
			case <-ctx.Done():
				log.Println("🛑 [ai_feed_writer_to_kafka] stopped")
				return

			case msg := <-stream.AutoReplyChan:
				delay := time.Duration(rand.Intn(21)+5) * time.Second
				time.Sleep(delay)

				payload := map[string]interface{}{
					"threadId": msg.ThreadID,
					"text":     msg.Text,
				}
				data, err := json.Marshal(payload)
				if err != nil {
					log.Printf("❌ Failed to marshal autoreply: %v", err)
					continue
				}

				if err := producer.Publish(ctx, msg.ThreadID, data); err != nil {
					log.Printf("❌ Kafka publish failed: %v", err)
				} else {
					log.Printf("✅ Autoreply published (threadID=%s)", msg.ThreadID)
				}
			}
		}
	}()
}
