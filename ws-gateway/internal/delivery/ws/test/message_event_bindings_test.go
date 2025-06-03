package ws_test

import (
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	socketio "github.com/googollee/go-socket.io"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/delivery/ws"
	"github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/mocks"
	"github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/ports"
)

func TestSocket_Message_Bindings_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMock := mocks.NewMockAuthService(ctrl)
	chatMock := mocks.NewMockChatService(ctrl)
	kafkaMock := mocks.NewMockKafkaProducer(ctrl)
	hub := ws.NewHub()

	server := socketio.NewServer(nil)
	ws.RegisterSocketHandlers(server, authMock, chatMock, kafkaMock, hub)
	srv := httptest.NewServer(server)
	defer srv.Close()

	authMock.EXPECT().
		ValidateToken(gomock.Any(), "valid-token").
		Return(int64(123), nil).
		Times(1)

	chatMock.EXPECT().
		GetBindingsByChat(gomock.Any(), int64(1)).
		Return(nil, errors.New("chat service error")).
		Times(1)

	wsURL := "ws" + srv.URL[len("http"):] + "/socket.io/?EIO=4&transport=websocket&token=valid-token"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	msg := `42["message",{"chatId":1,"text":"hello"}]`
	err = conn.WriteMessage(websocket.TextMessage, []byte(msg))
	require.NoError(t, err)

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, resp, err := conn.ReadMessage()
	require.NoError(t, err)

	require.Contains(t, string(resp), "internal error")
}

func TestSocket_Message_TriggersAdviceAndAutoreplyKafkaProduce(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMock := mocks.NewMockAuthService(ctrl)
	chatMock := mocks.NewMockChatService(ctrl)
	kafkaMock := mocks.NewMockKafkaProducer(ctrl)
	hub := ws.NewHub()

	server := socketio.NewServer(nil)
	ws.RegisterSocketHandlers(server, authMock, chatMock, kafkaMock, hub)
	srv := httptest.NewServer(server)
	defer srv.Close()

	authMock.EXPECT().
		ValidateToken(gomock.Any(), "valid-token").
		Return(int64(123), nil).
		Times(1)

	chatMock.EXPECT().
		GetBindingsByChat(gomock.Any(), int64(1)).
		Return([]ports.ChatBinding{
			{UserID: 100, Type: "advice"},
			{UserID: 101, Type: "autoreply"},
		}, nil).
		Times(1)

	kafkaMock.EXPECT().
		Produce(gomock.Any(), "chat.message.persist", gomock.Any()).Times(1)

	kafkaMock.EXPECT().
		Produce(gomock.Any(), "chat.message.ai.advice-request", gomock.Any()).Times(1)

	kafkaMock.EXPECT().
		Produce(gomock.Any(), "chat.message.ai.autoreply-request", gomock.Any()).Times(1)

	wsURL := "ws" + srv.URL[len("http"):] + "/socket.io/?EIO=4&transport=websocket&token=valid-token"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	msg := `42["message",{"chatId":1,"text":"hello"}]`
	err = conn.WriteMessage(websocket.TextMessage, []byte(msg))
	require.NoError(t, err)

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, _, _ = conn.ReadMessage()

	require.True(t, true)
}
