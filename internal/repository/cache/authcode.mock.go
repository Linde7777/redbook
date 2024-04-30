// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/repository/cache/authcode.go
//
// Generated by this command:
//
//	mockgen -source=./internal/repository/cache/authcode.go -package=cache -destination=./internal/repository/cache/authcode.mock.go
//

// Package cache is a generated GoMock package.
package cache

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockAuthCodeCache is a mock of AuthCodeCache interface.
type MockAuthCodeCache struct {
	ctrl     *gomock.Controller
	recorder *MockAuthCodeCacheMockRecorder
}

// MockAuthCodeCacheMockRecorder is the mock recorder for MockAuthCodeCache.
type MockAuthCodeCacheMockRecorder struct {
	mock *MockAuthCodeCache
}

// NewMockAuthCodeCache creates a new mock instance.
func NewMockAuthCodeCache(ctrl *gomock.Controller) *MockAuthCodeCache {
	mock := &MockAuthCodeCache{ctrl: ctrl}
	mock.recorder = &MockAuthCodeCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthCodeCache) EXPECT() *MockAuthCodeCacheMockRecorder {
	return m.recorder
}

// HasExceedSendLimitError mocks base method.
func (m *MockAuthCodeCache) HasExceedSendLimitError() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasExceedSendLimitError")
	ret0, _ := ret[0].(bool)
	return ret0
}

// HasExceedSendLimitError indicates an expected call of HasExceedSendLimitError.
func (mr *MockAuthCodeCacheMockRecorder) HasExceedSendLimitError() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasExceedSendLimitError", reflect.TypeOf((*MockAuthCodeCache)(nil).HasExceedSendLimitError))
}

// Key mocks base method.
func (m *MockAuthCodeCache) Key(businessName, phoneNumber string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Key", businessName, phoneNumber)
	ret0, _ := ret[0].(string)
	return ret0
}

// Key indicates an expected call of Key.
func (mr *MockAuthCodeCacheMockRecorder) Key(businessName, phoneNumber any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Key", reflect.TypeOf((*MockAuthCodeCache)(nil).Key), businessName, phoneNumber)
}

// Set mocks base method.
func (m *MockAuthCodeCache) Set(ctx context.Context, businessName, phoneNumber, authCode string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, businessName, phoneNumber, authCode)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Set indicates an expected call of Set.
func (mr *MockAuthCodeCacheMockRecorder) Set(ctx, businessName, phoneNumber, authCode any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockAuthCodeCache)(nil).Set), ctx, businessName, phoneNumber, authCode)
}

// Verify mocks base method.
func (m *MockAuthCodeCache) Verify(ctx context.Context, businessName, phoneNumber, authCode string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Verify", ctx, businessName, phoneNumber, authCode)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Verify indicates an expected call of Verify.
func (mr *MockAuthCodeCacheMockRecorder) Verify(ctx, businessName, phoneNumber, authCode any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Verify", reflect.TypeOf((*MockAuthCodeCache)(nil).Verify), ctx, businessName, phoneNumber, authCode)
}