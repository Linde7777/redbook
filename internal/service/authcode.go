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

func NewAuthCodeService(repo repository.AuthCodeRepository) *AuthCodeService {
	return &AuthCodeService{repo: repo}
}

func (svc *AuthCodeService) SendAuthCode(ctx context.Context, businessName, phoneNumber string) error {
	authCode := svc.generateAuthCode()
	err := svc.repo.Set(ctx, businessName, phoneNumber, authCode)
	if err != nil {
		return err
	}

	const templateID = "todo"
	return svc.sms.Send(ctx, templateID, []string{phoneNumber}, authCode)
}

func (svc *AuthCodeService) VerifyAuthCode(ctx context.Context, businessName, phoneNumber, authCode string) error {
	return svc.repo.Verify(ctx, businessName, phoneNumber, authCode)
}

func (svc *AuthCodeService) generateAuthCode() string {
	rawAuthCode := rand.Intn(1000000)
	return fmt.Sprintf("%06d", rawAuthCode)
}
