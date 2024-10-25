package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	resty "github.com/go-resty/resty/v2"
)

// MockClient is a mock of Client interface
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// R mocks base method
func (m *MockClient) R() *resty.Request {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "R")
	ret0, _ := ret[0].(*resty.Request)
	return ret0
}

// R indicates an expected call of R
func (mr *MockClientMockRecorder) R() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "R", reflect.TypeOf((*MockClient)(nil).R))
}
