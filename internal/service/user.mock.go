// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/service/user.go
//
// Generated by this command:
//
//	mockgen -source=./internal/service/user.go -package=service -destination=./internal/service/user.mock.go
//

// Package service is a generated GoMock package.
package service

import (
	context "context"
	domain "main/internal/domain"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockUserService is a mock of UserService interface.
type MockUserService struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceMockRecorder
}

// MockUserServiceMockRecorder is the mock recorder for MockUserService.
type MockUserServiceMockRecorder struct {
	mock *MockUserService
}

// NewMockUserService creates a new mock instance.
func NewMockUserService(ctrl *gomock.Controller) *MockUserService {
	mock := &MockUserService{ctrl: ctrl}
	mock.recorder = &MockUserServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserService) EXPECT() *MockUserServiceMockRecorder {
	return m.recorder
}

// LoginByPassword mocks base method.
func (m *MockUserService) LoginByPassword(ctx context.Context, user domain.User) (domain.User, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoginByPassword", ctx, user)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// LoginByPassword indicates an expected call of LoginByPassword.
func (mr *MockUserServiceMockRecorder) LoginByPassword(ctx, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoginByPassword", reflect.TypeOf((*MockUserService)(nil).LoginByPassword), ctx, user)
}

// SearchOrCreateUserByPhoneNumber mocks base method.
func (m *MockUserService) SearchOrCreateUserByPhoneNumber(ctx context.Context, phoneNumber string) (domain.User, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchOrCreateUserByPhoneNumber", ctx, phoneNumber)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SearchOrCreateUserByPhoneNumber indicates an expected call of SearchOrCreateUserByPhoneNumber.
func (mr *MockUserServiceMockRecorder) SearchOrCreateUserByPhoneNumber(ctx, phoneNumber any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchOrCreateUserByPhoneNumber", reflect.TypeOf((*MockUserService)(nil).SearchOrCreateUserByPhoneNumber), ctx, phoneNumber)
}

// Signup mocks base method.
func (m *MockUserService) Signup(ctx context.Context, user domain.User) (domain.User, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Signup", ctx, user)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Signup indicates an expected call of Signup.
func (mr *MockUserServiceMockRecorder) Signup(ctx, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Signup", reflect.TypeOf((*MockUserService)(nil).Signup), ctx, user)
}
