// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/repository/dao/user.go
//
// Generated by this command:
//
//	mockgen -source=./internal/repository/dao/user.go -package=dao -destination=./internal/repository/dao/user.mock.go
//

// Package dao is a generated GoMock package.
package dao

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockUserDAO is a mock of UserDAO interface.
type MockUserDAO struct {
	ctrl     *gomock.Controller
	recorder *MockUserDAOMockRecorder
}

// MockUserDAOMockRecorder is the mock recorder for MockUserDAO.
type MockUserDAOMockRecorder struct {
	mock *MockUserDAO
}

// NewMockUserDAO creates a new mock instance.
func NewMockUserDAO(ctrl *gomock.Controller) *MockUserDAO {
	mock := &MockUserDAO{ctrl: ctrl}
	mock.recorder = &MockUserDAOMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserDAO) EXPECT() *MockUserDAOMockRecorder {
	return m.recorder
}

// Insert mocks base method.
func (m *MockUserDAO) Insert(ctx context.Context, user *User) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, user)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert indicates an expected call of Insert.
func (mr *MockUserDAOMockRecorder) Insert(ctx, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockUserDAO)(nil).Insert), ctx, user)
}

// SearchUserByEmail mocks base method.
func (m *MockUserDAO) SearchUserByEmail(ctx context.Context, email string) (User, bool, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchUserByEmail", ctx, email)
	ret0, _ := ret[0].(User)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(int)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// SearchUserByEmail indicates an expected call of SearchUserByEmail.
func (mr *MockUserDAOMockRecorder) SearchUserByEmail(ctx, email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchUserByEmail", reflect.TypeOf((*MockUserDAO)(nil).SearchUserByEmail), ctx, email)
}

// SearchUserByPhoneNumber mocks base method.
func (m *MockUserDAO) SearchUserByPhoneNumber(ctx context.Context, number string) (User, bool, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchUserByPhoneNumber", ctx, number)
	ret0, _ := ret[0].(User)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(int)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// SearchUserByPhoneNumber indicates an expected call of SearchUserByPhoneNumber.
func (mr *MockUserDAOMockRecorder) SearchUserByPhoneNumber(ctx, number any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchUserByPhoneNumber", reflect.TypeOf((*MockUserDAO)(nil).SearchUserByPhoneNumber), ctx, number)
}
