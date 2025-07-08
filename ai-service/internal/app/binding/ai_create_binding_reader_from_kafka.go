package app

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/dto"
	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/ports"
	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/stream"
)

// ai_create_binding_reader_from_kafka
// Читает сообщения из Kafka топика chat.binding.init и пишет в BindingInitChan
func RunAiCreateBindingReaderFromKafka(ctx context.Context, concurrency int, reader ports.KafkaReader) {
	for i := 0; i < concurrency; i++ {
		go func(workerID int) {
			log.Printf("📥 [ai_create_binding_reader_from_kafka #%d] started", workerID)

			for {
				select {
				case <-ctx.Done():
					log.Printf("🛑 [ai_create_binding_reader_from_kafka #%d] stopped", workerID)
					return

				default:
					msg, err := reader.ReadMessage(ctx)
					if err != nil {
						log.Printf("❌ [reader #%d] Kafka read error: %v", workerID, err)
						continue
					}

					var payload dto.AiBindingInitPayload
					if err := json.Unmarshal(msg, &payload); err != nil {
						log.Printf("❌ [reader #%d] Invalid JSON: %v", workerID, err)
						continue
					}

					stream.BindingInitChan <- payload
				}
			}
		}(i)
	}
}
