// Code generated by MockGen. DO NOT EDIT.
// Source: proto/service.pb.go

// Package mock_pl_dwojciechowski_proto is a generated GoMock package.
package mock_pl_dwojciechowski_proto

import (
	context "context"
	proto "dominikw.pl/wnc_plugin/proto"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
	metadata "google.golang.org/grpc/metadata"
	reflect "reflect"
)

// MockFileServiceClient is a mock of FileServiceClient interface
type MockFileServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockFileServiceClientMockRecorder
}

// MockFileServiceClientMockRecorder is the mock recorder for MockFileServiceClient
type MockFileServiceClientMockRecorder struct {
	mock *MockFileServiceClient
}

// NewMockFileServiceClient creates a new mock instance
func NewMockFileServiceClient(ctrl *gomock.Controller) *MockFileServiceClient {
	mock := &MockFileServiceClient{ctrl: ctrl}
	mock.recorder = &MockFileServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFileServiceClient) EXPECT() *MockFileServiceClientMockRecorder {
	return m.recorder
}

// Navigate mocks base method
func (m *MockFileServiceClient) Navigate(ctx context.Context, in *proto.Path, opts ...grpc.CallOption) (*proto.FileResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Navigate", varargs...)
	ret0, _ := ret[0].(*proto.FileResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Navigate indicates an expected call of Navigate
func (mr *MockFileServiceClientMockRecorder) Navigate(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Navigate", reflect.TypeOf((*MockFileServiceClient)(nil).Navigate), varargs...)
}

// MockFileServiceServer is a mock of FileServiceServer interface
type MockFileServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockFileServiceServerMockRecorder
}

// MockFileServiceServerMockRecorder is the mock recorder for MockFileServiceServer
type MockFileServiceServerMockRecorder struct {
	mock *MockFileServiceServer
}

// NewMockFileServiceServer creates a new mock instance
func NewMockFileServiceServer(ctrl *gomock.Controller) *MockFileServiceServer {
	mock := &MockFileServiceServer{ctrl: ctrl}
	mock.recorder = &MockFileServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFileServiceServer) EXPECT() *MockFileServiceServerMockRecorder {
	return m.recorder
}

// Navigate mocks base method
func (m *MockFileServiceServer) Navigate(arg0 context.Context, arg1 *proto.Path) (*proto.FileResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Navigate", arg0, arg1)
	ret0, _ := ret[0].(*proto.FileResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Navigate indicates an expected call of Navigate
func (mr *MockFileServiceServerMockRecorder) Navigate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Navigate", reflect.TypeOf((*MockFileServiceServer)(nil).Navigate), arg0, arg1)
}

// MockCommandServiceClient is a mock of CommandServiceClient interface
type MockCommandServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockCommandServiceClientMockRecorder
}

// MockCommandServiceClientMockRecorder is the mock recorder for MockCommandServiceClient
type MockCommandServiceClientMockRecorder struct {
	mock *MockCommandServiceClient
}

// NewMockCommandServiceClient creates a new mock instance
func NewMockCommandServiceClient(ctrl *gomock.Controller) *MockCommandServiceClient {
	mock := &MockCommandServiceClient{ctrl: ctrl}
	mock.recorder = &MockCommandServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCommandServiceClient) EXPECT() *MockCommandServiceClientMockRecorder {
	return m.recorder
}

// Execute mocks base method
func (m *MockCommandServiceClient) Execute(ctx context.Context, in *proto.Command, opts ...grpc.CallOption) (*proto.Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Execute", varargs...)
	ret0, _ := ret[0].(*proto.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute
func (mr *MockCommandServiceClientMockRecorder) Execute(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockCommandServiceClient)(nil).Execute), varargs...)
}

// MockCommandServiceServer is a mock of CommandServiceServer interface
type MockCommandServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockCommandServiceServerMockRecorder
}

// MockCommandServiceServerMockRecorder is the mock recorder for MockCommandServiceServer
type MockCommandServiceServerMockRecorder struct {
	mock *MockCommandServiceServer
}

// NewMockCommandServiceServer creates a new mock instance
func NewMockCommandServiceServer(ctrl *gomock.Controller) *MockCommandServiceServer {
	mock := &MockCommandServiceServer{ctrl: ctrl}
	mock.recorder = &MockCommandServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCommandServiceServer) EXPECT() *MockCommandServiceServerMockRecorder {
	return m.recorder
}

