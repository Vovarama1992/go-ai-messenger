package app

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/dto"
	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/stream"
	"github.com/segmentio/kafka-go"
)

func RunAdviceReaderFromKafka(ctx context.Context, concurrency int, reader *kafka.Reader) {
	for i := 0; i < concurrency; i++ {
		go func(workerID int) {
			log.Printf("📥 [advice_reader_from_kafka #%d] started", workerID)

			for {
				select {
				case <-ctx.Done():
					log.Printf("🛑 [advice_reader_from_kafka #%d] stopped", workerID)
					return

				default:
					m, err := reader.ReadMessage(ctx)
					if err != nil {
						log.Printf("❌ [advice_reader_from_kafka #%d] read error: %v", workerID, err)
						continue
					}

					var payload dto.AdviceRequestPayload
					if err := json.Unmarshal(m.Value, &payload); err != nil {
						log.Printf("❌ [advice_reader_from_kafka #%d] invalid JSON: %v", workerID, err)
						continue
					}

					stream.AdviceRequestChan <- payload
				}
			}
		}(i)
	}
}
