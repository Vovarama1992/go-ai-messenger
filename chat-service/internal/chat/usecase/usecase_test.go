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

var ErrNotFound = errors.New("chat not found")

func newChatServiceWithMocks(t *testing.T) (*gomock.Controller, *mocks.MockChatRepository, *usecase.ChatService) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockChatRepository(ctrl)
	mockBindingRepo := mocks.NewMockChatBindingRepository(ctrl)
	mockAdvicePublisher := mocks.NewMockAdvicePublisher(ctrl)

	service := usecase.NewChatService(mockRepo, mockBindingRepo, mockAdvicePublisher)

	return ctrl, mockRepo, service
}

func TestCreateChat_Success(t *testing.T) {
	ctrl, mockRepo, service := newChatServiceWithMocks(t)
	defer ctrl.Finish()

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, chat *model.Chat) error {
			chat.ID = 1
			chat.CreatedAt = time.Now().Unix()
			return nil
		})

	chatType := model.ChatTypePrivate
	chat, err := service.CreateChat(context.Background(), 123, chatType)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if chat.CreatorID != 123 {
		t.Fatalf("expected CreatorID 123, got %d", chat.CreatorID)
	}
	if chat.Type != chatType {
		t.Fatalf("expected Type %s, got %s", chatType, chat.Type)
	}
	if chat.ID == 0 {
		t.Fatal("expected chat ID to be set")
	}
}

func TestCreateChat_InvalidType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := usecase.NewChatService(nil, nil, nil)

	invalidType := model.ChatType("invalid")

	_, err := service.CreateChat(context.Background(), 123, invalidType)
	if err == nil {
		t.Fatal("expected error for invalid chat type, got nil")
	}
}

func TestGetChatByID_Success(t *testing.T) {
	ctrl, mockRepo, service := newChatServiceWithMocks(t)
	defer ctrl.Finish()

	expectedChat := &model.Chat{
		ID:        1,
		CreatorID: 123,
		Type:      model.ChatTypePrivate,
		CreatedAt: time.Now().Unix(),
	}

	mockRepo.EXPECT().
		FindByID(gomock.Any(), int64(1)).
		Return(expectedChat, nil)

	chat, err := service.GetChatByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if chat.ID != 1 {
		t.Fatalf("expected chat ID 1, got %d", chat.ID)
	}
}

func TestGetChatByID_NotFound(t *testing.T) {
	ctrl, mockRepo, service := newChatServiceWithMocks(t)
	defer ctrl.Finish()

	mockRepo.EXPECT().
		FindByID(gomock.Any(), int64(1)).
		Return(nil, ErrNotFound)

	_, err := service.GetChatByID(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRequestAdvice_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChatRepo := mocks.NewMockChatRepository(ctrl)
	mockBindingRepo := mocks.NewMockChatBindingRepository(ctrl)
	mockAdvice := mocks.NewMockAdvicePublisher(ctrl)

	binding := &model.ChatBinding{
		ChatID:   1,
		UserID:   123,
		Type:     model.AIBindingAdvice,
		ThreadID: "thread_abc123",
	}

	mockBindingRepo.EXPECT().
		FindByUserAndChat(gomock.Any(), int64(123), int64(1)).
		Return(binding, nil)

	mockAdvice.EXPECT().
		PublishAdviceRequest("thread_abc123").
		Return(nil)

	service := usecase.NewChatService(mockChatRepo, mockBindingRepo, mockAdvice)

	err := service.RequestAdvice(context.Background(), 123, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRequestAdvice_BindingNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChatRepo := mocks.NewMockChatRepository(ctrl)
	mockBindingRepo := mocks.NewMockChatBindingRepository(ctrl)
	mockAdvice := mocks.NewMockAdvicePublisher(ctrl)

	mockBindingRepo.EXPECT().
		FindByUserAndChat(gomock.Any(), int64(123), int64(1)).
		Return(nil, errors.New("not found"))

	service := usecase.NewChatService(mockChatRepo, mockBindingRepo, mockAdvice)

	err := service.RequestAdvice(context.Background(), 123, 1)
	if err == nil {
		t.Fatal("expected error for missing binding, got nil")
	}
}

func TestRequestAdvice_WrongType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChatRepo := mocks.NewMockChatRepository(ctrl)
	mockBindingRepo := mocks.NewMockChatBindingRepository(ctrl)
	mockAdvice := mocks.NewMockAdvicePublisher(ctrl)

	binding := &model.ChatBinding{
		ChatID:   1,
		UserID:   123,
		Type:     model.AIBindingAutoreply, // ⛔ не advice
		ThreadID: "thread_abc123",
	}

	mockBindingRepo.EXPECT().
		FindByUserAndChat(gomock.Any(), int64(123), int64(1)).
		Return(binding, nil)

	service := usecase.NewChatService(mockChatRepo, mockBindingRepo, mockAdvice)

	err := service.RequestAdvice(context.Background(), 123, 1)
	if err == nil {
		t.Fatal("expected error for wrong binding type, got nil")
	}
}

func TestRequestAdvice_PublishFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChatRepo := mocks.NewMockChatRepository(ctrl)
	mockBindingRepo := mocks.NewMockChatBindingRepository(ctrl)
	mockAdvice := mocks.NewMockAdvicePublisher(ctrl)

	binding := &model.ChatBinding{
		ChatID:   1,
		UserID:   123,
		Type:     model.AIBindingAdvice,
		ThreadID: "thread_abc123",
	}

	mockBindingRepo.EXPECT().
		FindByUserAndChat(gomock.Any(), int64(123), int64(1)).
		Return(binding, nil)

	mockAdvice.EXPECT().
		PublishAdviceRequest("thread_abc123").
		Return(errors.New("kafka error"))

	service := usecase.NewChatService(mockChatRepo, mockBindingRepo, mockAdvice)

	err := service.RequestAdvice(context.Background(), 123, 1)
	if err == nil {
		t.Fatal("expected error from publisher, got nil")
	}
}
