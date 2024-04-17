package service

import (
	"context"
	"main/internal/domain"
	"main/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (svc *UserService) Signup(ctx context.Context, user *domain.User) error {
	return svc.repo.Create(ctx, user)
}
