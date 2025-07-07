package ws_test

import (
	"errors"
	"testing"

	"github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/delivery/ws"
	fake "github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/delivery/ws/fakes"
	"github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/mocks"
	"github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/ports"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestHandleConnect_ValidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	auth := mocks.NewMockAuthService(ctrl)
	hub := mocks.NewMockHub(ctrl)
	conn := mocks.NewMockConn(ctrl)

	// ожидаем вызовы
	auth.EXPECT().
		ValidateToken(gomock.Any(), "valid-token").
		Return(int64(42), "me@email.com", nil)

	hub.EXPECT().
		Register(int64(42), conn)

	conn.EXPECT().
		SetContext(gomock.Any())

	conn.EXPECT().
		Emit("connected", gomock.Any()).
		Return(nil)

	err := ws.HandleConnect(auth, hub, conn, "valid-token")
	require.NoError(t, err)
}

func TestHandleConnect_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	auth := mocks.NewMockAuthService(ctrl)
	hub := mocks.NewMockHub(ctrl)
	conn := mocks.NewMockConn(ctrl)

	auth.EXPECT().
		ValidateToken(gomock.Any(), "invalid-token").
		Return(int64(0), "", errors.New("unauthorized"))

	err := ws.HandleConnect(auth, hub, conn, "invalid-token")
	require.EqualError(t, err, "unauthorized")
}

func TestMakeDisconnectHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := mocks.NewMockHub(ctrl)

	fc := &fake.FakeConn{
		Ctx: struct {
			ID    int64
			Email string
		}{ID: 123, Email: "x"},
	}

	hub.EXPECT().Unregister(int64(123))

	handler := ws.MakeDisconnectHandler(hub)
	handler(fc, "client-left")
}

func TestMakeMessageHandler_ValidInput(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	chatService := mocks.NewMockChatService(ctrl)
	kafka := mocks.NewMockKafkaProducer(ctrl)

	conn := &fake.FakeConn{
		Ctx: ws.UserCtx{ID: 42, Email: "test@example.com"},
	}

	topicPersist := "topic-persist"
	topicFeed := "topic-feed"

	handler := ws.MakeMessageHandler(chatService, kafka, topicPersist, topicFeed)

	msg := map[string]interface{}{
		"chatId": float64(123),
		"text":   "Hello!",
	}

	chatService.EXPECT().
		GetBindingsByChat(gomock.Any(), int64(123)).
		Return([]ports.ChatBinding{
			{ThreadID: "abc", BindingType: "feed"},
		}, nil)

	kafka.EXPECT().
		Produce(gomock.Any(), topicPersist, gomock.Any())

	kafka.EXPECT().
		Produce(gomock.Any(), topicFeed, gomock.Any())

	handler(conn, msg)
}

func TestMakeDisconnectHandler_NoContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := mocks.NewMockHub(ctrl)
	fc := &fake.FakeConn{
		Ctx: nil, // Нет контекста
	}

	// НЕ ожидаем вызова Unregister
	// hub.EXPECT().Unregister(gomock.Any()).Times(0) // можно опустить, по умолчанию

	handler := ws.MakeDisconnectHandler(hub)
	handler(fc, "client-left")
}
func TestMakeMessageHandler_NoUserCtx(t *testing.T) {
	conn := &fake.FakeConnWithEmitLog{Ctx: nil}

	handler := ws.MakeMessageHandler(nil, nil, "", "")
	handler(conn, map[string]interface{}{
		"chatId": float64(1),
		"text":   "hello",
	})

	found := false
	for _, ev := range conn.Emitted {
		if ev.Event == "error" {
			found = true
			break
		}
	}
	require.True(t, found, "expected 'error' event to be emitted")
}

func TestMakeMessageHandler_InvalidChatID(t *testing.T) {
	conn := &fake.FakeConnWithEmitLog{
		Ctx: struct {
			ID    int64
			Email string
		}{ID: 1, Email: "a@b.com"},
	}

	handler := ws.MakeMessageHandler(nil, nil, "", "")

	handler(conn, map[string]interface{}{
		"chatId": "not-a-float",
		"text":   "hello",
	})

	found := false
	for _, ev := range conn.Emitted {
		if ev.Event == "error" && len(ev.Data) > 0 && ev.Data[0] == "invalid chatId" {
			found = true
			break
		}
	}

	require.True(t, found, "expected 'invalid chatId' error to be emitted")
}

func TestMakeMessageHandler_EmptyText(t *testing.T) {
	conn := &fake.FakeConnWithEmitLog{
		Ctx: struct {
			ID    int64
			Email string
		}{ID: 1, Email: "a@b.com"},
	}

	handler := ws.MakeMessageHandler(nil, nil, "", "")

	handler(conn, map[string]interface{}{
		"chatId": float64(123),
		"text":   "",
	})

	found := false
	for _, ev := range conn.Emitted {
		if ev.Event == "error" && len(ev.Data) > 0 && ev.Data[0] == "empty message" {
			found = true
			break
		}
	}

	require.True(t, found, "expected 'empty message' error to be emitted")
}
