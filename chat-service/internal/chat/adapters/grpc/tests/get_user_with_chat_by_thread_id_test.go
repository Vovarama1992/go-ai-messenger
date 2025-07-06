package grpc_test

import (
	"context"
	"errors"
	"testing"

	chatgrpc "github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/adapters/grpc"
	mock_ports "github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/mocks"
	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/model"
	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/ports"
	"github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/usecase"
	"github.com/Vovarama1992/go-ai-messenger/proto/chatpb"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type testEnv struct {
	handler        *chatgrpc.ChatHandler
	mockUserClient *mock_ports.MockUserClient
	mockBinding    *mock_ports.MockChatBindingRepository
}

func setupTest(t *testing.T) *testEnv {
	ctrl := gomock.NewController(t)

	mockUserClient := mock_ports.NewMockUserClient(ctrl)
	mockBindingRepo := mock_ports.NewMockChatBindingRepository(ctrl)

	bindingService := usecase.NewChatBindingService(mockBindingRepo, nil, nil)
	userService := usecase.NewUserService(mockUserClient)
	handler := chatgrpc.NewChatHandler(nil, bindingService, userService)

	return &testEnv{
		handler:        handler,
		mockUserClient: mockUserClient,
		mockBinding:    mockBindingRepo,
	}
}

func TestGetUserWithChatByThreadID(t *testing.T) {
	threadID := "test-thread-123"
	expectedBinding := &model.ChatBinding{
		UserID: 42,
		ChatID: 99,
	}
	expectedUser := &ports.User{
		ID:    42,
		Email: "user@example.com",
	}

	t.Run("success", func(t *testing.T) {
		env := setupTest(t)

		env.mockBinding.EXPECT().
			FindByThreadID(gomock.Any(), threadID).
			Return(expectedBinding, nil)

		env.mockUserClient.EXPECT().
			GetUserByID(gomock.Any(), int64(42)).
			Return(expectedUser, nil)

		req := &chatpb.GetUserWithChatByThreadIDRequest{ThreadId: threadID}
		resp, err := env.handler.GetUserWithChatByThreadID(context.Background(), req)

		require.NoError(t, err)
		require.Equal(t, int64(42), resp.UserId)
		require.Equal(t, int64(99), resp.ChatId)
		require.Equal(t, "user@example.com", resp.UserEmail)
	})

	t.Run("binding not found", func(t *testing.T) {
		env := setupTest(t)

		env.mockBinding.EXPECT().
			FindByThreadID(gomock.Any(), threadID).
			Return(nil, errors.New("not found"))

		req := &chatpb.GetUserWithChatByThreadIDRequest{ThreadId: threadID}
		_, err := env.handler.GetUserWithChatByThreadID(context.Background(), req)

		require.Error(t, err)
	})

	t.Run("user not found", func(t *testing.T) {
		env := setupTest(t)

		env.mockBinding.EXPECT().
			FindByThreadID(gomock.Any(), threadID).
			Return(expectedBinding, nil)

		env.mockUserClient.EXPECT().
			GetUserByID(gomock.Any(), int64(42)).
			Return(nil, errors.New("user not found"))

		req := &chatpb.GetUserWithChatByThreadIDRequest{ThreadId: threadID}
		_, err := env.handler.GetUserWithChatByThreadID(context.Background(), req)

		require.Error(t, err)
	})
}
