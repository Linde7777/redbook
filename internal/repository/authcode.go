package repository

import (
	"context"
	"fmt"
	"main/internal/repository/cache"
	"sync"
)

type AuthCodeRepository interface {
	Set(ctx context.Context, businessName, phoneNumber, authCode string) (httpCode int, err error)
	Verify(ctx context.Context, businessName, phoneNumber, authCode string) (httpCode int, err error)
}

type AuthCodeRepositoryWithCache struct {
	cache cache.AuthCodeCache
}

var (
	authCodeRepoWithCacheOnce sync.Once
	authCodeRepoWithCache     *AuthCodeRepositoryWithCache
	_                         AuthCodeRepository = (*AuthCodeRepositoryWithCache)(nil)
)

// NewAuthCodeRepositoryWithCache 为了适配wire，只能返回接口，而不是返回具体实现
func NewAuthCodeRepositoryWithCache(cache cache.AuthCodeCache) AuthCodeRepository {
	authCodeRepoWithCacheOnce.Do(func() {
		authCodeRepoWithCache = &AuthCodeRepositoryWithCache{
			cache: cache,
		}
	})
	return authCodeRepoWithCache
}

func (c *AuthCodeRepositoryWithCache) Set(ctx context.Context, businessName, phoneNumber, authCode string) (httpCode int, err error) {
	httpCode, err = c.cache.Set(ctx, businessName, phoneNumber, authCode)
	if c.cache.HasExceedSendLimitError() {
		fmt.Println("有攻击者")
	}
	return httpCode, err
}

func (c *AuthCodeRepositoryWithCache) Verify(ctx context.Context, businessName, phoneNumber, authCode string) (httpCode int, err error) {
	return c.cache.Verify(ctx, businessName, phoneNumber, authCode)
}
