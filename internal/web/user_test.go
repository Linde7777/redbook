package web

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"main/internal/domain"
	"main/internal/service"
	servicemock "main/internal/service/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

func generateEncryptedPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

func TestUserHandler_Signup(t *testing.T) {
	email1 := "abc@gmail.com"
	password1 := "123456789"
	encryptedPassword1 := generateEncryptedPassword(password1)
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
				userService := servicemock.NewMockUserService(ctrl)
				userService.EXPECT().Signup(gomock.Any(), &domain.User{
					Email:    email1,
					Password: password1,
				}).Return(domain.User{
					ID:       1,
					Email:    email1,
					Password: encryptedPassword1,
				}, http.StatusOK, nil)
				authCodeService := servicemock.NewMockAuthCodeService(ctrl)
				return userService, authCodeService
			},
			reqBuilder: func(t *testing.T) *http.Request {
				reqBodyStruct := ReqSignup{
					Email:           email1,
					Password:        password1,
					ConfirmPassword: password1,
				}
				reqBodyBytes, err := json.Marshal(reqBodyStruct)
				assert.NoError(t, err)
				reqBody := bytes.NewReader(reqBodyBytes)
				req := httptest.NewRequest(http.MethodPost, "/v1/user/signup", reqBody)
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			expectedHTTPCode: http.StatusOK,
			expectedResponse: "signup success",
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
