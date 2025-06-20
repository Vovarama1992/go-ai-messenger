// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/ports (interfaces: ChatRepository)
//
// Generated by this command:
//
//	mockgen github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/ports ChatRepository
//

// Package mock_ports is a generated GoMock package.
package mock_ports

import (
	context "context"
	reflect "reflect"

	model "github.com/Vovarama1992/go-ai-messenger/chat-service/internal/chat/model"
	gomock "go.uber.org/mock/gomock"
)

// MockChatRepository is a mock of ChatRepository interface.
type MockChatRepository struct {
	ctrl     *gomock.Controller
	recorder *MockChatRepositoryMockRecorder
	isgomock struct{}
}

// MockChatRepositoryMockRecorder is the mock recorder for MockChatRepository.
type MockChatRepositoryMockRecorder struct {
	mock *MockChatRepository
}

// NewMockChatRepository creates a new mock instance.
func NewMockChatRepository(ctrl *gomock.Controller) *MockChatRepository {
	mock := &MockChatRepository{ctrl: ctrl}
	mock.recorder = &MockChatRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChatRepository) EXPECT() *MockChatRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockChatRepository) Create(ctx context.Context, chat *model.Chat) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, chat)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockChatRepositoryMockRecorder) Create(ctx, chat any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockChatRepository)(nil).Create), ctx, chat)
}

// FindByID mocks base method.
func (m *MockChatRepository) FindByID(ctx context.Context, id int64) (*model.Chat, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", ctx, id)
	ret0, _ := ret[0].(*model.Chat)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockChatRepositoryMockRecorder) FindByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockChatRepository)(nil).FindByID), ctx, id)
}