// Execute mocks base method
func (m *MockCommandServiceServer) Execute(arg0 context.Context, arg1 *proto.Command) (*proto.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", arg0, arg1)
	ret0, _ := ret[0].(*proto.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute
func (mr *MockCommandServiceServerMockRecorder) Execute(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockCommandServiceServer)(nil).Execute), arg0, arg1)
}

// MockLogViewerServiceClient is a mock of LogViewerServiceClient interface
type MockLogViewerServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockLogViewerServiceClientMockRecorder
}

// MockLogViewerServiceClientMockRecorder is the mock recorder for MockLogViewerServiceClient
type MockLogViewerServiceClientMockRecorder struct {
	mock *MockLogViewerServiceClient
}

// NewMockLogViewerServiceClient creates a new mock instance
func NewMockLogViewerServiceClient(ctrl *gomock.Controller) *MockLogViewerServiceClient {
	mock := &MockLogViewerServiceClient{ctrl: ctrl}
	mock.recorder = &MockLogViewerServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLogViewerServiceClient) EXPECT() *MockLogViewerServiceClientMockRecorder {
	return m.recorder
}

// GetLogs mocks base method
func (m *MockLogViewerServiceClient) GetLogs(ctx context.Context, in *proto.LogFileLocation, opts ...grpc.CallOption) (proto.LogViewerService_GetLogsClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetLogs", varargs...)
	ret0, _ := ret[0].(proto.LogViewerService_GetLogsClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLogs indicates an expected call of GetLogs
func (mr *MockLogViewerServiceClientMockRecorder) GetLogs(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLogs", reflect.TypeOf((*MockLogViewerServiceClient)(nil).GetLogs), varargs...)
}

// MockLogViewerService_GetLogsClient is a mock of LogViewerService_GetLogsClient interface
type MockLogViewerService_GetLogsClient struct {
	ctrl     *gomock.Controller
	recorder *MockLogViewerService_GetLogsClientMockRecorder
}

// MockLogViewerService_GetLogsClientMockRecorder is the mock recorder for MockLogViewerService_GetLogsClient
type MockLogViewerService_GetLogsClientMockRecorder struct {
	mock *MockLogViewerService_GetLogsClient
}

// NewMockLogViewerService_GetLogsClient creates a new mock instance
func NewMockLogViewerService_GetLogsClient(ctrl *gomock.Controller) *MockLogViewerService_GetLogsClient {
	mock := &MockLogViewerService_GetLogsClient{ctrl: ctrl}
	mock.recorder = &MockLogViewerService_GetLogsClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLogViewerService_GetLogsClient) EXPECT() *MockLogViewerService_GetLogsClientMockRecorder {
	return m.recorder
}

// Recv mocks base method
func (m *MockLogViewerService_GetLogsClient) Recv() (*proto.LogLine, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*proto.LogLine)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv
func (mr *MockLogViewerService_GetLogsClientMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockLogViewerService_GetLogsClient)(nil).Recv))
}

// Header mocks base method
func (m *MockLogViewerService_GetLogsClient) Header() (metadata.MD, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Header")
	ret0, _ := ret[0].(metadata.MD)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Header indicates an expected call of Header
func (mr *MockLogViewerService_GetLogsClientMockRecorder) Header() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockLogViewerService_GetLogsClient)(nil).Header))
}

// Trailer mocks base method
func (m *MockLogViewerService_GetLogsClient) Trailer() metadata.MD {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trailer")
	ret0, _ := ret[0].(metadata.MD)
	return ret0
}

// Trailer indicates an expected call of Trailer
func (mr *MockLogViewerService_GetLogsClientMockRecorder) Trailer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trailer", reflect.TypeOf((*MockLogViewerService_GetLogsClient)(nil).Trailer))
}

// CloseSend mocks base method
func (m *MockLogViewerService_GetLogsClient) CloseSend() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseSend")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseSend indicates an expected call of CloseSend
func (mr *MockLogViewerService_GetLogsClientMockRecorder) CloseSend() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseSend", reflect.TypeOf((*MockLogViewerService_GetLogsClient)(nil).CloseSend))
}

// Context mocks base method
func (m *MockLogViewerService_GetLogsClient) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context
func (mr *MockLogViewerService_GetLogsClientMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockLogViewerService_GetLogsClient)(nil).Context))
}

// SendMsg mocks base method
func (m_2 *MockLogViewerService_GetLogsClient) SendMsg(m interface{}) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "SendMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg
func (mr *MockLogViewerService_GetLogsClientMockRecorder) SendMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockLogViewerService_GetLogsClient)(nil).SendMsg), m)
}

