package ws_test

import (
	"testing"

	"github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/delivery/ws"
	kafkaadapter "github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/infra/kafka"
	"github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/mocks"
	"go.uber.org/mock/gomock"
)

func TestForwardListener_HandleMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := mocks.NewMockHub(ctrl)
	chatService := mocks.NewMockChatService(ctrl)

	// Создаём слушатель, топик не важен, handler проверяем напрямую
	fl := ws.NewForwardListener(hub, chatService, "topic")

	// Подготовим тестовое сообщение
	msg := kafkaadapter.ForwardMessage{
		ChatID:   123,
		SenderID: 42,
		Text:     "test message",
	}

	// Ожидаем вызов hub.SendToRoom с нужными параметрами
	hub.EXPECT().
		SendToRoom(msg.ChatID, "message", gomock.Any()).
		Times(1)

	// Вызовем внутренний обработчик вручную
	fl.HandleMessage(msg)
}
