package app

import (
	"context"
	"log"

	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/dto"
	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/ports"
	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/stream"
	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/usecase"
)

func RunAdviceReaderFromGpt(ctx context.Context, concurrency int, gpt ports.GptClient) {
	sem := make(chan struct{}, concurrency)

	go func() {
		log.Printf("🤖 [advice_reader_from_gpt] started (%d concurrent)", concurrency)

		for {
			select {
			case <-ctx.Done():
				log.Println("🛑 [advice_reader_from_gpt] stopped")
				return

			case payload := <-stream.AdviceRequestChan:
				sem <- struct{}{}

				go func(payload dto.AdviceRequestPayload) {
					defer func() { <-sem }()

					result, err := usecase.ProcessAdviceRequest(ctx, payload, gpt)
					if err != nil {
						log.Printf("❌ GPT error on advice: %v", err)
						return
					}

					stream.AdviceResponseChan <- result
				}(payload)
			}
		}
	}()
}
