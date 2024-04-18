package cache

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"main/internal/domain"
	"math/rand"
	"time"
)

type UserCache interface {
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	SetUserByEmail(ctx context.Context, user domain.User) error
}

type RedisCache struct {
	cmd redis.Cmdable

	// 大部分kv的缓存时间可以设置为相同的值
	commonExpireDuration time.Duration
}

func (c *RedisCache) randCommonExpDuration() time.Duration {
	return c.commonExpireDuration + time.Duration(rand.Int63n(int64(c.commonExpireDuration)))
}

func NewRedisCache() *RedisCache {
	return &RedisCache{}
}

func (c *RedisCache) keyUserEmail(email string) string {
	return "user:email:" + email
}

func (c *RedisCache) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	data, err := c.cmd.Get(ctx, c.keyUserEmail(email)).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var user domain.User
	err = json.Unmarshal(data, &user)
	return user, err
}

func (c *RedisCache) SetUserByEmail(ctx context.Context, user domain.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return c.cmd.Set(ctx, c.keyUserEmail(user.Email), data, c.randCommonExpDuration()).Err()
}
