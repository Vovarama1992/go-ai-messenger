package ws

import (
	kafkaadapter "github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/infra/kafka"
	"github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/ports"
)

type ForwardListener struct {
	consumer kafkaadapter.ForwardConsumerInterface
	hub      ports.Hub
}

func NewForwardListener(hub ports.Hub, chatService ports.ChatService, topic string) *ForwardListener {
	fl := &ForwardListener{
		hub: hub,
	}
	consumer := kafkaadapter.NewForwardConsumer(topic, fl.HandleMessage)
	fl.consumer = consumer
	return fl
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

func (fl *ForwardListener) HandleMessage(msg kafkaadapter.ForwardMessage) {
	fl.hub.SendToRoom(msg.ChatID, "message", msg)
}
