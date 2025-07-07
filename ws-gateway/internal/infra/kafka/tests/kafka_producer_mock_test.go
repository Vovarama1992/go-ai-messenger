package kafkaadapter_test

import (
	"context"
	"testing"

	"github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/mocks"
	"go.uber.org/mock/gomock"
)

func TestKafkaProducer_Produce(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProducer := mocks.NewMockKafkaProducer(ctrl)

	mockProducer.EXPECT().
		Produce(gomock.Any(), "test-topic", gomock.Any()).
		Return(nil).
		Times(1)

	err := mockProducer.Produce(context.Background(), "test-topic", map[string]interface{}{
		"key": "value",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
