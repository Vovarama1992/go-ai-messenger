package ws_test

import (
	"testing"

	"github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/delivery/ws"
	"github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHub_RegisterAndSend(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h := ws.NewHub()
	conn := mocks.NewMockConn(ctrl)

	conn.EXPECT().
		Emit("some_event", "data").
		Times(1)

	h.Register(1, conn)
	h.SendToRoom(123, "some_event", "data") // нет в комнате → ничего

	h.JoinRoom(1, 123)
	h.SendToRoom(123, "some_event", "data")
}

func TestHub_Unregister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h := ws.NewHub()
	conn := mocks.NewMockConn(ctrl)

	h.Register(1, conn)
	h.JoinRoom(1, 10)
	h.Unregister(1)

	require.False(t, h.HasConnection(1))

	// после удаления — вызовов быть не должно
	h.SendToRoom(10, "e", "d")
}

func TestHub_JoinAndLeaveRoom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h := ws.NewHub()
	h.JoinRoom(1, 10)
	h.JoinRoom(2, 10)

	h.LeaveRoom(1, 10)
	h.LeaveRoom(2, 10)

	conn1 := mocks.NewMockConn(ctrl)
	conn2 := mocks.NewMockConn(ctrl)

	h.Register(1, conn1)
	h.Register(2, conn2)

	// не должны срабатывать, так как из комнаты вышли
	h.SendToRoom(10, "x", "y")
}

func TestHub_HasConnection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h := ws.NewHub()
	require.False(t, h.HasConnection(123))

	conn := mocks.NewMockConn(ctrl)
	h.Register(123, conn)

	require.True(t, h.HasConnection(123))
}
