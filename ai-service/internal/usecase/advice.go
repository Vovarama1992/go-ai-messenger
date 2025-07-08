package usecase

import (
	"context"
	"log"

	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/dto"
	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/ports"
)

func ProcessAdviceRequest(ctx context.Context, payload dto.AdviceRequestPayload, gpt ports.GptClient) (dto.GptAdvice, error) {
	log.Printf("💬 [advice] requesting advice for threadID=%s", payload.ThreadID)

	answer, err := gpt.GetAdvice(ctx, payload.ThreadID)
	if err != nil {
		return dto.GptAdvice{}, err
	}

	return dto.GptAdvice{
		ThreadID: payload.ThreadID,
		Text:     answer,
	}, nil
}
