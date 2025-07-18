// Code generated by MockGen. DO NOT EDIT.
// Source: proto/chatpb/chat_grpc.pb.go
//
// Generated by this command:
//
//	mockgen -source=proto/chatpb/chat_grpc.pb.go -destination=internal/mocks/chatpb_mock.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	chatpb "github.com/Vovarama1992/go-ai-messenger/proto/chatpb"
	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockChatServiceClient is a mock of ChatServiceClient interface.
type MockChatServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockChatServiceClientMockRecorder
	isgomock struct{}
}

// MockChatServiceClientMockRecorder is the mock recorder for MockChatServiceClient.
type MockChatServiceClientMockRecorder struct {
	mock *MockChatServiceClient
}

// NewMockChatServiceClient creates a new mock instance.
func NewMockChatServiceClient(ctrl *gomock.Controller) *MockChatServiceClient {
	mock := &MockChatServiceClient{ctrl: ctrl}
	mock.recorder = &MockChatServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChatServiceClient) EXPECT() *MockChatServiceClientMockRecorder {
	return m.recorder
}

// GetBindingsByChat mocks base method.
func (m *MockChatServiceClient) GetBindingsByChat(ctx context.Context, in *chatpb.GetBindingsByChatRequest, opts ...grpc.CallOption) (*chatpb.GetBindingsByChatResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetBindingsByChat", varargs...)
	ret0, _ := ret[0].(*chatpb.GetBindingsByChatResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBindingsByChat indicates an expected call of GetBindingsByChat.
func (mr *MockChatServiceClientMockRecorder) GetBindingsByChat(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBindingsByChat", reflect.TypeOf((*MockChatServiceClient)(nil).GetBindingsByChat), varargs...)
}

// GetChatByID mocks base method.
func (m *MockChatServiceClient) GetChatByID(ctx context.Context, in *chatpb.GetChatByIDRequest, opts ...grpc.CallOption) (*chatpb.GetChatByIDResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetChatByID", varargs...)
	ret0, _ := ret[0].(*chatpb.GetChatByIDResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChatByID indicates an expected call of GetChatByID.
func (mr *MockChatServiceClientMockRecorder) GetChatByID(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChatByID", reflect.TypeOf((*MockChatServiceClient)(nil).GetChatByID), varargs...)
}

// GetThreadContext mocks base method.
func (m *MockChatServiceClient) GetThreadContext(ctx context.Context, in *chatpb.GetThreadContextRequest, opts ...grpc.CallOption) (*chatpb.GetThreadContextResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetThreadContext", varargs...)
	ret0, _ := ret[0].(*chatpb.GetThreadContextResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetThreadContext indicates an expected call of GetThreadContext.
func (mr *MockChatServiceClientMockRecorder) GetThreadContext(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetThreadContext", reflect.TypeOf((*MockChatServiceClient)(nil).GetThreadContext), varargs...)
}

// GetUserWithChatByThreadID mocks base method.
func (m *MockChatServiceClient) GetUserWithChatByThreadID(ctx context.Context, in *chatpb.GetUserWithChatByThreadIDRequest, opts ...grpc.CallOption) (*chatpb.GetUserWithChatByThreadIDResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetUserWithChatByThreadID", varargs...)
	ret0, _ := ret[0].(*chatpb.GetUserWithChatByThreadIDResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserWithChatByThreadID indicates an expected call of GetUserWithChatByThreadID.
func (mr *MockChatServiceClientMockRecorder) GetUserWithChatByThreadID(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserWithChatByThreadID", reflect.TypeOf((*MockChatServiceClient)(nil).GetUserWithChatByThreadID), varargs...)
}

// GetUsersByChatID mocks base method.
func (m *MockChatServiceClient) GetUsersByChatID(ctx context.Context, in *chatpb.GetUsersByChatIDRequest, opts ...grpc.CallOption) (*chatpb.GetUsersByChatIDResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetUsersByChatID", varargs...)
	ret0, _ := ret[0].(*chatpb.GetUsersByChatIDResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsersByChatID indicates an expected call of GetUsersByChatID.
func (mr *MockChatServiceClientMockRecorder) GetUsersByChatID(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersByChatID", reflect.TypeOf((*MockChatServiceClient)(nil).GetUsersByChatID), varargs...)
}

// MockChatServiceServer is a mock of ChatServiceServer interface.
type MockChatServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockChatServiceServerMockRecorder
	isgomock struct{}
}

// MockChatServiceServerMockRecorder is the mock recorder for MockChatServiceServer.
type MockChatServiceServerMockRecorder struct {
	mock *MockChatServiceServer
}

// NewMockChatServiceServer creates a new mock instance.
func NewMockChatServiceServer(ctrl *gomock.Controller) *MockChatServiceServer {
	mock := &MockChatServiceServer{ctrl: ctrl}
	mock.recorder = &MockChatServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChatServiceServer) EXPECT() *MockChatServiceServerMockRecorder {
	return m.recorder
}

// GetBindingsByChat mocks base method.
func (m *MockChatServiceServer) GetBindingsByChat(arg0 context.Context, arg1 *chatpb.GetBindingsByChatRequest) (*chatpb.GetBindingsByChatResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBindingsByChat", arg0, arg1)
	ret0, _ := ret[0].(*chatpb.GetBindingsByChatResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBindingsByChat indicates an expected call of GetBindingsByChat.
func (mr *MockChatServiceServerMockRecorder) GetBindingsByChat(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBindingsByChat", reflect.TypeOf((*MockChatServiceServer)(nil).GetBindingsByChat), arg0, arg1)
}

// GetChatByID mocks base method.
func (m *MockChatServiceServer) GetChatByID(arg0 context.Context, arg1 *chatpb.GetChatByIDRequest) (*chatpb.GetChatByIDResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChatByID", arg0, arg1)
	ret0, _ := ret[0].(*chatpb.GetChatByIDResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChatByID indicates an expected call of GetChatByID.
func (mr *MockChatServiceServerMockRecorder) GetChatByID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChatByID", reflect.TypeOf((*MockChatServiceServer)(nil).GetChatByID), arg0, arg1)
}

// GetThreadContext mocks base method.
func (m *MockChatServiceServer) GetThreadContext(arg0 context.Context, arg1 *chatpb.GetThreadContextRequest) (*chatpb.GetThreadContextResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetThreadContext", arg0, arg1)
	ret0, _ := ret[0].(*chatpb.GetThreadContextResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetThreadContext indicates an expected call of GetThreadContext.
func (mr *MockChatServiceServerMockRecorder) GetThreadContext(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetThreadContext", reflect.TypeOf((*MockChatServiceServer)(nil).GetThreadContext), arg0, arg1)
}

// GetUserWithChatByThreadID mocks base method.
func (m *MockChatServiceServer) GetUserWithChatByThreadID(arg0 context.Context, arg1 *chatpb.GetUserWithChatByThreadIDRequest) (*chatpb.GetUserWithChatByThreadIDResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserWithChatByThreadID", arg0, arg1)
	ret0, _ := ret[0].(*chatpb.GetUserWithChatByThreadIDResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserWithChatByThreadID indicates an expected call of GetUserWithChatByThreadID.
func (mr *MockChatServiceServerMockRecorder) GetUserWithChatByThreadID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserWithChatByThreadID", reflect.TypeOf((*MockChatServiceServer)(nil).GetUserWithChatByThreadID), arg0, arg1)
}

// GetUsersByChatID mocks base method.
func (m *MockChatServiceServer) GetUsersByChatID(arg0 context.Context, arg1 *chatpb.GetUsersByChatIDRequest) (*chatpb.GetUsersByChatIDResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsersByChatID", arg0, arg1)
	ret0, _ := ret[0].(*chatpb.GetUsersByChatIDResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsersByChatID indicates an expected call of GetUsersByChatID.
func (mr *MockChatServiceServerMockRecorder) GetUsersByChatID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersByChatID", reflect.TypeOf((*MockChatServiceServer)(nil).GetUsersByChatID), arg0, arg1)
}

// mustEmbedUnimplementedChatServiceServer mocks base method.
func (m *MockChatServiceServer) mustEmbedUnimplementedChatServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedChatServiceServer")
}

// mustEmbedUnimplementedChatServiceServer indicates an expected call of mustEmbedUnimplementedChatServiceServer.
func (mr *MockChatServiceServerMockRecorder) mustEmbedUnimplementedChatServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedChatServiceServer", reflect.TypeOf((*MockChatServiceServer)(nil).mustEmbedUnimplementedChatServiceServer))
}

// MockUnsafeChatServiceServer is a mock of UnsafeChatServiceServer interface.
type MockUnsafeChatServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeChatServiceServerMockRecorder
	isgomock struct{}
}

// MockUnsafeChatServiceServerMockRecorder is the mock recorder for MockUnsafeChatServiceServer.
type MockUnsafeChatServiceServerMockRecorder struct {
	mock *MockUnsafeChatServiceServer
}

// NewMockUnsafeChatServiceServer creates a new mock instance.
func NewMockUnsafeChatServiceServer(ctrl *gomock.Controller) *MockUnsafeChatServiceServer {
	mock := &MockUnsafeChatServiceServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeChatServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeChatServiceServer) EXPECT() *MockUnsafeChatServiceServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedChatServiceServer mocks base method.
func (m *MockUnsafeChatServiceServer) mustEmbedUnimplementedChatServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedChatServiceServer")
}

// mustEmbedUnimplementedChatServiceServer indicates an expected call of mustEmbedUnimplementedChatServiceServer.
func (mr *MockUnsafeChatServiceServerMockRecorder) mustEmbedUnimplementedChatServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedChatServiceServer", reflect.TypeOf((*MockUnsafeChatServiceServer)(nil).mustEmbedUnimplementedChatServiceServer))
}