// RecvMsg mocks base method
func (m_2 *MockLogViewerService_GetLogsClient) RecvMsg(m interface{}) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "RecvMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg
func (mr *MockLogViewerService_GetLogsClientMockRecorder) RecvMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockLogViewerService_GetLogsClient)(nil).RecvMsg), m)
}

// MockLogViewerServiceServer is a mock of LogViewerServiceServer interface
type MockLogViewerServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockLogViewerServiceServerMockRecorder
}

// MockLogViewerServiceServerMockRecorder is the mock recorder for MockLogViewerServiceServer
type MockLogViewerServiceServerMockRecorder struct {
	mock *MockLogViewerServiceServer
}

// NewMockLogViewerServiceServer creates a new mock instance
func NewMockLogViewerServiceServer(ctrl *gomock.Controller) *MockLogViewerServiceServer {
	mock := &MockLogViewerServiceServer{ctrl: ctrl}
	mock.recorder = &MockLogViewerServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLogViewerServiceServer) EXPECT() *MockLogViewerServiceServerMockRecorder {
	return m.recorder
}

// GetLogs mocks base method
func (m *MockLogViewerServiceServer) GetLogs(arg0 *proto.LogFileLocation, arg1 proto.LogViewerService_GetLogsServer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLogs", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetLogs indicates an expected call of GetLogs
func (mr *MockLogViewerServiceServerMockRecorder) GetLogs(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLogs", reflect.TypeOf((*MockLogViewerServiceServer)(nil).GetLogs), arg0, arg1)
}

// MockLogViewerService_GetLogsServer is a mock of LogViewerService_GetLogsServer interface
type MockLogViewerService_GetLogsServer struct {
	ctrl     *gomock.Controller
	recorder *MockLogViewerService_GetLogsServerMockRecorder
}

// MockLogViewerService_GetLogsServerMockRecorder is the mock recorder for MockLogViewerService_GetLogsServer
type MockLogViewerService_GetLogsServerMockRecorder struct {
	mock *MockLogViewerService_GetLogsServer
}

// NewMockLogViewerService_GetLogsServer creates a new mock instance
func NewMockLogViewerService_GetLogsServer(ctrl *gomock.Controller) *MockLogViewerService_GetLogsServer {
	mock := &MockLogViewerService_GetLogsServer{ctrl: ctrl}
	mock.recorder = &MockLogViewerService_GetLogsServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLogViewerService_GetLogsServer) EXPECT() *MockLogViewerService_GetLogsServerMockRecorder {
	return m.recorder
}

// Send mocks base method
func (m *MockLogViewerService_GetLogsServer) Send(arg0 *proto.LogLine) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send
func (mr *MockLogViewerService_GetLogsServerMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockLogViewerService_GetLogsServer)(nil).Send), arg0)
}

// SetHeader mocks base method
func (m *MockLogViewerService_GetLogsServer) SetHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetHeader indicates an expected call of SetHeader
func (mr *MockLogViewerService_GetLogsServerMockRecorder) SetHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeader", reflect.TypeOf((*MockLogViewerService_GetLogsServer)(nil).SetHeader), arg0)
}

// SendHeader mocks base method
func (m *MockLogViewerService_GetLogsServer) SendHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendHeader indicates an expected call of SendHeader
func (mr *MockLogViewerService_GetLogsServerMockRecorder) SendHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendHeader", reflect.TypeOf((*MockLogViewerService_GetLogsServer)(nil).SendHeader), arg0)
}

// SetTrailer mocks base method
func (m *MockLogViewerService_GetLogsServer) SetTrailer(arg0 metadata.MD) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTrailer", arg0)
}

// SetTrailer indicates an expected call of SetTrailer
func (mr *MockLogViewerService_GetLogsServerMockRecorder) SetTrailer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTrailer", reflect.TypeOf((*MockLogViewerService_GetLogsServer)(nil).SetTrailer), arg0)
}

// Context mocks base method
func (m *MockLogViewerService_GetLogsServer) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context
func (mr *MockLogViewerService_GetLogsServerMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockLogViewerService_GetLogsServer)(nil).Context))
}

// SendMsg mocks base method
func (m_2 *MockLogViewerService_GetLogsServer) SendMsg(m interface{}) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "SendMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg
func (mr *MockLogViewerService_GetLogsServerMockRecorder) SendMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockLogViewerService_GetLogsServer)(nil).SendMsg), m)
}

// RecvMsg mocks base method
func (m_2 *MockLogViewerService_GetLogsServer) RecvMsg(m interface{}) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "RecvMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg
func (mr *MockLogViewerService_GetLogsServerMockRecorder) RecvMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockLogViewerService_GetLogsServer)(nil).RecvMsg), m)
}
