package app

import (
	"context"
	"log"

	"github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/model"
	"github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/ports"
	"github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/stream"
)

func RunChannelsBetweener(ctx context.Context, chatService ports.ChatService) {
	workerCount := 10

	for i := 0; i < workerCount; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case advice := <-stream.PendingAdviceChan:
					userID, chatID, _, err := chatService.GetUserWithChatByThreadID(advice.ThreadID)
					if err != nil {
						log.Printf("❌ ошибка обогащения advice: %v", err)
						continue
					}

					enriched := model.EnrichedAdvice{
						UserID: userID,
						ChatID: chatID,
						Text:   advice.Text,
					}
					stream.EnrichedAdviceChan <- enriched
				}
			}
		}()
	}
}
