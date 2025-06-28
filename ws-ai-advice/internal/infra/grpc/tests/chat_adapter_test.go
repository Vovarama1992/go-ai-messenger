package grpcadapter_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/Vovarama1992/go-ai-messenger/proto/chatpb"
	grpcadapter "github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/infra/grpc"
	"github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/mocks"
)

func TestChatServiceAdapter_GetUserWithChatByThreadID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockChatServiceClient(ctrl)
	service := grpcadapter.NewChatService(mockClient)

	t.Run("success", func(t *testing.T) {
		mockClient.EXPECT().
			GetUserWithChatByThreadID(context.Background(), &chatpb.GetUserWithChatByThreadIDRequest{
				ThreadId: "abc",
			}).
			Return(&chatpb.GetUserWithChatByThreadIDResponse{
				UserId:    123,
				ChatId:    456,
				UserEmail: "user@example.com",
			}, nil)

		uid, cid, email, err := service.GetUserWithChatByThreadID("abc")
		require.NoError(t, err)
		require.Equal(t, int64(123), uid)
		require.Equal(t, int64(456), cid)
		require.Equal(t, "user@example.com", email)
	})

	t.Run("grpc returns error", func(t *testing.T) {
		mockClient.EXPECT().
			GetUserWithChatByThreadID(context.Background(), &chatpb.GetUserWithChatByThreadIDRequest{
				ThreadId: "fail",
			}).
			Return(nil, errors.New("boom"))

		_, _, _, err := service.GetUserWithChatByThreadID("fail")
		require.Error(t, err)
		require.Contains(t, err.Error(), "boom")
	})

	t.Run("grpc returns nil, nil", func(t *testing.T) {
		mockClient.EXPECT().
			GetUserWithChatByThreadID(context.Background(), &chatpb.GetUserWithChatByThreadIDRequest{
				ThreadId: "nil-case",
			}).
			Return(nil, nil)

		_, _, _, err := service.GetUserWithChatByThreadID("nil-case")
		require.Error(t, err)
		require.Contains(t, err.Error(), "nil response from chat service")
	})

	t.Run("long threadID", func(t *testing.T) {
		longID := make([]byte, 1024)
		for i := range longID {
			longID[i] = 'x'
		}

		threadID := string(longID)

		mockClient.EXPECT().
			GetUserWithChatByThreadID(context.Background(), &chatpb.GetUserWithChatByThreadIDRequest{
				ThreadId: threadID,
			}).
			Return(&chatpb.GetUserWithChatByThreadIDResponse{
				UserId:    1,
				ChatId:    2,
				UserEmail: "long@example.com",
			}, nil)

		uid, cid, email, err := service.GetUserWithChatByThreadID(threadID)
		require.NoError(t, err)
		require.Equal(t, int64(1), uid)
		require.Equal(t, int64(2), cid)
		require.Equal(t, "long@example.com", email)
	})
}
