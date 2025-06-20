// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/ports (interfaces: AdvicePublisher)
//
// Generated by this command:
//
//	mockgen github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/ports AdvicePublisher
//

// Package mock_ports is a generated GoMock package.
package mock_ports

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockAdvicePublisher is a mock of AdvicePublisher interface.
type MockAdvicePublisher struct {
	ctrl     *gomock.Controller
	recorder *MockAdvicePublisherMockRecorder
	isgomock struct{}
}

// MockAdvicePublisherMockRecorder is the mock recorder for MockAdvicePublisher.
type MockAdvicePublisherMockRecorder struct {
	mock *MockAdvicePublisher
}

// NewMockAdvicePublisher creates a new mock instance.
func NewMockAdvicePublisher(ctrl *gomock.Controller) *MockAdvicePublisher {
	mock := &MockAdvicePublisher{ctrl: ctrl}
	mock.recorder = &MockAdvicePublisherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAdvicePublisher) EXPECT() *MockAdvicePublisherMockRecorder {
	return m.recorder
}

// PublishAdviceRequest mocks base method.
func (m *MockAdvicePublisher) PublishAdviceRequest(threadID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublishAdviceRequest", threadID)
	ret0, _ := ret[0].(error)
	return ret0
}

// PublishAdviceRequest indicates an expected call of PublishAdviceRequest.
func (mr *MockAdvicePublisherMockRecorder) PublishAdviceRequest(threadID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishAdviceRequest", reflect.TypeOf((*MockAdvicePublisher)(nil).PublishAdviceRequest), threadID)
}
