package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"net/http"
)

//go:embed lua/set_auth_code.lua
var luaSetCode string

//go:embed lua/verify_auth_code.lua
var luaVerifyCode string

type AuthCodeCache interface {
	Set(ctx context.Context, businessName, phoneNumber, authCode string) (httpCode int, err error)
	Verify(ctx context.Context, businessName, phoneNumber, authCode string) (httpCode int, err error)
	HasExceedSendLimitError() bool
	Key(businessName, phoneNumber string) string
}

type RedisAuthCodeCache struct {
	cmd redis.Cmdable
}

// NewRedisAuthCodeCache 为了适配wire，只能返回接口，而不是返回具体实现
func NewRedisAuthCodeCache(cmd redis.Cmdable) AuthCodeCache {
	return &RedisAuthCodeCache{cmd: cmd}
}

func (c *RedisAuthCodeCache) Key(businessName, phoneNumber string) string {
	return fmt.Sprintf("authcode:%s:%s", businessName, phoneNumber)
}

var errSentinel error
var errExceedSendLimit = errors.New("超过一分钟内发送次数限制，有攻击者")

func (c *RedisAuthCodeCache) HasExceedSendLimitError() bool {
	return errors.Is(errSentinel, errExceedSendLimit)
}

// Set 在验证码发送次数超过限制时，为防止调用者粗心把这个错误返回给前端，
// 不会返回错误，调用者判断是否有这种错误需要调用 HasExceedSendLimitError。
// 正常用户受到前端限制，不可能在一分钟内请求发送多次验证码，
// 我们需要对攻击者隐藏这个错误，增加攻击者成本
func (c *RedisAuthCodeCache) Set(ctx context.Context, businessName, phoneNumber, authCode string) (httpCode int, err error) {
	res := c.cmd.Eval(ctx, luaSetCode, []string{c.Key(businessName, phoneNumber)}, authCode).String()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	switch res {
	case "err no expire":
		return http.StatusInternalServerError, errors.New("验证码key存在，但无过期时间")
	case "err exceed send limit":
		// todo: zap error
		fmt.Println(errExceedSendLimit)
		errSentinel = errExceedSendLimit
		return http.StatusOK, nil
	case "ok":
		return http.StatusOK, nil
	default:
		// 不要把错误暴露给调用者，调用者可能会不小心返回给前端
		fmt.Println("未知错误，检查lua脚本！")
		return http.StatusInternalServerError, errors.New("未知错误，请查看日志")
	}
}

func (c *RedisAuthCodeCache) Verify(ctx context.Context, businessName, phoneNumber, authCode string) (httpCode int, err error) {
	res := c.cmd.Eval(ctx, luaVerifyCode, []string{c.Key(businessName, phoneNumber)}, authCode).String()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	switch res {
	case "err exceed verify limit":
		return http.StatusBadRequest, errors.New("验证次数耗尽，请重新获取验证码")
	case "err auth code not match":
		return http.StatusBadRequest, errors.New("验证码不匹配，请重新输入")
	case "ok":
		return http.StatusOK, nil
	default:
		// 不要把错误暴露给调用者，调用者可能会不小心返回给前端
		fmt.Println("未知错误，检查lua脚本！")
		return http.StatusInternalServerError, errors.New("未知错误，请查看日志")
	}
}
