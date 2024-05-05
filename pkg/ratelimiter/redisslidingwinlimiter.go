package ratelimiter

import (
	"context"
	_ "embed"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisSlidingWinLimiter struct {
	cmd                redis.Cmdable
	threshold          int
	windowSizeInMinute int
}

var _ RateLimiter = &RedisSlidingWinLimiter{}

func NewRedisSlidingWinLimiter(cmd redis.Cmdable,
	threshold, windowSize int) *RedisSlidingWinLimiter {

	return &RedisSlidingWinLimiter{
		cmd:                cmd,
		threshold:          threshold,
		windowSizeInMinute: windowSize,
	}
}

//go:embed redisslidingwinlimiter.lua
var luaScript string

func (l *RedisSlidingWinLimiter) Limit(ctx context.Context, limitKey string) (bool, error) {
	res, err := l.cmd.Eval(ctx, luaScript, []string{limitKey},
		time.Now().UnixMilli(), l.threshold, l.windowSizeInMinute).Bool()
	if err != nil {
		return false, err
	}
	return res, nil
}
