package usecase

import (
	"context"
	"log"

	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/dto"
	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/ports"
)

var ResultChan = make(chan dto.ThreadResult, 100)

func ProcessBindingInit(ctx context.Context, payload dto.AiBindingInitPayload, gpt ports.GptClient) (dto.ThreadResult, error) {
	threadID, err := gpt.CreateThreadForUserAndChat(
		ctx,
		payload.UserEmail,
		payload.Messages,
	)
	if err != nil {
		return dto.ThreadResult{}, err
	}

	log.Printf("✅ GPT thread created: chatID=%d userID=%d threadID=%s", payload.ChatID, payload.UserID, threadID)

	return dto.ThreadResult{
		ChatID:   payload.ChatID,
		UserID:   payload.UserID,
		ThreadID: threadID,
	}, nil
}
