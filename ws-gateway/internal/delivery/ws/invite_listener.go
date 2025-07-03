package ws

import (
	"log"

	kafkaadapter "github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/infra/kafka"
)

type InviteListener struct {
	consumer kafkaadapter.InviteConsumerInterface
}

func NewInviteListener(hub *Hub, topic string) *InviteListener {
	consumer := kafkaadapter.NewInviteConsumer(topic, func(msg kafkaadapter.InviteMessage) {
		if hub.HasConnection(msg.UserID) {
			hub.sockets[msg.UserID].Emit("invite", map[string]interface{}{
				"chatId": msg.ChatID,
			})
		} else {
			log.Printf("ℹ️ user %d offline, skipping invite push", msg.UserID)
		}
	})

	return &InviteListener{consumer: consumer}
}

func (il *InviteListener) Start() {
	il.consumer.Start()
}

func (il *InviteListener) Stop() {
	il.consumer.Stop()
}
