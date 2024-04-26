package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"testing"
)

func TestRedisAuthCodeCache_Set(t *testing.T) {

	key := func(businessName, phoneNumber string) string {
		return fmt.Sprintf("authcode:%s:%s", businessName, phoneNumber)
	}
	testCases := []struct {
		name             string
		mock             func(ctrl *gomock.Controller) redis.Cmdable
		ctx              context.Context
		businessName     string
		phoneNumber      string
		authCode         string
		expectedHTTPCode int
		expectedErr      error
	}{
		{
			name: "成功",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				mockCmdable := NewMockCmdable(ctrl)
				mockCmd := redis.NewCmd(context.Background())
				mockCmd.SetErr(nil)

				// 一开始是用 mockCmd.SetVal(0)，报错redis: unexpected type=int for Int，
				// 原本的业务代码是这么调用的：
				// res, err := c.cmd.Eval(ctx, luaSetCode, []string{c.key(businessName, phoneNumber)}, authCode).Int()
				// err就是从Int()方法里面来的，点进去看，发现需要返回int64
				mockCmd.SetVal(int64(0))

				mockCmdable.EXPECT().
					Eval(gomock.Any(), luaSetCode,
						[]string{key("test", "123456789")},
						"123456").
					Return(mockCmd)
				return mockCmdable
			},
			ctx:              context.Background(),
			businessName:     "test",
			phoneNumber:      "123456789",
			authCode:         "123456",
			expectedHTTPCode: http.StatusOK,
			expectedErr:      nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockCmd := tc.mock(ctrl)
			authCodeCache := NewRedisAuthCodeCache(mockCmd)
			httpCode, err := authCodeCache.Set(tc.ctx, tc.businessName, tc.phoneNumber, tc.authCode)
			assert.Equal(t, tc.expectedHTTPCode, httpCode)
			assert.Equal(t, tc.expectedErr, err)

		})
	}
}
