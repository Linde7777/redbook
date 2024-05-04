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

// RedisSlidWinLimiterBuilder 是基于Redis的滑动窗口限流器
type RedisSlidWinLimiterBuilder struct {
	cmd        redis.Cmdable
	keyFunc    func(ctx *gin.Context) string
	threshold  int
	windowSize int
}

var _ LimiterBuilder = &RedisSlidWinLimiterBuilder{}

// NewRedisSlidingWindowsLimiterBuilder
// 会创建一个基于Redis的集群限流器。
// 创建一个ip限流器，参数keyFunc示例:
//
//	keyFunc = func(ctx *gin.Context) string {
//		return fmt.Sprintf("ip-limiter:%s", ctx.ClientIP())
//	}
func NewRedisSlidingWindowsLimiterBuilder(cmd redis.Cmdable,
	threshold, windowSize int, keyFunc func(ctx *gin.Context) string) LimiterBuilder {
	if keyFunc == nil {
		panic("keyFunc is nil")
	}

	return &RedisSlidWinLimiterBuilder{
		cmd:        cmd,
		keyFunc:    keyFunc,
		threshold:  threshold,
		windowSize: windowSize,
	}
}

//go:embed ratelimiter.lua
var luascript string

func (l *RedisSlidWinLimiterBuilder) limit(ctx *gin.Context, key string) (bool, error) {
	res, err := l.cmd.Eval(ctx, luascript, []string{key}, l.threshold, l.windowSize).Bool()
	if err != nil {
		// todo: 启动单机限流
		return false, err
	}
	return res, nil
}

func (l *RedisSlidWinLimiterBuilder) Build() gin.HandlerFunc {
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
