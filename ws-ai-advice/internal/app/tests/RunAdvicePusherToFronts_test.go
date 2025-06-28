package app_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/app"
	"github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/mocks"
	"github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/model"
	"github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/stream"
)

func TestRunAdviceWriterToFronts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockHub := mocks.NewMockAdviceHub(ctrl)

	done := make(chan struct{}, 1)

	mockHub.EXPECT().
		Send(int64(123), "gpt-advice", map[string]any{
			"chatId": int64(456),
			"text":   "ответ ГПТ",
		}).
		Do(func(userID int64, event string, payload any) {
			done <- struct{}{}
		}).
		Times(1)

	go app.RunAdvicePusherToFronts(ctx, mockHub) // 🔧 исправлено имя функции

	stream.EnrichedAdviceChan <- model.EnrichedAdvice{
		UserID: 123,
		ChatID: 456,
		Text:   "ответ ГПТ",
	}

	select {
	case <-done:
		require.True(t, true, "hub.Send was called")
	case <-time.After(time.Second):
		t.Fatal("timeout: hub.Send was not called")
	}
}
