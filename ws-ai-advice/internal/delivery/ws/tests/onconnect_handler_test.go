package ws_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/delivery/ws"
	"github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/mocks"
)

func TestOnConnectHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHub := mocks.NewMockAdviceHub(ctrl)
	mockAuth := mocks.NewMockAuthClient(ctrl)
	mockConn := mocks.NewMockConn(ctrl)

	mockConn.EXPECT().GetToken().Return("valid-token").Times(1)
	mockAuth.EXPECT().ValidateToken(gomock.Any(), "valid-token").Return(int64(123), "user@example.com", nil).Times(1)
	mockHub.EXPECT().Register(int64(123), mockConn).Times(1)
	mockConn.EXPECT().SetContext(gomock.Any()).Times(1)
	mockConn.EXPECT().Emit("connected", gomock.Any()).Times(1)

	handler := ws.TestableOnConnectHandler(mockHub, mockAuth)
	err := handler(mockConn)
	require.NoError(t, err)
}

func TestOnConnectHandler_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHub := mocks.NewMockAdviceHub(ctrl)
	mockAuth := mocks.NewMockAuthClient(ctrl)
	mockConn := mocks.NewMockConn(ctrl)

	mockConn.EXPECT().GetToken().Return("invalid-token").Times(1)
	mockAuth.EXPECT().ValidateToken(gomock.Any(), "invalid-token").Return(int64(0), "", errors.New("unauthorized")).Times(1)

	handler := ws.TestableOnConnectHandler(mockHub, mockAuth)
	err := handler(mockConn)
	require.Error(t, err)
}
