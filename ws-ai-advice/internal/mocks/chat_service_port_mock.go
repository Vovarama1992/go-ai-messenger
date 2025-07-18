// Code generated by MockGen. DO NOT EDIT.
// Source: ws-ai-advice//internal/ports/chat_service_port.go
//
// Generated by this command:
//
//	mockgen -source=ws-ai-advice//internal/ports/chat_service_port.go -destination=ws-ai-advice//internal/mocks/chat_service_port_mock.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockChatService is a mock of ChatService interface.
type MockChatService struct {
	ctrl     *gomock.Controller
	recorder *MockChatServiceMockRecorder
	isgomock struct{}
}

// MockChatServiceMockRecorder is the mock recorder for MockChatService.
type MockChatServiceMockRecorder struct {
	mock *MockChatService
}

// NewMockChatService creates a new mock instance.
func NewMockChatService(ctrl *gomock.Controller) *MockChatService {
	mock := &MockChatService{ctrl: ctrl}
	mock.recorder = &MockChatServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChatService) EXPECT() *MockChatServiceMockRecorder {
	return m.recorder
}

// GetThreadContext mocks base method.
func (m *MockChatService) GetThreadContext(threadID string) (int64, string, int64, []int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetThreadContext", threadID)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(int64)
	ret3, _ := ret[3].([]int64)
	ret4, _ := ret[4].(error)
	return ret0, ret1, ret2, ret3, ret4
}

// GetThreadContext indicates an expected call of GetThreadContext.
func (mr *MockChatServiceMockRecorder) GetThreadContext(threadID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetThreadContext", reflect.TypeOf((*MockChatService)(nil).GetThreadContext), threadID)
}

// GetUserWithChatByThreadID mocks base method.
func (m *MockChatService) GetUserWithChatByThreadID(threadID string) (int64, int64, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserWithChatByThreadID", threadID)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(string)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// GetUserWithChatByThreadID indicates an expected call of GetUserWithChatByThreadID.
func (mr *MockChatServiceMockRecorder) GetUserWithChatByThreadID(threadID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserWithChatByThreadID", reflect.TypeOf((*MockChatService)(nil).GetUserWithChatByThreadID), threadID)
}
