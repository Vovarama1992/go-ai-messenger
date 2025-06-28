package kafka_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"

	kafkaadapter "github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/infra/kafka"
	"github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/model"
	"github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/stream"
)

type fakeReader struct {
	messages []kafka.Message
	index    int
	errOnce  bool
}

func (f *fakeReader) ReadMessage(ctx context.Context) (kafka.Message, error) {
	// вернём ошибку один раз, если нужно
	if f.errOnce {
		f.errOnce = false
		return kafka.Message{}, errors.New("read error")
	}

	if f.index >= len(f.messages) {
		<-ctx.Done()
		return kafka.Message{}, ctx.Err()
	}

	msg := f.messages[f.index]
	f.index++
	return msg, nil
}

func (f *fakeReader) Close() error { return nil }

func TestAdviceConsumer_Start_PushesToChan(t *testing.T) {
	// valидное сообщение
	valid := model.GptAdvice{
		ThreadID: "t1",
		Text:     "Привет",
	}
	rawValid, _ := json.Marshal(valid)

	// невалидный JSON
	rawInvalid := []byte(`{not-json}`)

	// пустое сообщение
	rawEmpty := []byte(``)

	reader := &fakeReader{
		messages: []kafka.Message{
			{Value: rawInvalid},
			{Value: rawEmpty},
			{Value: rawValid},
		},
		errOnce: true, // первая попытка вернёт ошибку
	}

	consumer := kafkaadapter.NewAdviceConsumer(reader)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := consumer.Start(ctx)
	require.NoError(t, err)

	select {
	case got := <-stream.PendingAdviceChan:
		require.Equal(t, valid.ThreadID, got.ThreadID)
		require.Equal(t, valid.Text, got.Text)
	case <-time.After(1 * time.Second):
		t.Fatal("timeout: message not received")
	}
}
