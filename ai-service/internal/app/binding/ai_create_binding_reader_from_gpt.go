package app

import (
	"context"
	"log"

	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/dto"
	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/ports"
	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/stream"
	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/usecase"
)

// ai_create_binding_reader_from_gpt
// Читает BindingInitPayload из канала, вызывает GPT и пишет ThreadResult в канал
func RunAiCreateBindingReaderFromGpt(ctx context.Context, concurrency int, gpt ports.GptService) {
	sem := make(chan struct{}, concurrency)

	go func() {
		log.Printf("🤖 [ai_create_binding_reader_from_gpt] started (%d concurrent)", concurrency)

		for {
			select {
			case <-ctx.Done():
				log.Println("🛑 [ai_create_binding_reader_from_gpt] stopped")
				return

			case payload := <-stream.BindingInitChan:
				sem <- struct{}{}

				go func(payload dto.AiBindingInitPayload) {
					defer func() { <-sem }()

					res, err := usecase.ProcessBindingInit(ctx, payload, gpt)
					if err != nil {
						log.Printf("❌ GPT processing error: %v", err)
						return
					}

					stream.BindingResultChan <- res
				}(payload)
			}
		}
	}()
}
