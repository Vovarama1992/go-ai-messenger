package ws

import (
	kafkaadapter "github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/infra/kafka"
	"github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/ports"
)

type ForwardListener struct {
	consumer kafkaadapter.ForwardConsumerInterface
}

func NewForwardListener(hub *Hub, chatService ports.ChatService, topic string) *ForwardListener {
	consumer := kafkaadapter.NewForwardConsumer(topic, func(msg kafkaadapter.ForwardMessage) {

		hub.SendToRoom(msg.ChatID, "message", msg)
	})

	return &ForwardListener{
		consumer: consumer,
	}
}

func NewForwardListenerWithConsumer(consumer kafkaadapter.ForwardConsumerInterface) *ForwardListener {
	return &ForwardListener{
		consumer: consumer,
	}
}

func (fl *ForwardListener) Start() {
	fl.consumer.Start()
}

func (fl *ForwardListener) Stop() {
	fl.consumer.Stop()
}
