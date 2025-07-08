package binding_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	binding "github.com/Vovarama1992/go-ai-messenger/ai-service/internal/app/binding"
	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/dto"
	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/mocks"
	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/stream"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRunAiCreateBindingReaderFromKafka(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReader := mocks.NewMockKafkaReader(ctrl)

	payload := dto.AiBindingInitPayload{
		UserID:    1,
		UserEmail: "test@example.com",
		ChatID:    2,
		Messages: []dto.ChatMessage{
			{SenderEmail: "a@b", Text: "hello"},
		},
	}
	raw, _ := json.Marshal(payload)

	mockReader.EXPECT().ReadMessage(gomock.Any()).Return(raw, nil).Times(1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	origChan := stream.BindingInitChan
	testChan := make(chan dto.AiBindingInitPayload, 1)
	stream.BindingInitChan = testChan
	defer func() { stream.BindingInitChan = origChan }()

	go binding.RunAiCreateBindingReaderFromKafka(ctx, 1, mockReader)

	select {
	case got := <-testChan:
		require.Equal(t, payload.ChatID, got.ChatID)
		require.Equal(t, payload.UserEmail, got.UserEmail)
	case <-time.After(2 * time.Second):
		t.Fatal("no message received from channel")
	}
}

func TestRunAiCreateBindingReaderFromGpt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGpt := mocks.NewMockGptClient(ctrl)

	payload := dto.AiBindingInitPayload{
		UserID:    1,
		UserEmail: "test@example.com",
		ChatID:    2,
		Messages: []dto.ChatMessage{
			{SenderEmail: "a@b", Text: "hello"},
		},
	}
	threadID := "thread-abc"

	mockGpt.
		EXPECT().
		CreateThreadForUserAndChat(gomock.Any(), payload.UserEmail, payload.Messages).
		Return(threadID, nil).
		Times(1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// подменяем каналы
	origIn := stream.BindingInitChan
	origOut := stream.BindingResultChan

	in := make(chan dto.AiBindingInitPayload, 1)
	out := make(chan dto.ThreadResult, 1)

	stream.BindingInitChan = in
	stream.BindingResultChan = out

	defer func() {
		stream.BindingInitChan = origIn
		stream.BindingResultChan = origOut
	}()

	// кидаем тестовый payload
	in <- payload

	go binding.RunAiCreateBindingReaderFromGpt(ctx, 1, mockGpt)

	select {
	case got := <-out:
		require.Equal(t, payload.ChatID, got.ChatID)
		require.Equal(t, threadID, got.ThreadID)
	case <-time.After(2 * time.Second):
		t.Fatal("no thread result received")
	}
}

func TestRunAiCreateBindingWriterToKafka(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProducer := mocks.NewMockKafkaProducer(ctrl)

	payload := dto.ThreadResult{
		UserID:   1,
		ChatID:   2,
		ThreadID: "thread-abc",
	}

	mockProducer.
		EXPECT().
		Publish(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	origChan := stream.BindingResultChan
	testChan := make(chan dto.ThreadResult, 1)
	stream.BindingResultChan = testChan
	defer func() { stream.BindingResultChan = origChan }()

	testChan <- payload

	go binding.RunAiCreateBindingWriterToKafka(ctx, mockProducer, "dummy-topic")

	select {
	case <-time.After(1 * time.Second):
		// если Publish не вызовется — тест упадёт по EXPECT выше
	case <-ctx.Done():
		t.Fatal("context canceled too early")
	}
}
