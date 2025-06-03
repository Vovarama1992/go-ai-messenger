package ws_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/delivery/ws"
	"github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/mocks"
)

func TestHub_RegisterAndSend(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h := ws.NewHub()
	mockConn := mocks.NewMockConn(ctrl)

	mockConn.EXPECT().
		Emit("test_event", "payload").
		Times(1)

	h.Register(1, mockConn)
	h.Send(1, "test_event", "payload")

	require.True(t, true)
}

func TestHub_Unregister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h := ws.NewHub()
	mockConn := mocks.NewMockConn(ctrl)

	h.Register(1, mockConn)
	h.Unregister(1)

	// Send после Unregister не должен вызывать Emit
	mockConn.EXPECT().
		Emit("test_event", "payload").
		Times(0)

	h.Send(1, "test_event", "payload")
}
