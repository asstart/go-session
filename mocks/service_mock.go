// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	session "github.com/asstart/go-session"
	gomock "github.com/golang/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// AddAttributes mocks base method.
func (m *MockService) AddAttributes(ctx context.Context, sid string, keyAndValues ...interface{}) (*session.Session, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, sid}
	for _, a := range keyAndValues {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddAttributes", varargs...)
	ret0, _ := ret[0].(*session.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddAttributes indicates an expected call of AddAttributes.
func (mr *MockServiceMockRecorder) AddAttributes(ctx, sid interface{}, keyAndValues ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, sid}, keyAndValues...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddAttributes", reflect.TypeOf((*MockService)(nil).AddAttributes), varargs...)
}

// CreateAnonymSession mocks base method.
func (m *MockService) CreateAnonymSession(ctx context.Context, cc session.CookieConf, sc session.Conf, keyAndValues ...interface{}) (*session.Session, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, cc, sc}
	for _, a := range keyAndValues {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateAnonymSession", varargs...)
	ret0, _ := ret[0].(*session.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAnonymSession indicates an expected call of CreateAnonymSession.
func (mr *MockServiceMockRecorder) CreateAnonymSession(ctx, cc, sc interface{}, keyAndValues ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, cc, sc}, keyAndValues...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAnonymSession", reflect.TypeOf((*MockService)(nil).CreateAnonymSession), varargs...)
}

// CreateUserSession mocks base method.
func (m *MockService) CreateUserSession(ctx context.Context, uid string, cc session.CookieConf, sc session.Conf, keyAndValues ...interface{}) (*session.Session, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, uid, cc, sc}
	for _, a := range keyAndValues {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateUserSession", varargs...)
	ret0, _ := ret[0].(*session.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUserSession indicates an expected call of CreateUserSession.
func (mr *MockServiceMockRecorder) CreateUserSession(ctx, uid, cc, sc interface{}, keyAndValues ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, uid, cc, sc}, keyAndValues...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserSession", reflect.TypeOf((*MockService)(nil).CreateUserSession), varargs...)
}

// InvalidateSession mocks base method.
func (m *MockService) InvalidateSession(ctx context.Context, sid string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InvalidateSession", ctx, sid)
	ret0, _ := ret[0].(error)
	return ret0
}

// InvalidateSession indicates an expected call of InvalidateSession.
func (mr *MockServiceMockRecorder) InvalidateSession(ctx, sid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InvalidateSession", reflect.TypeOf((*MockService)(nil).InvalidateSession), ctx, sid)
}

// LoadSession mocks base method.
func (m *MockService) LoadSession(ctx context.Context, sid string) (*session.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadSession", ctx, sid)
	ret0, _ := ret[0].(*session.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadSession indicates an expected call of LoadSession.
func (mr *MockServiceMockRecorder) LoadSession(ctx, sid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadSession", reflect.TypeOf((*MockService)(nil).LoadSession), ctx, sid)
}

// RemoveAttributes mocks base method.
func (m *MockService) RemoveAttributes(ctx context.Context, sid string, keys ...string) (*session.Session, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, sid}
	for _, a := range keys {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RemoveAttributes", varargs...)
	ret0, _ := ret[0].(*session.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RemoveAttributes indicates an expected call of RemoveAttributes.
func (mr *MockServiceMockRecorder) RemoveAttributes(ctx, sid interface{}, keys ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, sid}, keys...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveAttributes", reflect.TypeOf((*MockService)(nil).RemoveAttributes), varargs...)
}
