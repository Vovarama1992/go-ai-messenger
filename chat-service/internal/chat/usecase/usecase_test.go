package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/mocks"
	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/model"
	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/usecase"
	"go.uber.org/mock/gomock"
)

var ErrNotFound = errors.New("chat not found")

func TestCreateChat_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockChatRepository(ctrl)
	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, chat *model.Chat) error {
			chat.ID = 1
			chat.CreatedAt = time.Now().Unix()
			return nil
		})

	service := usecase.NewChatService(mockRepo)

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
	service := usecase.NewChatService(nil)

	invalidType := model.ChatType("invalid")

	_, err := service.CreateChat(context.Background(), 123, invalidType)
	if err == nil {
		t.Fatal("expected error for invalid chat type, got nil")
	}
}

func TestGetChatByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockChatRepository(ctrl)
	expectedChat := &model.Chat{
		ID:        1,
		CreatorID: 123,
		Type:      model.ChatTypePrivate,
		CreatedAt: time.Now().Unix(),
	}

	mockRepo.EXPECT().
		FindByID(gomock.Any(), int64(1)).
		Return(expectedChat, nil)

	service := usecase.NewChatService(mockRepo)

	chat, err := service.GetChatByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if chat.ID != 1 {
		t.Fatalf("expected chat ID 1, got %d", chat.ID)
	}
}

func TestGetChatByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockChatRepository(ctrl)
	mockRepo.EXPECT().
		FindByID(gomock.Any(), int64(1)).
		Return(nil, ErrNotFound) // Определи ErrNotFound в usecase

	service := usecase.NewChatService(mockRepo)

	_, err := service.GetChatByID(context.Background(), 1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
