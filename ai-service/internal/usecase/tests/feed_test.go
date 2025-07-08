package usecase_test

import (
	"context"
	"testing"
	"time"

	app "github.com/Vovarama1992/go-ai-messenger/ai-service/internal/app/feed"
	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/dto"
	mocks "github.com/Vovarama1992/go-ai-messenger/ai-service/internal/mocks"
	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/stream"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRunAiFeedReaderFromGpt_Autoreply(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGpt := mocks.NewMockGptClient(ctrl)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	payload := dto.AiFeedPayload{
		ThreadID:    "thread-1",
		SenderEmail: "user@example.com",
		Text:        "что делать?",
		BindingType: "autoreply",
	}

	expectedReply := "ответ от GPT"

	mockGpt.EXPECT().
		SendMessageAndGetAutoreply(gomock.Any(), payload.ThreadID, payload.SenderEmail, payload.Text).
		Return(expectedReply, nil)

	// стартуем обработчик
	app.RunAiFeedReaderFromGpt(ctx, 1, mockGpt)

	// отправляем в канал
	stream.FeedChan <- payload

	// читаем из результата
	select {
	case res := <-stream.AutoReplyChan:
		require.Equal(t, payload.ThreadID, res.ThreadID)
		require.Equal(t, expectedReply, res.Text)
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for response")
	}
}
