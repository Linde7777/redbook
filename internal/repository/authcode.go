package repository

import (
	"context"
	"fmt"
	"main/internal/repository/cache"
)

type AuthCodeRepository struct {
	cache cache.AuthCodeCache
}

func NewAuthCodeRepository(cache cache.AuthCodeCache) *AuthCodeRepository {
	return &AuthCodeRepository{cache: cache}
}

func (c *AuthCodeRepository) Set(ctx context.Context, businessName, phoneNumber, authCode string) (httpCode int, err error) {
	httpCode, err = c.cache.Set(ctx, businessName, phoneNumber, authCode)
	if c.cache.HasExceedSendLimitError() {
		fmt.Println("有攻击者")
	}
	return httpCode, err
}

func (c *AuthCodeRepository) Verify(ctx context.Context, businessName, phoneNumber, authCode string) (httpCode int, err error) {
	return c.cache.Verify(ctx, businessName, phoneNumber, authCode)
}
