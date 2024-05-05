package service

import (
	"context"
	"fmt"
	"main/internal/repository"
	"math/rand"
)

type AuthCodeService interface {
	SendAuthCode(ctx context.Context, businessName, phoneNumber string) (httpCode int, err error)
	VerifyAuthCode(ctx context.Context, businessName, phoneNumber, authCode string) (httpCode int, err error)
}

type AuthCodeServiceV1 struct {
	repo repository.AuthCodeRepository
	sms  SMSService
}

// NewAuthCodeServiceV1 为了适配wire，只能返回接口，而不是返回具体实现
func NewAuthCodeServiceV1(repo repository.AuthCodeRepository, sms SMSService) AuthCodeService {
	return &AuthCodeServiceV1{
		repo: repo,
		sms:  sms,
	}
}

func (svc *AuthCodeServiceV1) SendAuthCode(ctx context.Context, businessName, phoneNumber string) (httpCode int, err error) {
	authCode := svc.generateAuthCode()
	httpCode, err = svc.repo.Set(ctx, businessName, phoneNumber, authCode)
	if err != nil {
		return httpCode, err
	}

	const templateID = "todo"
	return svc.sms.Send(ctx, templateID, []string{phoneNumber}, authCode)
}

func (svc *AuthCodeServiceV1) VerifyAuthCode(ctx context.Context, businessName, phoneNumber, authCode string) (httpCode int, err error) {
	return svc.repo.Verify(ctx, businessName, phoneNumber, authCode)
}

func (svc *AuthCodeServiceV1) generateAuthCode() string {
	rawAuthCode := rand.Intn(1000000)
	return fmt.Sprintf("%06d", rawAuthCode)
}
