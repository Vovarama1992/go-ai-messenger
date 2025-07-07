package ws

import (
	"log"

	kafkaadapter "github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/infra/kafka"
	"github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/ports"
)

type InviteListener struct {
	Consumer kafkaadapter.InviteConsumerInterface
}

func NewInviteListener(hub ports.Hub, topic string) *InviteListener {
	consumer := kafkaadapter.NewInviteConsumer(topic, func(msg kafkaadapter.InviteMessage) {
		if hub.HasConnection(msg.UserID) {
			conn := hub.GetConn(msg.UserID)
			if conn != nil {
				conn.Emit("invite", map[string]interface{}{
					"chatId": msg.ChatID,
				})
			}
		} else {
			log.Printf("ℹ️ user %d offline, skipping invite push", msg.UserID)
		}
	})

	return &InviteListener{Consumer: consumer}
}

func (il *InviteListener) Start() {
	il.Consumer.Start()
}

func (il *InviteListener) Stop() {
	il.Consumer.Stop()
}
