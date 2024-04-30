package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"main/internal/domain"
	"main/internal/testutils"
	"net/http"
	"testing"
)

func TestRedisUserCache_GetUserByEmail(t *testing.T) {
	testCases := []struct {
		name             string
		mock             func(ctrl *gomock.Controller) redis.Cmdable
		inputEmail       string
		expectedUser     domain.User
		expectedHTTPCode int
	}{
		{
			name: "成功",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				mockCmdable := testutils.NewMockCmdable(ctrl)
				mockCMD := redis.NewStringCmd(context.Background())
				mockCMD.SetErr(nil)
				mockCMD.SetVal(testutils.StructToString(domain.User{
					ID:    1,
					Email: "abc@gmail.com",
				}))
				mockCmdable.EXPECT().Get(gomock.Any(), "user:email:abc@gmail.com").Return(mockCMD)
				return mockCmdable
			},
			inputEmail: "abc@gmail.com",
			expectedUser: domain.User{
				ID:    1,
				Email: "abc@gmail.com",
			},
			expectedHTTPCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			redisClient := tc.mock(ctrl)
			redisUserCache := NewRedisUserCache(redisClient)
			user, httpCode, err := redisUserCache.GetUserByEmail(context.Background(), tc.inputEmail)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedUser, user)
			assert.Equal(t, tc.expectedHTTPCode, httpCode)
		})
	}
}
