package middlewares

import (
	_ "embed"
	"github.com/gin-gonic/gin"
	"main/pkg/ratelimiter"
	"net/http"
)

type RateLimiterMiddlewareBuilder struct {
	limiter    ratelimiter.RateLimiter
	keyGenFunc func(ctx *gin.Context) string
}

// NewRedisSlidingWindowsLimiterBuilder
// 会创建一个基于Redis的滑动窗口集群限流器。
// 创建一个ip限流器，参数keyGenFunc示例:
//
//	keyGenFunc = func(ctx *gin.Context) string {
//		return fmt.Sprintf("ip-ratelimiter:%s", ctx.ClientIP())
//	}
func NewRedisSlidingWindowsLimiterBuilder(limiter ratelimiter.RateLimiter, keyGenFunc func(ctx *gin.Context) string) RateLimiterMiddlewareBuilder {
	return NewRedisSlidingWindowsLimiterBuilder(limiter, keyGenFunc)
}

func (l *RateLimiterMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := l.keyGenFunc(ctx)
		if limit, err := l.limiter.Limit(ctx, key); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else if limit {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit"})
			return
		}
		ctx.Next()
	}
}
