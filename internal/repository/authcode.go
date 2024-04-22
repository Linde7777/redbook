package repository

import (
	"context"
	"fmt"
	"main/internal/repository/cache"
)

type AuthCodeRepository struct {
	cache cache.AuthCodeCache
}

func (c *AuthCodeRepository) Set(ctx context.Context, businessName, phoneNumber, authCode string) error {
	err := c.cache.Set(ctx, businessName, phoneNumber, authCode)
	if c.cache.HasExceedSendLimitError() {
		fmt.Println("有攻击者")
	}
	return err
}

func (c *AuthCodeRepository) Verify(ctx context.Context, businessName, phoneNumber, authCode string) error {
	return c.cache.Verify(ctx, businessName, phoneNumber, authCode)
}
