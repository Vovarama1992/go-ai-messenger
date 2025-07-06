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

func newChatServiceWithMocks(t *testing.T) (
	*gomock.Controller,
	*mocks.MockChatRepository,
	*usecase.ChatService,
) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockChatRepository(ctrl)
	mockBindingRepo := mocks.NewMockChatBindingRepository(ctrl)
	mockAdvicePublisher := mocks.NewMockAdvicePublisher(ctrl)
	mockBroker := mocks.NewMockMessageBroker(ctrl)

	service := usecase.NewChatService(mockBroker, mockRepo, mockBindingRepo, mockAdvicePublisher)

	return ctrl, mockRepo, service
}

func TestCreateChat_Success(t *testing.T) {
	ctrl, mockRepo, service := newChatServiceWithMocks(t)
	defer ctrl.Finish()

	mockRepo.EXPECT().
		Create(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, chat *model.Chat, memberIDs []int64) error {
			chat.ID = 1
			chat.CreatedAt = time.Now().Unix()
			return nil
		})

	chatType := model.ChatTypePrivate
	chat, err := service.CreateChat(context.Background(), 123, chatType, []int64{123})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if chat.CreatorID != 123 {
		t.Fatalf("expected CreatorID 123, got %d", chat.CreatorID)
	}
	if chat.ChatType != chatType {
		t.Fatalf("expected Type %s, got %s", chatType, chat.ChatType)
	}
	if chat.ID == 0 {
		t.Fatal("expected chat ID to be set")
	}
}

func TestCreateChat_InvalidType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := usecase.NewChatService(nil, nil, nil, nil)

	invalidType := model.ChatType("invalid")

	_, err := service.CreateChat(context.Background(), 123, invalidType, []int64{123})
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
		ChatType:  model.ChatTypePrivate,
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
		ChatID:      1,
		UserID:      123,
		BindingType: model.AIBindingAdvice,
		ThreadID:    "thread_abc123",
	}

	mockBindingRepo.EXPECT().
		FindByUserAndChat(gomock.Any(), int64(123), int64(1)).
		Return(binding, nil)

	mockAdvice.EXPECT().
		PublishAdviceRequest("thread_abc123").
		Return(nil)

	service := usecase.NewChatService(nil, mockChatRepo, mockBindingRepo, mockAdvice)

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

	service := usecase.NewChatService(nil, mockChatRepo, mockBindingRepo, mockAdvice)

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
		ChatID:      1,
		UserID:      123,
		BindingType: model.AIBindingAutoreply,
		ThreadID:    "thread_abc123",
	}

	mockBindingRepo.EXPECT().
		FindByUserAndChat(gomock.Any(), int64(123), int64(1)).
		Return(binding, nil)

	service := usecase.NewChatService(nil, mockChatRepo, mockBindingRepo, mockAdvice)

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
		ChatID:      1,
		UserID:      123,
		BindingType: model.AIBindingAdvice,
		ThreadID:    "thread_abc123",
	}

	mockBindingRepo.EXPECT().
		FindByUserAndChat(gomock.Any(), int64(123), int64(1)).
		Return(binding, nil)

	mockAdvice.EXPECT().
		PublishAdviceRequest("thread_abc123").
		Return(errors.New("kafka error"))

	service := usecase.NewChatService(nil, mockChatRepo, mockBindingRepo, mockAdvice)

	err := service.RequestAdvice(context.Background(), 123, 1)
	if err == nil {
		t.Fatal("expected error from publisher, got nil")
	}
}

func TestSendInvite_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockChatRepository(ctrl)
	mockBindingRepo := mocks.NewMockChatBindingRepository(ctrl)
	mockAdvice := mocks.NewMockAdvicePublisher(ctrl)
	mockBroker := mocks.NewMockMessageBroker(ctrl)

	service := usecase.NewChatService(mockBroker, mockRepo, mockBindingRepo, mockAdvice)

	chatID := int64(1)
	userIDs := []int64{10, 20}
	topic := "invite-topic"

	mockRepo.EXPECT().
		SendInvite(gomock.Any(), chatID, userIDs).
		Return(nil)

	for _, uid := range userIDs {
		mockBroker.EXPECT().
			SendInvite(gomock.Any(), uid, topic).
			Return(nil)
	}

	err := service.SendInvite(context.Background(), chatID, userIDs, topic)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSendInvite_RepoFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockChatRepository(ctrl)
	mockBindingRepo := mocks.NewMockChatBindingRepository(ctrl)
	mockAdvice := mocks.NewMockAdvicePublisher(ctrl)
	mockBroker := mocks.NewMockMessageBroker(ctrl)

	service := usecase.NewChatService(mockBroker, mockRepo, mockBindingRepo, mockAdvice)

	mockRepo.EXPECT().
		SendInvite(gomock.Any(), int64(1), gomock.Any()).
		Return(errors.New("db error"))

	err := service.SendInvite(context.Background(), 1, []int64{10}, "topic")
	if err == nil {
		t.Fatal("expected error from SendInvite, got nil")
	}
}

