package service

import (
	"context"
	"golang.org/x/crypto/bcrypt"
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
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	return svc.repo.Create(ctx, user)
}

func (svc *UserService) Login(ctx context.Context, inputUser *domain.User) (*domain.User, error) {
	dbUser, err := svc.repo.GetUserByEmail(ctx, inputUser.Email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(inputUser.Password), []byte(dbUser.Password))
	if err != nil {
		return nil, err
	}
	return dbUser, nil
}
