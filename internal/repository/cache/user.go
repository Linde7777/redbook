package cache

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"main/internal/domain"
	"math/rand"
	"net/http"
	"time"
)

type UserCache interface {
	GetUserByEmail(ctx context.Context, email string) (user domain.User, httpCode int, err error)
	SetUserByEmail(ctx context.Context, user domain.User) (httpCode int, err error)
}

type RedisUserCache struct {
	cmd redis.Cmdable

	// 大部分kv的缓存时间可以设置为相同的值
	commonExpireDuration time.Duration
}

// NewRedisUserCache 为了适配wire，只能返回接口，而不是返回具体实现
func NewRedisUserCache(cmd redis.Cmdable) UserCache {
	return &RedisUserCache{
		cmd:                  cmd,
		commonExpireDuration: 15 * time.Minute,
	}
}

func (c *RedisUserCache) randCommonExpDuration() time.Duration {
	return c.commonExpireDuration + time.Duration(rand.Int63n(int64(c.commonExpireDuration)))
}

func (c *RedisUserCache) keyUserEmail(email string) string {
	return "user:email:" + email
}

func (c *RedisUserCache) GetUserByEmail(ctx context.Context, email string) (user domain.User, httpCode int, err error) {
	data, err := c.cmd.Get(ctx, c.keyUserEmail(email)).Bytes()
	if err != nil {
		return domain.User{}, http.StatusNotFound, err
	}
	err = json.Unmarshal(data, &user)
	return user, http.StatusOK, err
}

func (c *RedisUserCache) SetUserByEmail(ctx context.Context, user domain.User) (httpCode int, err error) {
	data, err := json.Marshal(user)
	if err != nil {
		return http.StatusBadRequest, err
	}
	err = c.cmd.Set(ctx, c.keyUserEmail(user.Email), data, c.randCommonExpDuration()).Err()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
