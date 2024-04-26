package web

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"main/internal/domain"
	"main/internal/service"
	"main/internal/testutils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_Signup(t *testing.T) {
	testCases := []struct {
		name             string
		mock             func(ctrl *gomock.Controller) (service.UserService, service.AuthCodeService)
		reqBuilder       func(t *testing.T) *http.Request
		expectedHTTPCode int
		expectedResponse string
	}{
		{
			name: "正常情况",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.AuthCodeService) {
				userService := service.NewMockUserService(ctrl)
				userService.EXPECT().Signup(gomock.Any(), domain.User{
					Email:    "abc@gmail.com",
					Password: "123456789",
				}).Return(domain.User{
					ID:       1,
					Email:    "abc@gmail.com",
					Password: testutils.GenEncryptedPassword("123456789"),
				}, http.StatusOK, nil)
				authCodeService := service.NewMockAuthCodeService(ctrl)
				return userService, authCodeService
			},
			reqBuilder: func(t *testing.T) *http.Request {
				reqBodyStruct := ReqSignup{
					Email:           "abc@gmail.com",
					Password:        "123456789",
					ConfirmPassword: "123456789",
				}
				return testutils.GenHTTPJSONReq("POST", "/v1/user/signup", reqBodyStruct)
			},
			expectedHTTPCode: http.StatusOK,
			expectedResponse: "signup success",
		},
		{
			name: "邮箱格式错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.AuthCodeService) {
				return nil, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				reqBodyStruct := ReqSignup{
					Email:           "abc@.com",
					Password:        "123456789",
					ConfirmPassword: "123456789",
				}
				return testutils.GenHTTPJSONReq("POST", "/v1/user/signup", reqBodyStruct)
			},
			expectedHTTPCode: http.StatusBadRequest,
			expectedResponse: "Key: 'ReqSignup.Email' Error:Field validation for 'Email' failed on the 'email' tag",
		},
		{
			name: "密码长度过短",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.AuthCodeService) {
				return nil, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				reqBodyStruct := ReqSignup{
					Email:           "abc@gmail.com",
					Password:        "123",
					ConfirmPassword: "123",
				}
				return testutils.GenHTTPJSONReq("POST", "/v1/user/signup", reqBodyStruct)
			},
			expectedHTTPCode: http.StatusBadRequest,
			expectedResponse: "Key: 'ReqSignup.Password' Error:Field validation for 'Password' failed on the 'min' tag",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userService, authCodeService := testCase.mock(ctrl)
			userHandler := NewUserHandler(userService, authCodeService)

			router := gin.Default()
			userHandler.RegisterRoutes(router)

			req := testCase.reqBuilder(t)
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			assert.Equal(t, testCase.expectedHTTPCode, recorder.Code)
			assert.Equal(t, testCase.expectedResponse, recorder.Body.String())
		})
	}
}
