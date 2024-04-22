package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

//go:embed lua/set_auth_code.lua
var luaSetCode string

//go:embed lua/verify_auth_code.lua
var luaVerifyCode string

type AuthCodeCache interface {
	Set(ctx context.Context, businessName, phoneNumber, authCode string) error
	Verify(ctx context.Context, businessName, phoneNumber, authCode string) error
	HasExceedSendLimitError() bool
}

type RedisAuthCodeCache struct {
	cmd redis.Cmdable
}

func NewRedisAuthCodeCache(cmd redis.Cmdable) *RedisAuthCodeCache {
	return &RedisAuthCodeCache{cmd: cmd}
}

func (c *RedisAuthCodeCache) key(businessName, phoneNumber string) string {
	return "authcode:" + businessName + ":" + phoneNumber
}

var errSentinel error
var errExceedSendLimit = errors.New("超过一分钟内发送次数限制，有攻击者")

func (c *RedisAuthCodeCache) HasExceedSendLimitError() bool {
	return errors.Is(errSentinel, errExceedSendLimit)
}

func (c *RedisAuthCodeCache) Set(ctx context.Context, businessName, phoneNumber, authCode string) error {
	res, err := c.cmd.Eval(ctx, luaSetCode, []string{c.key(businessName, phoneNumber)}, authCode).Int()
	if err != nil {
		return err
	}
	switch res {
	case 100:
		return errors.New("验证码key存在，但无过期时间")
	case 200:
		// 不对调用者直接返回这个错误，调用者有可能粗心把这个错误暴露给前端（攻击者）。调用在需要调用 HasExceedSendLimitError 判断是否有这个错误
		// todo: zap error
		fmt.Println(errExceedSendLimit)
		errSentinel = errExceedSendLimit
		return nil
	default:
		return nil
	}
}

func (c *RedisAuthCodeCache) Verify(ctx context.Context, businessName, phoneNumber, authCode string) error {
	res, err := c.cmd.Eval(ctx, luaVerifyCode, []string{c.key(businessName, phoneNumber)}, authCode).Int()
	if err != nil {
		return err
	}
	switch res {
	case 300:
		return errors.New("验证次数耗尽，请重新获取验证码")
	case 400:
		return errors.New("验证码不匹配，请重新输入")
	default:
		return nil
	}
}
