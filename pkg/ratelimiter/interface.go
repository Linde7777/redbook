package ratelimiter

import "context"

type RateLimiter interface {
	// Limit 返回true即触发限流。
	// 需要传入参数key，而不是把它设置为限流器实现内部的一个字段，
	// 让实现调用Limit的时候引用内部字段key，
	// 是因为一个实现可能会有限制多个key的需求，比如做中间件时按IP限流。
	Limit(ctx context.Context, key string) (bool, error)
}
