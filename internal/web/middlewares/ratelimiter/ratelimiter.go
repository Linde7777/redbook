package middlewares

import (
	_ "embed"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
)

type LimiterBuilder interface {
	// Limit 返回true即触发限流
	limit(ctx *gin.Context, key string) (bool, error)
	Build() gin.HandlerFunc
}

type RedisLimiter struct {
	cmd        redis.Cmdable
	keyFunc    func(ctx *gin.Context) string
	threshold  int
	windowSize int
}

var _ LimiterBuilder = &RedisLimiter{}

func NewRedisLimiter(cmd redis.Cmdable, keyFunc func(ctx *gin.Context) string,
	threshold, windowSize int) LimiterBuilder {

	return &RedisLimiter{
		cmd:        cmd,
		keyFunc:    keyFunc,
		threshold:  threshold,
		windowSize: windowSize,
	}
}

//go:embed ratelimiter.lua
var luascript string

func (l *RedisLimiter) limit(ctx *gin.Context, key string) (bool, error) {
	res, err := l.cmd.Eval(ctx, luascript, []string{key}, l.threshold, l.windowSize).Bool()
	if err != nil {
		// todo: 启动单机限流
		return false, err
	}
	return res, nil
}

func (l *RedisLimiter) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := l.keyFunc(ctx)
		if limit, err := l.limit(ctx, key); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else if limit {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit"})
			return
		}
		ctx.Next()
	}
}
