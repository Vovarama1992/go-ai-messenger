package app_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/app"
	"github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/mocks"
	"github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/model"
	"github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/stream"
)

func TestRunAdviceReaderFromChan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockChat := mocks.NewMockChatService(ctrl)
	mockChat.EXPECT().
		GetUserWithChatByThreadID("thread-abc").
		Return(int64(123), int64(456), "user@example.com", nil).
		Times(1)

	go app.RunChannelsBetweener(ctx, mockChat)

	stream.PendingAdviceChan <- model.GptAdvice{
		ThreadID: "thread-abc",
		Text:     "Привет",
	}

	select {
	case enriched := <-stream.EnrichedAdviceChan:
		require.Equal(t, int64(123), enriched.UserID)
		require.Equal(t, int64(456), enriched.ChatID)
		require.Equal(t, "Привет", enriched.Text)
	case <-time.After(1 * time.Second):
		t.Fatal("timeout: enriched advice not received")
	}
}

func TestRunAdviceReaderFromChan_ChatServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockChat := mocks.NewMockChatService(ctrl)
	mockChat.EXPECT().
		GetUserWithChatByThreadID("bad-thread").
		Return(int64(0), int64(0), "", errors.New("boom")).
		Times(1)

	go app.RunChannelsBetweener(ctx, mockChat)

	stream.PendingAdviceChan <- model.GptAdvice{
		ThreadID: "bad-thread",
		Text:     "не дойдёт",
	}

	select {
	case msg := <-stream.EnrichedAdviceChan:
		t.Fatalf("ожидали, что ничего не придёт, но получили: %+v", msg)
	case <-time.After(300 * time.Millisecond):
		// всё ок — ничего не пришло
	}
}
