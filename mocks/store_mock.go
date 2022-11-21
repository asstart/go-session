// Code generated by MockGen. DO NOT EDIT.
// Source: session/store.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	session "github.com/asstart/go-session"
	gomock "github.com/golang/mock/gomock"
)

// MockSessionStore is a mock of SessionStore interface.
type MockSessionStore struct {
	ctrl     *gomock.Controller
	recorder *MockSessionStoreMockRecorder
}

// MockSessionStoreMockRecorder is the mock recorder for MockSessionStore.
type MockSessionStoreMockRecorder struct {
	mock *MockSessionStore
}

// NewMockSessionStore creates a new mock instance.
func NewMockSessionStore(ctrl *gomock.Controller) *MockSessionStore {
	mock := &MockSessionStore{ctrl: ctrl}
	mock.recorder = &MockSessionStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionStore) EXPECT() *MockSessionStoreMockRecorder {
	return m.recorder
}

// Invalidate mocks base method.
func (m *MockSessionStore) Invalidate(ctx context.Context, sid string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Invalidate", ctx, sid)
	ret0, _ := ret[0].(error)
	return ret0
}

// Invalidate indicates an expected call of Invalidate.
func (mr *MockSessionStoreMockRecorder) Invalidate(ctx, sid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Invalidate", reflect.TypeOf((*MockSessionStore)(nil).Invalidate), ctx, sid)
}

// Load mocks base method.
func (m *MockSessionStore) Load(ctx context.Context, sid string) (session.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Load", ctx, sid)
	ret0, _ := ret[0].(session.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Load indicates an expected call of Load.
func (mr *MockSessionStoreMockRecorder) Load(ctx, sid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Load", reflect.TypeOf((*MockSessionStore)(nil).Load), ctx, sid)
}

// Save mocks base method.
func (m *MockSessionStore) Save(ctx context.Context, s session.Session) (session.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, s)
	ret0, _ := ret[0].(session.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Save indicates an expected call of Save.
func (mr *MockSessionStoreMockRecorder) Save(ctx, s interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockSessionStore)(nil).Save), ctx, s)
}

// Update mocks base method.
func (m *MockSessionStore) Update(ctx context.Context, s session.Session) (session.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, s)
	ret0, _ := ret[0].(session.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockSessionStoreMockRecorder) Update(ctx, s interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockSessionStore)(nil).Update), ctx, s)
}
