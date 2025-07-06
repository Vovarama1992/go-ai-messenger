package postgres_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/dto"
	"github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/infra/postgres"
	"github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/mocks"
	"github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/model"
	"github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/usecase"
)

func TestHandle_WithThreadID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMessageRepo(ctrl)
	mockUser := mocks.NewMockUserClient(ctrl)
	mockChat := mocks.NewMockChatClient(ctrl)

	service := usecase.NewMessageService(mockRepo, mockUser)
	processor := postgres.NewDefaultMessageProcessor(service, mockChat)

	payload := dto.IncomingMessage{
		ThreadID: "thread-123",
		Text:     "From thread",
	}

	mockChat.
		EXPECT().
		GetThreadInfo(gomock.Any(), "thread-123").
		Return(model.ThreadContext{
			UserID: 123,
			ChatID: 456,
		}, nil)

	mockRepo.
		EXPECT().
		Save(gomock.Any()).
		Return(nil)

	err := processor.Handle(context.Background(), payload)
	assert.NoError(t, err)
}

func TestHandle_SaveError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMessageRepo(ctrl)
	mockUser := mocks.NewMockUserClient(ctrl)
	mockChat := mocks.NewMockChatClient(ctrl)

	service := usecase.NewMessageService(mockRepo, mockUser)
	processor := postgres.NewDefaultMessageProcessor(service, mockChat)

	payload := dto.IncomingMessage{
		ChatID:   1,
		SenderID: 2,
		Text:     "error case",
	}

	mockRepo.
		EXPECT().
		Save(gomock.Any()).
		Return(errors.New("insert failed"))

	err := processor.Handle(context.Background(), payload)
	assert.Error(t, err)
	assert.EqualError(t, err, "insert failed")
}

func TestHandle_InvalidThreadInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMessageRepo(ctrl)
	mockUser := mocks.NewMockUserClient(ctrl)
	mockChat := mocks.NewMockChatClient(ctrl)

	service := usecase.NewMessageService(mockRepo, mockUser)
	processor := postgres.NewDefaultMessageProcessor(service, mockChat)

	payload := dto.IncomingMessage{
		ThreadID: "unknown-thread",
		Text:     "fail thread",
	}

	mockChat.
		EXPECT().
		GetThreadInfo(gomock.Any(), "unknown-thread").
		Return(model.ThreadContext{}, errors.New("not found"))

	err := processor.Handle(context.Background(), payload)
	assert.Error(t, err)
	assert.EqualError(t, err, "not found")
}
