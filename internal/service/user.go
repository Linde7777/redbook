package service

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"main/internal/domain"
	"main/internal/repository"
	"net/http"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (svc *UserService) Signup(ctx context.Context, user *domain.User) (httpCode int, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return http.StatusServiceUnavailable, err
	}
	user.Password = string(hash)
	return svc.repo.Create(ctx, user)
}

func (svc *UserService) Login(ctx context.Context, inputUser *domain.User) (
	dbUser domain.User, httpCode int, err error) {

	dbUser, err = svc.repo.GetUserByEmail(ctx, inputUser.Email)
	if err != nil {
		return domain.User{}, http.StatusInternalServerError, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(inputUser.Password), []byte(dbUser.Password))
	if err != nil {
		return domain.User{}, http.StatusInternalServerError, err
	}
	return dbUser, http.StatusOK, nil
}