func TestSendInvite_BrokerFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockChatRepository(ctrl)
	mockBindingRepo := mocks.NewMockChatBindingRepository(ctrl)
	mockAdvice := mocks.NewMockAdvicePublisher(ctrl)
	mockBroker := mocks.NewMockMessageBroker(ctrl)

	service := usecase.NewChatService(mockBroker, mockRepo, mockBindingRepo, mockAdvice)

	userIDs := []int64{10, 20}
	topic := "invite-topic"

	mockRepo.EXPECT().
		SendInvite(gomock.Any(), int64(1), userIDs).
		Return(nil)

	// Один из брокеров фейлится
	mockBroker.EXPECT().
		SendInvite(gomock.Any(), gomock.Any(), topic).
		Return(errors.New("kafka down"))

	mockBroker.EXPECT().
		SendInvite(gomock.Any(), gomock.Any(), topic).
		Return(nil)

	err := service.SendInvite(context.Background(), 1, userIDs, topic)
	if err != nil {
		t.Fatalf("should not fail even if broker fails: %v", err)
	}
}

func TestAcceptInvite_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockChatRepository(ctrl)
	mockBindingRepo := mocks.NewMockChatBindingRepository(ctrl)
	mockAdvice := mocks.NewMockAdvicePublisher(ctrl)
	mockBroker := mocks.NewMockMessageBroker(ctrl)

	service := usecase.NewChatService(mockBroker, mockRepo, mockBindingRepo, mockAdvice)

	chatID := int64(1)
	userID := int64(123)

	mockRepo.EXPECT().
		AcceptInvite(gomock.Any(), chatID, userID).
		Return(nil)

	err := service.AcceptInvite(context.Background(), chatID, userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAcceptInvite_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockChatRepository(ctrl)
	service := usecase.NewChatService(nil, mockRepo, nil, nil)

	mockRepo.EXPECT().
		AcceptInvite(gomock.Any(), int64(1), int64(123)).
		Return(errors.New("db error"))

	err := service.AcceptInvite(context.Background(), 1, 123)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetParticipants_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockChatRepository(ctrl)
	service := usecase.NewChatService(nil, mockRepo, nil, nil)

	expected := []int64{1, 2, 3}

	mockRepo.EXPECT().
		GetChatParticipants(gomock.Any(), int64(10)).
		Return(expected, nil)

	participants, err := service.GetParticipants(context.Background(), 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(participants) != len(expected) {
		t.Fatalf("expected %d participants, got %d", len(expected), len(participants))
	}
}

func TestGetParticipants_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockChatRepository(ctrl)
	service := usecase.NewChatService(nil, mockRepo, nil, nil)

	mockRepo.EXPECT().
		GetChatParticipants(gomock.Any(), int64(10)).
		Return(nil, errors.New("db error"))

	_, err := service.GetParticipants(context.Background(), 10)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetPendingInvites_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockChatRepository(ctrl)
	service := usecase.NewChatService(nil, mockRepo, nil, nil)

	expected := []model.Chat{
		{ID: 1, CreatorID: 10, ChatType: model.ChatTypePrivate, CreatedAt: time.Now().Unix()},
		{ID: 2, CreatorID: 20, ChatType: model.ChatTypeGroup, CreatedAt: time.Now().Unix()},
	}

	mockRepo.EXPECT().
		GetPendingInvites(gomock.Any(), int64(123)).
		Return(expected, nil)

	invites, err := service.GetPendingInvites(context.Background(), 123)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(invites) != len(expected) {
		t.Fatalf("expected %d invites, got %d", len(expected), len(invites))
	}
	if invites[0].ID != expected[0].ID || invites[1].CreatorID != expected[1].CreatorID {
		t.Fatal("returned invites don't match expected data")
	}
}

func TestGetPendingInvites_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockChatRepository(ctrl)
	service := usecase.NewChatService(nil, mockRepo, nil, nil)

	mockRepo.EXPECT().
		GetPendingInvites(gomock.Any(), int64(123)).
		Return(nil, errors.New("db error"))

	_, err := service.GetPendingInvites(context.Background(), 123)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
