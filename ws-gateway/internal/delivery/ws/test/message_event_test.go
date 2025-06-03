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
)

func TestMessageEvent_MissingToken(t *testing.T) {
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

	wsURL := "ws" + srv.URL[len("http"):] + "/socket.io/?EIO=4&transport=websocket"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	msg := `42["message",{"chatId":1,"text":"hello"}]`
	err = conn.WriteMessage(websocket.TextMessage, []byte(msg))
	require.NoError(t, err)

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, resp, err := conn.ReadMessage()
	require.NoError(t, err)

	require.Contains(t, string(resp), "missing token")
}

func TestMessageEvent_InvalidToken(t *testing.T) {
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
		ValidateToken(gomock.Any(), "bad-token").
		Return(int64(0), errors.New("unauthorized")).
		Times(1)

	wsURL := "ws" + srv.URL[len("http"):] + "/socket.io/?EIO=4&transport=websocket&token=bad-token"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	msg := `42["message",{"chatId":1,"text":"hello"}]`
	err = conn.WriteMessage(websocket.TextMessage, []byte(msg))
	require.NoError(t, err)

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, resp, err := conn.ReadMessage()
	require.NoError(t, err)

	require.Contains(t, string(resp), "unauthorized")
}
