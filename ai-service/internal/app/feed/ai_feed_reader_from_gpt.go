package app

import (
	"context"
	"fmt"
	"log"

	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/dto"
	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/ports"
	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/stream"
	"github.com/sashabaranov/go-openai"
)

// ai_feed_reader_from_gpt
// Читает из FeedChan и шлёт сообщение в GPT-тред
func RunAiFeedReaderFromGpt(ctx context.Context, concurrency int, gpt ports.GptService) {
	sem := make(chan struct{}, concurrency)

	go func() {
		log.Printf("🤖 [ai_feed_reader_from_gpt] started (%d concurrent)", concurrency)

		for {
			select {
			case <-ctx.Done():
				log.Println("🛑 [ai_feed_reader_from_gpt] stopped")
				return

			case payload := <-stream.FeedChan:
				sem <- struct{}{}

				go func(p dto.AiFeedPayload) {
					defer func() { <-sem }()

					switch p.BindingType {
					case "autoreply":
						reply, err := gpt.SendMessageAndGetAutoreply(ctx, p.ThreadID, p.SenderEmail, p.Text)
						if err != nil {
							log.Printf("❌ GPT autoreply error for thread %s: %v", p.ThreadID, err)
							return
						}

						stream.AutoReplyChan <- dto.AiAutoReplyResult{
							ThreadID: p.ThreadID,
							Text:     reply,
						}

						log.Printf("✅ Autoreply pushed to channel for thread %s", p.ThreadID)

					default:
						message := fmt.Sprintf("бро лови новое сообщение из диалога. емейл: %s, текст: %s", p.SenderEmail, p.Text)

						err := gpt.SendMessageToThread(ctx, p.ThreadID, openai.ChatMessageRoleUser, message)
						if err != nil {
							log.Printf("❌ GPT send error for thread %s: %v", p.ThreadID, err)
							return
						}
						log.Printf("✅ Sent message to GPT thread %s", p.ThreadID)
					}
				}(payload)
			}
		}
	}()
}
