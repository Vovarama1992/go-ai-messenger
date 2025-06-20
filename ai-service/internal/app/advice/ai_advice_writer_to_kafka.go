package app

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/ports"
	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/stream"
)

// RunAdviceWriterToKafka
// Читает dto.GptAdvice из канала AdviceResponseChan и публикует
// в Kafka-топик TOPIC_AI_ADVICE_RESPONSE
func RunAdviceWriterToKafka(
	ctx context.Context,
	producer ports.KafkaProducer,
	topic string,
) {
	go func() {
		log.Printf("📤 [advice_writer_to_kafka] started (topic=%s)", topic)

		for {
			select {
			case <-ctx.Done():
				log.Println("🛑 [advice_writer_to_kafka] stopped")
				return

			case advice := <-stream.AdviceResponseChan:
				data, err := json.Marshal(advice)
				if err != nil {
					log.Printf("❌ marshal advice: %v", err)
					continue
				}
				if err := producer.Publish(ctx, topic, data); err != nil {
					log.Printf("❌ publish advice: %v", err)
					continue
				}
				log.Printf("✅ published advice for threadID=%s", advice.ThreadID)
			}
		}
	}()
}
