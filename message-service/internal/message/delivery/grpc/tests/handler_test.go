package messagegrpc_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	messagegrpc "github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/delivery/grpc"
	"github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/model"
	"github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/usecase"
	"github.com/Vovarama1992/go-ai-messenger/proto/messagepb"

	"github.com/Vovarama1992/go-ai-messenger/message-service/internal/message/mocks"
)

func TestHandler_GetMessagesByChat_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMessageRepo(ctrl)
	mockUser := mocks.NewMockUserClient(ctrl)

	svc := usecase.NewMessageService(mockRepo, mockUser)
	handler := messagegrpc.NewMessageHandler(svc)

	mockRepo.
		EXPECT().
		GetByChat(int64(100), 100, 0).
		Return([]model.Message{
			{
				ID:          1,
				ChatID:      100,
				SenderID:    42,
				Content:     "hello",
				AIGenerated: false,
				CreatedAt:   time.Now(),
			},
		}, nil)

	mockUser.
		EXPECT().
		GetUserEmailByID(gomock.Any(), int64(42)).
		Return("user@example.com", nil)

	resp, err := handler.GetMessagesByChat(context.Background(), &messagepb.GetMessagesRequest{ChatId: 100})
	require.NoError(t, err)
	require.Len(t, resp.Messages, 1)
	assert.Equal(t, "hello", resp.Messages[0].Content)
	assert.Equal(t, "user@example.com", resp.Messages[0].SenderEmail)
}

func TestHandler_GetMessagesByChat_ErrorFromRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMessageRepo(ctrl)
	mockUser := mocks.NewMockUserClient(ctrl)

	svc := usecase.NewMessageService(mockRepo, mockUser)
	handler := messagegrpc.NewMessageHandler(svc)

	mockRepo.
		EXPECT().
		GetByChat(int64(100), 100, 0).
		Return(nil, errors.New("db fail"))

	resp, err := handler.GetMessagesByChat(context.Background(), &messagepb.GetMessagesRequest{ChatId: 100})
	require.Nil(t, resp)
	assert.EqualError(t, err, "rpc error: code = Internal desc = db fail")
}
