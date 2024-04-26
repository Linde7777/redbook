package intergrationtests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"main/internal/repository/cache"
	"main/internal/testutils"
	"main/internal/web"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

//todo:
// 1.写一个脚本，启动mysql和redis的docker容器
// 2. 改动本目录下的wire.go，使其连接到docker容器的mysql和redis

func Test(t *testing.T) {
	redisClient := InitRedis()
	server := InitWebServer()
	testCases := []struct {
		name             string
		before           func(t *testing.T)
		after            func(t *testing.T)
		phoneNumber      string
		expectedHTTPCode int
		expectedResult   string
	}{
		{
			name: "成功",
			before: func(t *testing.T) {
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()
				temp := cache.NewRedisAuthCodeCache(nil)
				key := temp.Key("test", "123456789")

				authCode, err := redisClient.Get(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, len(authCode) == 6)

				ttl, err := redisClient.TTL(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, ttl > 14*time.Minute)

				err = redisClient.Del(ctx, key).Err()
				assert.NoError(t, err)
			},
			phoneNumber:      "123456789",
			expectedHTTPCode: http.StatusOK,
			expectedResult:   "send auth code success",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)
			req := testutils.GenHTTPJSONReq("POST",
				"/v1/user/send-login-sms-auth-code",
				web.ReqSendLoginSMSAuthCode{PhoneNumber: tc.phoneNumber})
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)
			assert.Equal(t, tc.expectedHTTPCode, recorder.Code)
			assert.Equal(t, tc.expectedResult, recorder.Body.String())
		})
	}
}
