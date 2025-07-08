package feed_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	app "github.com/Vovarama1992/go-ai-messenger/ai-service/internal/app/feed"
	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/dto"
	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/mocks"
	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/stream"
)

func TestRunAiFeedReaderFromKafka(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReader := mocks.NewMockKafkaReader(ctrl)

	payload := dto.AiFeedPayload{
		ThreadID:    "thread-1",
		SenderEmail: "user@example.com",
		Text:        "что делать?",
		BindingType: "autoreply",
	}
	raw, _ := json.Marshal(payload)

	mockReader.EXPECT().ReadMessage(gomock.Any()).Return(raw, nil).Times(1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	origChan := stream.FeedChan
	testChan := make(chan dto.AiFeedPayload, 1)
	stream.FeedChan = testChan
	defer func() { stream.FeedChan = origChan }()

	go app.RunAiFeedReaderFromKafka(ctx, 1, mockReader)

	select {
	case got := <-testChan:
		require.Equal(t, payload.Text, got.Text)
		require.Equal(t, payload.BindingType, got.BindingType)
	case <-time.After(2 * time.Second):
		t.Fatal("no message received from channel")
	}
}

func TestRunAiFeedReaderFromGpt_Autoreply(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGpt := mocks.NewMockGptClient(ctrl)

	payload := dto.AiFeedPayload{
		ThreadID:    "thread-1",
		SenderEmail: "user@example.com",
		Text:        "что делать?",
		BindingType: "autoreply",
	}
	expectedReply := "ответ от GPT"

	mockGpt.
		EXPECT().
		SendMessageAndGetAutoreply(gomock.Any(), payload.ThreadID, payload.SenderEmail, payload.Text).
		Return(expectedReply, nil).
		Times(1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	origIn := stream.FeedChan
	origOut := stream.AutoReplyChan

	in := make(chan dto.AiFeedPayload, 1)
	out := make(chan dto.AiAutoReplyResult, 1)

	stream.FeedChan = in
	stream.AutoReplyChan = out

	defer func() {
		stream.FeedChan = origIn
		stream.AutoReplyChan = origOut
	}()

	in <- payload

	go app.RunAiFeedReaderFromGpt(ctx, 1, mockGpt)

	select {
	case got := <-out:
		require.Equal(t, payload.ThreadID, got.ThreadID)
		require.Equal(t, expectedReply, got.Text)
	case <-time.After(2 * time.Second):
		t.Fatal("no reply received")
	}
}

func TestRunAiFeedWriterToKafka(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWriter := mocks.NewMockKafkaProducer(ctrl)

	payload := dto.AiAutoReplyResult{
		ThreadID: "thread-1",
		Text:     "ответ от GPT",
	}

	mockWriter.
		EXPECT().
		Publish(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	origChan := stream.AutoReplyChan
	testChan := make(chan dto.AiAutoReplyResult, 1)
	stream.AutoReplyChan = testChan
	defer func() { stream.AutoReplyChan = origChan }()

	testChan <- payload

	go app.RunAiFeedWriterToKafka(ctx, mockWriter)

	select {
	case <-time.After(1 * time.Second):
		// проверка сработает по EXPECT
	case <-ctx.Done():
		t.Fatal("context canceled too early")
	}
}
