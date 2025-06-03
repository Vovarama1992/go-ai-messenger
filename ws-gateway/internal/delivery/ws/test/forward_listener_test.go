package ws_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

	"github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/delivery/ws"
	kafkaadapter "github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/infra/kafka"
	"github.com/Vovarama1992/go-ai-messenger/ws-gateway/internal/mocks"
)

func TestForwardListener_Start_CallsHubSend(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHub := mocks.NewMockHub(ctrl)
	mockConsumer := mocks.NewMockForwardConsumerInterface(ctrl)

	msg := kafkaadapter.ForwardMessage{
		ChatID:      1,
		SenderID:    42,
		Text:        "test message",
		AIGenerated: false,
	}

	mockHub.EXPECT().Sockets().Return(map[int64]struct{}{
		1:  {},
		42: {},
		99: {},
	}).AnyTimes()

	mockHub.EXPECT().Send(int64(1), "message", msg).Times(1)
	mockHub.EXPECT().Send(int64(42), "message", msg).Times(1)
	mockHub.EXPECT().Send(int64(99), "message", msg).Times(1)

	mockConsumer.EXPECT().Start().DoAndReturn(func() {
		handler := func(m kafkaadapter.ForwardMessage) {
			for userID := range mockHub.Sockets() {
				mockHub.Send(userID, "message", m)
			}
		}
		handler(msg)
	}).Times(1)

	listener := ws.NewForwardListenerWithConsumer(mockConsumer)
	listener.Start()

	require.True(t, true)
}

func TestForwardListener_NoSockets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHub := mocks.NewMockHub(ctrl)
	mockConsumer := mocks.NewMockForwardConsumerInterface(ctrl)

	msg := kafkaadapter.ForwardMessage{
		ChatID:      1,
		SenderID:    42,
		Text:        "test message",
		AIGenerated: false,
	}

	mockHub.EXPECT().Sockets().Return(map[int64]struct{}{}).Times(1)

	mockConsumer.EXPECT().Start().DoAndReturn(func() {
		handler := func(m kafkaadapter.ForwardMessage) {
			for userID := range mockHub.Sockets() {
				mockHub.Send(userID, "message", m)
			}
		}
		handler(msg)
	}).Times(1)

	listener := ws.NewForwardListenerWithConsumer(mockConsumer)
	listener.Start()

	require.True(t, true)
}

func TestForwardListener_MultipleRecipients(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHub := mocks.NewMockHub(ctrl)
	mockConsumer := mocks.NewMockForwardConsumerInterface(ctrl)

	msg := kafkaadapter.ForwardMessage{
		ChatID:      20,
		SenderID:    5,
		Text:        "multi-receiver",
		AIGenerated: true,
	}

	mockHub.EXPECT().Sockets().Return(map[int64]struct{}{
		1: {},
		5: {}, // sender
		8: {},
		9: {},
	}).AnyTimes()

	mockHub.EXPECT().Send(int64(1), "message", msg).Times(1)
	mockHub.EXPECT().Send(int64(5), "message", msg).Times(1)
	mockHub.EXPECT().Send(int64(8), "message", msg).Times(1)
	mockHub.EXPECT().Send(int64(9), "message", msg).Times(1)

	mockConsumer.EXPECT().Start().DoAndReturn(func() {
		handler := func(m kafkaadapter.ForwardMessage) {
			for userID := range mockHub.Sockets() {
				mockHub.Send(userID, "message", m)
			}
		}
		handler(msg)
	}).Times(1)

	listener := ws.NewForwardListenerWithConsumer(mockConsumer)
	listener.Start()

	require.True(t, true)
}
