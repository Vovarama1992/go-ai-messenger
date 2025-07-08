package app

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/dto"
	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/ports"
	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/stream"
)

func RunAiFeedReaderFromKafka(ctx context.Context, concurrency int, reader ports.KafkaReader) {
	for i := 0; i < concurrency; i++ {
		go func(workerID int) {
			log.Printf("📥 [ai_feed_reader_from_kafka #%d] started", workerID)

			for {
				select {
				case <-ctx.Done():
					log.Printf("🛑 [ai_feed_reader_from_kafka #%d] stopped", workerID)
					return

				default:
					raw, err := reader.ReadMessage(ctx)
					if err != nil {
						log.Printf("❌ [reader #%d] Kafka read error: %v", workerID, err)
						continue
					}

					var payload dto.AiFeedPayload
					if err := json.Unmarshal(raw, &payload); err != nil {
						log.Printf("❌ [reader #%d] Invalid JSON: %v", workerID, err)
						continue
					}

					stream.FeedChan <- payload
				}
			}
		}(i)
	}
}
