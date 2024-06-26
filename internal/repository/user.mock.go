// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/repository/user.go
//
// Generated by this command:
//
//	mockgen -source=./internal/repository/user.go -package=repository -destination=./internal/repository/user.mock.go
//

// Package repository is a generated GoMock package.
package repository

import (
	context "context"
	domain "main/internal/domain"
	dao "main/internal/repository/dao"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockUserRepository) Create(ctx context.Context, inputDomainUser domain.User) (domain.User, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, inputDomainUser)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Create indicates an expected call of Create.
func (mr *MockUserRepositoryMockRecorder) Create(ctx, inputDomainUser any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUserRepository)(nil).Create), ctx, inputDomainUser)
}

// SearchUserByEmail mocks base method.
func (m *MockUserRepository) SearchUserByEmail(ctx context.Context, email string) (domain.User, bool, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchUserByEmail", ctx, email)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(int)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// SearchUserByEmail indicates an expected call of SearchUserByEmail.
func (mr *MockUserRepositoryMockRecorder) SearchUserByEmail(ctx, email any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchUserByEmail", reflect.TypeOf((*MockUserRepository)(nil).SearchUserByEmail), ctx, email)
}

// SearchUserByPhoneNumber mocks base method.
func (m *MockUserRepository) SearchUserByPhoneNumber(ctx context.Context, phoneNumber string) (domain.User, bool, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchUserByPhoneNumber", ctx, phoneNumber)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(int)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// SearchUserByPhoneNumber indicates an expected call of SearchUserByPhoneNumber.
func (mr *MockUserRepositoryMockRecorder) SearchUserByPhoneNumber(ctx, phoneNumber any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchUserByPhoneNumber", reflect.TypeOf((*MockUserRepository)(nil).SearchUserByPhoneNumber), ctx, phoneNumber)
}

// toDaoUser mocks base method.
func (m *MockUserRepository) toDaoUser(domainUser domain.User) dao.User {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "toDaoUser", domainUser)
	ret0, _ := ret[0].(dao.User)
	return ret0
}

// toDaoUser indicates an expected call of toDaoUser.
func (mr *MockUserRepositoryMockRecorder) toDaoUser(domainUser any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "toDaoUser", reflect.TypeOf((*MockUserRepository)(nil).toDaoUser), domainUser)
}

// toDomainUser mocks base method.
func (m *MockUserRepository) toDomainUser(daoUser dao.User) domain.User {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "toDomainUser", daoUser)
	ret0, _ := ret[0].(domain.User)
	return ret0
}

// toDomainUser indicates an expected call of toDomainUser.
func (mr *MockUserRepositoryMockRecorder) toDomainUser(daoUser any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "toDomainUser", reflect.TypeOf((*MockUserRepository)(nil).toDomainUser), daoUser)
}
