package kafka

import (
	kafkaadapter "github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/infra/kafka"
)

type MockInviteConsumer struct {
	handler func(kafkaadapter.InviteMessage)
}

func (m *MockInviteConsumer) Start() {}
func (m *MockInviteConsumer) Stop()  {}
func (m *MockInviteConsumer) SetHandler(h func(kafkaadapter.InviteMessage)) {
	m.handler = h
}
func (m *MockInviteConsumer) Trigger(msg kafkaadapter.InviteMessage) {
	if m.handler != nil {
		m.handler(msg)
	}
}
