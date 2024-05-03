package middlewares

import (
	_ "embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
)

type LimiterBuilder interface {
	// Limit 返回true即触发限流
	limit(ctx *gin.Context, key string) (bool, error)
	Build() gin.HandlerFunc
	SetLimitKeyFunc(func(ctx *gin.Context) string)
}

type RedisLimiterBuilder struct {
	cmd        redis.Cmdable
	keyFunc    func(ctx *gin.Context) string
	threshold  int
	windowSize int
}

var _ LimiterBuilder = &RedisLimiterBuilder{}

// 创建一个基于Redis的集群限流器，默认限制IP，如需要限制其他key，需使用
func NewRedisLimiterBuilder(cmd redis.Cmdable,
	threshold, windowSize int) LimiterBuilder {

	return &RedisLimiterBuilder{
		cmd: cmd,
		keyFunc: func(ctx *gin.Context) string {
			return fmt.Sprintf("ip-limiter:%s", ctx.ClientIP())
		},
		threshold:  threshold,
		windowSize: windowSize,
	}
}

func (l *RedisLimiterBuilder) SetLimitKeyFunc(f func(ctx *gin.Context) string) {
	l.keyFunc = f
}

//go:embed ratelimiter.lua
var luascript string

func (l *RedisLimiterBuilder) limit(ctx *gin.Context, key string) (bool, error) {
	res, err := l.cmd.Eval(ctx, luascript, []string{key}, l.threshold, l.windowSize).Bool()
	if err != nil {
		// todo: 启动单机限流
		return false, err
	}
	return res, nil
}

func (l *RedisLimiterBuilder) Build() gin.HandlerFunc {
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
