package ws_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/delivery/ws"
	"github.com/Vovarama1992/go-ai-messenger/ws-ai-advice/internal/mocks"
	"go.uber.org/mock/gomock"
)

func TestOnConnectHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mocks.NewMockAuthClient(ctrl)
	mockHub := mocks.NewMockAdviceHub(ctrl)
	mockConn := mocks.NewMockConn(ctrl)

	mockConn.EXPECT().GetToken().Return("valid-token")
	mockAuth.EXPECT().
		ValidateToken(context.Background(), "valid-token").
		Return(int64(123), "user@example.com", nil)
	mockHub.EXPECT().Register(int64(123), mockConn)
	mockConn.EXPECT().SetContext(gomock.Any())
	mockConn.EXPECT().Emit("connected", gomock.Any())

	handler := ws.TestableOnConnectHandler(mockHub, mockAuth)
	err := handler(mockConn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestOnConnectHandler_MissingToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mocks.NewMockAuthClient(ctrl)
	mockHub := mocks.NewMockAdviceHub(ctrl)
	mockConn := mocks.NewMockConn(ctrl)

	mockConn.EXPECT().GetToken().Return("")

	handler := ws.TestableOnConnectHandler(mockHub, mockAuth)
	err := handler(mockConn)
	if err == nil || err.Error() != "missing token" {
		t.Fatalf("expected 'missing token', got: %v", err)
	}
}

func TestOnConnectHandler_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mocks.NewMockAuthClient(ctrl)
	mockHub := mocks.NewMockAdviceHub(ctrl)
	mockConn := mocks.NewMockConn(ctrl)

	mockConn.EXPECT().GetToken().Return("bad")
	mockAuth.EXPECT().
		ValidateToken(context.Background(), "bad").
		Return(int64(0), "", errors.New("unauthorized"))

	handler := ws.TestableOnConnectHandler(mockHub, mockAuth)
	err := handler(mockConn)
	if err == nil || err.Error() != "unauthorized: unauthorized" {
		t.Fatalf("expected auth error, got: %v", err)
	}
}
