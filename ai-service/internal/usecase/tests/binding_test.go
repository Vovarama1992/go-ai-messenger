package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/dto"
	mocks "github.com/Vovarama1992/go-ai-messenger/ai-service/internal/mocks"
	"github.com/Vovarama1992/go-ai-messenger/ai-service/internal/usecase"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestProcessBindingInit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGpt := mocks.NewMockGptClient(ctrl)

	payload := dto.AiBindingInitPayload{
		UserEmail: "test@example.com",
		UserID:    1,
		ChatID:    2,
		Messages: []dto.ChatMessage{
			{SenderEmail: "a@b", Text: "hello"},
		},
	}

	mockGpt.
		EXPECT().
		CreateThreadForUserAndChat(gomock.Any(), payload.UserEmail, payload.Messages).
		Return("thread-xyz", nil)

	res, err := usecase.ProcessBindingInit(context.Background(), payload, mockGpt)
	require.NoError(t, err)
	require.Equal(t, "thread-xyz", res.ThreadID)
	require.Equal(t, payload.ChatID, res.ChatID)
	require.Equal(t, payload.UserID, res.UserID)
}

func TestProcessBindingInit_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGpt := mocks.NewMockGptClient(ctrl)

	payload := dto.AiBindingInitPayload{
		UserEmail: "test@example.com",
		UserID:    1,
		ChatID:    2,
		Messages:  []dto.ChatMessage{},
	}

	mockGpt.
		EXPECT().
		CreateThreadForUserAndChat(gomock.Any(), payload.UserEmail, payload.Messages).
		Return("", errors.New("fail"))

	res, err := usecase.ProcessBindingInit(context.Background(), payload, mockGpt)
	require.Error(t, err)
	require.Empty(t, res.ThreadID)
}
