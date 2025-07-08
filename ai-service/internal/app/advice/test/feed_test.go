package advice_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	advice "github.com/Vovarama1992/go-ai-messenger/ai-service/internal/app/advice"
	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/dto"
	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/mocks"
	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/stream"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRunAdviceReaderFromKafka(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReader := mocks.NewMockKafkaReader(ctrl)

	payload := dto.AdviceRequestPayload{ThreadID: "thread-123"}
	raw, _ := json.Marshal(payload)

	mockReader.
		EXPECT().
		ReadMessage(gomock.Any()).
		Return(raw, nil).
		Times(1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	origChan := stream.AdviceRequestChan
	testChan := make(chan dto.AdviceRequestPayload, 1)
	stream.AdviceRequestChan = testChan
	defer func() { stream.AdviceRequestChan = origChan }()

	go advice.RunAdviceReaderFromKafka(ctx, 1, mockReader)

	select {
	case got := <-testChan:
		require.Equal(t, payload.ThreadID, got.ThreadID)
	case <-time.After(2 * time.Second):
		t.Fatal("no message received from channel")
	}
}

func TestRunAdviceReaderFromGpt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGpt := mocks.NewMockGptClient(ctrl)

	threadID := "thread-xyz"
	expectedText := "важный совет"
	payload := dto.AdviceRequestPayload{ThreadID: threadID}

	mockGpt.
		EXPECT().
		GetAdvice(gomock.Any(), threadID).
		Return(expectedText, nil).
		Times(1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	origIn := stream.AdviceRequestChan
	origOut := stream.AdviceResponseChan
	in := make(chan dto.AdviceRequestPayload, 1)
	out := make(chan dto.GptAdvice, 1)
	stream.AdviceRequestChan = in
	stream.AdviceResponseChan = out
	defer func() {
		stream.AdviceRequestChan = origIn
		stream.AdviceResponseChan = origOut
	}()

	in <- payload

	go advice.RunAdviceReaderFromGpt(ctx, 1, mockGpt)

	select {
	case res := <-out:
		require.Equal(t, threadID, res.ThreadID)
		require.Equal(t, expectedText, res.Text)
	case <-time.After(2 * time.Second):
		t.Fatal("no thread result received")
	}
}
func TestRunAdviceWriterToKafka(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProducer := mocks.NewMockKafkaProducer(ctrl)

	payload := dto.GptAdvice{
		ThreadID: "thread-xyz",
		Text:     "готовый совет",
	}
	topic := "topic-ai-advice-response"

	mockProducer.
		EXPECT().
		Publish(gomock.Any(), topic, gomock.Any()).
		Return(nil).
		Times(1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	origChan := stream.AdviceResponseChan
	testChan := make(chan dto.GptAdvice, 1)
	stream.AdviceResponseChan = testChan
	defer func() { stream.AdviceResponseChan = origChan }()

	testChan <- payload

	go advice.RunAdviceWriterToKafka(ctx, mockProducer, topic)

	select {
	case <-time.After(2 * time.Second):
		// SUCCESS: mock.Expectation выполнена
	case <-ctx.Done():
		t.Fatal("context cancelled too early")
	}
}
