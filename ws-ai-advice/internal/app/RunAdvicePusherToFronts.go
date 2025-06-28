package app

import (
	"context"

	"github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/ports"
	"github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/stream"
)

func RunAdvicePusherToFronts(ctx context.Context, hub ports.AdviceHub) {
	for {
		select {
		case <-ctx.Done():
			return
		case enriched := <-stream.EnrichedAdviceChan:
			hub.Send(enriched.UserID, "gpt-advice", map[string]any{
				"chatId": enriched.ChatID,
				"text":   enriched.Text,
			})
		}
	}
}
