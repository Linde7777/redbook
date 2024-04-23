package service

import (
	"context"
	"fmt"
	"main/internal/repository"
	"main/internal/service/sms"
	"math/rand"
)

type AuthCodeService struct {
	repo repository.AuthCodeRepository
	sms  sms.Service
}

func NewAuthCodeService(repo repository.AuthCodeRepository, sms sms.Service) *AuthCodeService {
	return &AuthCodeService{
		repo: repo,
		sms:  sms,
	}
}

func (svc *AuthCodeService) SendAuthCode(ctx context.Context, businessName, phoneNumber string) (httpCode int, err error) {
	authCode := svc.generateAuthCode()
	httpCode, err = svc.repo.Set(ctx, businessName, phoneNumber, authCode)
	if err != nil {
		return httpCode, err
	}

	const templateID = "todo"
	return svc.sms.Send(ctx, templateID, []string{phoneNumber}, authCode)
}

func (svc *AuthCodeService) VerifyAuthCode(ctx context.Context, businessName, phoneNumber, authCode string) (httpCode int, err error) {
	return svc.repo.Verify(ctx, businessName, phoneNumber, authCode)
}

func (svc *AuthCodeService) generateAuthCode() string {
	rawAuthCode := rand.Intn(1000000)
	return fmt.Sprintf("%06d", rawAuthCode)
}
