// Code generated by MockGen. DO NOT EDIT.
// Source: ai-service//internal/ports/gpt_client.go
//
// Generated by this command:
//
//	mockgen -source=ai-service//internal/ports/gpt_client.go -destination=ai-service//internal/mocks/gpt_client_mock.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	dto "github.com/Vovarama1992/go-ai-messenger/ai-service/internal/dto"
	gomock "go.uber.org/mock/gomock"
)

// MockGptClient is a mock of GptClient interface.
type MockGptClient struct {
	ctrl     *gomock.Controller
	recorder *MockGptClientMockRecorder
	isgomock struct{}
}

// MockGptClientMockRecorder is the mock recorder for MockGptClient.
type MockGptClientMockRecorder struct {
	mock *MockGptClient
}

// NewMockGptClient creates a new mock instance.
func NewMockGptClient(ctrl *gomock.Controller) *MockGptClient {
	mock := &MockGptClient{ctrl: ctrl}
	mock.recorder = &MockGptClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGptClient) EXPECT() *MockGptClientMockRecorder {
	return m.recorder
}

// CreateThreadForUserAndChat mocks base method.
func (m *MockGptClient) CreateThreadForUserAndChat(ctx context.Context, userEmail string, messages []dto.ChatMessage) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateThreadForUserAndChat", ctx, userEmail, messages)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateThreadForUserAndChat indicates an expected call of CreateThreadForUserAndChat.
func (mr *MockGptClientMockRecorder) CreateThreadForUserAndChat(ctx, userEmail, messages any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateThreadForUserAndChat", reflect.TypeOf((*MockGptClient)(nil).CreateThreadForUserAndChat), ctx, userEmail, messages)
}

// GetAdvice mocks base method.
func (m *MockGptClient) GetAdvice(ctx context.Context, threadID string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAdvice", ctx, threadID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAdvice indicates an expected call of GetAdvice.
func (mr *MockGptClientMockRecorder) GetAdvice(ctx, threadID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAdvice", reflect.TypeOf((*MockGptClient)(nil).GetAdvice), ctx, threadID)
}

// SendMessageAndGetAutoreply mocks base method.
func (m *MockGptClient) SendMessageAndGetAutoreply(ctx context.Context, threadID, senderEmail, text string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessageAndGetAutoreply", ctx, threadID, senderEmail, text)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendMessageAndGetAutoreply indicates an expected call of SendMessageAndGetAutoreply.
func (mr *MockGptClientMockRecorder) SendMessageAndGetAutoreply(ctx, threadID, senderEmail, text any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessageAndGetAutoreply", reflect.TypeOf((*MockGptClient)(nil).SendMessageAndGetAutoreply), ctx, threadID, senderEmail, text)
}

// SendMessageToThread mocks base method.
func (m *MockGptClient) SendMessageToThread(ctx context.Context, threadID, role, content string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessageToThread", ctx, threadID, role, content)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMessageToThread indicates an expected call of SendMessageToThread.
func (mr *MockGptClientMockRecorder) SendMessageToThread(ctx, threadID, role, content any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessageToThread", reflect.TypeOf((*MockGptClient)(nil).SendMessageToThread), ctx, threadID, role, content)
}
