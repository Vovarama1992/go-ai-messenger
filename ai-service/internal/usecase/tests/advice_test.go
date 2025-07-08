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

func TestProcessAdviceRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGpt := mocks.NewMockGptClient(ctrl)

	threadID := "thread-123"
	expectedText := "вот тебе совет"
	payload := dto.AdviceRequestPayload{ThreadID: threadID}

	mockGpt.
		EXPECT().
		GetAdvice(gomock.Any(), threadID).
		Return(expectedText, nil)

	resp, err := usecase.ProcessAdviceRequest(context.Background(), payload, mockGpt)
	require.NoError(t, err)
	require.Equal(t, threadID, resp.ThreadID)
	require.Equal(t, expectedText, resp.Text)
}

func TestProcessAdviceRequest_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGpt := mocks.NewMockGptClient(ctrl)
	threadID := "thread-456"
	payload := dto.AdviceRequestPayload{ThreadID: threadID}

	mockGpt.
		EXPECT().
		GetAdvice(gomock.Any(), threadID).
		Return("", errors.New("fail"))

	resp, err := usecase.ProcessAdviceRequest(context.Background(), payload, mockGpt)
	require.Error(t, err)
	require.Empty(t, resp.Text)
}
