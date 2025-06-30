package ws

import (
	kafkaadapter "github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/infra/kafka"
)

type ForwardListener struct {
	consumer kafkaadapter.ForwardConsumerInterface
}

func NewForwardListener(hub *Hub, topic string) *ForwardListener {
	consumer := kafkaadapter.NewForwardConsumer(topic, func(msg kafkaadapter.ForwardMessage) {
		for userID := range hub.sockets {
			hub.Send(userID, "message", msg)
		}
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
