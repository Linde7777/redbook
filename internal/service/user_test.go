package service

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"main/internal/domain"
	"main/internal/repository"
	"main/internal/testutils"
	"net/http"
	"testing"
)

func Test(t *testing.T) {
	pass := testutils.GenEncryptedPassword("123456789")
	fmt.Print(pass)
}

func TestUserServiceV1_LoginByPassword(t *testing.T) {
	testCases := []struct {
		name             string
		mock             func(ctrl *gomock.Controller) repository.UserRepository
		ctx              context.Context
		email            string
		password         string
		expectedUser     domain.User
		expectedHTTPCode int
		expectedErr      error
	}{
		{
			name: "正常情况",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repository.NewMockUserRepository(ctrl)
				repo.EXPECT().
					SearchUserByEmail(gomock.Any(), "abc@gmail.com").
					Return(domain.User{
						Email: "abc@gmail.com",

						// bcrypt对于同一个字符串，每次生成的hash都不一样，所以只能提前运行一次，把结果写死
						Password: "$2a$10$T/ec6oBgBg4nBZeHUUjA9u5I/5/16U9zBDIUG15aEXKmv61HJtAPy",

						PhoneNumber: "13332999999",
					}, true, http.StatusOK, nil)
				return repo
			},
			ctx:      nil,
			email:    "abc@gmail.com",
			password: "123456789",
			expectedUser: domain.User{
				Email: "abc@gmail.com",

				// bcrypt对于同一个字符串，每次生成的hash都不一样，所以只能提前运行一次，把结果写死
				Password: "$2a$10$T/ec6oBgBg4nBZeHUUjA9u5I/5/16U9zBDIUG15aEXKmv61HJtAPy",

				PhoneNumber: "13332999999",
			},
			expectedHTTPCode: http.StatusOK,
			expectedErr:      nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := tc.mock(ctrl)
			svc := NewUserServiceV1(repo)
			user, httpCode, err := svc.LoginByPassword(tc.ctx, domain.User{
				Email:    tc.email,
				Password: tc.password,
			})
			assert.Equal(t, tc.expectedHTTPCode, httpCode)
			assert.Equal(t, tc.expectedUser, user)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
