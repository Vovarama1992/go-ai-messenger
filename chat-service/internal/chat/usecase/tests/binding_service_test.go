package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	mocks "github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/mocks"
	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/model"
	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/usecase"
	"go.uber.org/mock/gomock"
)

func TestUpdateBinding_CreatesAndPushes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockChatBindingRepository(ctrl)
	mockBroker := mocks.NewMockMessageBroker(ctrl)
	mockClient := mocks.NewMockMessageClient(ctrl)

	service := usecase.NewChatBindingService(mockRepo, mockBroker, mockClient)

	ctx := context.Background()
	userID := int64(123)
	chatID := int64(456)
	email := "user@example.com"
	bindingType := model.AIBindingAdvice

	mockRepo.EXPECT().
		FindByUserAndChat(ctx, userID, chatID).
		Return(nil, errors.New("not found"))

	mockRepo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(nil)

	mockClient.EXPECT().
		GetMessagesByChat(ctx, chatID).
		Return([]*model.ChatMessage{
			{
				SenderEmail: "user@example.com",
				Text:        "Hello",
				SentAt:      time.Now().Unix(),
			},
		}, nil)

	mockBroker.EXPECT().
		SendAiBindingInit(gomock.Any(), gomock.Any()).
		Return(nil)

	err := service.UpdateBinding(ctx, email, userID, chatID, bindingType)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHandleThreadCreated_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockChatBindingRepository(ctrl)
	mockBroker := mocks.NewMockMessageBroker(ctrl)
	mockClient := mocks.NewMockMessageClient(ctrl)

	service := usecase.NewChatBindingService(mockRepo, mockBroker, mockClient)

	res := model.ThreadResult{
		ChatID:   1,
		UserID:   2,
		ThreadID: "thread-abc",
	}

	mockRepo.EXPECT().
		UpdateThreadID(gomock.Any(), res.ChatID, res.UserID, res.ThreadID).
		Return(nil)

	service.HandleThreadCreated(res)
}

func TestHandleThreadCreated_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockChatBindingRepository(ctrl)
	mockBroker := mocks.NewMockMessageBroker(ctrl)
	mockClient := mocks.NewMockMessageClient(ctrl)

	service := usecase.NewChatBindingService(mockRepo, mockBroker, mockClient)

	res := model.ThreadResult{
		ChatID:   1,
		UserID:   2,
		ThreadID: "thread-xyz",
	}

	mockRepo.EXPECT().
		UpdateThreadID(gomock.Any(), res.ChatID, res.UserID, res.ThreadID).
		Return(errors.New("db error"))

	service.HandleThreadCreated(res)
}
