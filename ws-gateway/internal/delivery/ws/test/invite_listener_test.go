package ws

import (
	"testing"

	kafkaadapter "github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/infra/kafka"
	fake "github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/infra/kafka/fakes"
	"github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/mocks"
	"go.uber.org/mock/gomock"
)

func TestInviteListener_HandleMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hub := mocks.NewMockHub(ctrl)
	conn := mocks.NewMockConn(ctrl)

	mockConsumer := &fake.MockInviteConsumer{}

	testUserID := int64(42)
	testData := map[string]interface{}{"chatId": int64(100)}

	hub.EXPECT().HasConnection(testUserID).Return(true)
	hub.EXPECT().GetConn(testUserID).Return(conn)
	conn.EXPECT().Emit("invite", testData).Return(nil)

	mockConsumer.SetHandler(func(msg kafkaadapter.InviteMessage) {
		if hub.HasConnection(msg.UserID) {
			conn := hub.GetConn(msg.UserID)
			if conn != nil {
				conn.Emit("invite", map[string]interface{}{
					"chatId": msg.ChatID,
				})
			}
		}
	})

	mockConsumer.Trigger(kafkaadapter.InviteMessage{
		UserID: testUserID,
		ChatID: 100,
	})
}
