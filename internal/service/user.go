package service

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"main/internal/domain"
	"main/internal/repository"
	"net/http"
)

type UserService interface {
	Signup(ctx context.Context, user domain.User) (output domain.User, httpCode int, err error)
	LoginByPassword(ctx context.Context, user domain.User) (dbUser domain.User, httpCode int, err error)
	SearchOrCreateUserByPhoneNumber(ctx context.Context, phoneNumber string) (user domain.User, httpCode int, err error)
}

type UserServiceV1 struct {
	repo repository.UserRepository
}

// NewUserServiceV1 为了适配wire，只能返回接口，而不是返回具体实现
func NewUserServiceV1(repo repository.UserRepository) UserService {
	return &UserServiceV1{repo: repo}
}

func (svc *UserServiceV1) Signup(ctx context.Context, user domain.User) (output domain.User, httpCode int, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, http.StatusServiceUnavailable, err
	}
	user.Password = string(hash)
	return svc.repo.Create(ctx, user)
}

func (svc *UserServiceV1) LoginByPassword(ctx context.Context, inputUser domain.User) (
	dbUser domain.User, httpCode int, err error) {

	// 不暴露是否用户不存在，提高攻击者成本
	dbUser, httpCode, err = svc.repo.SearchUserByEmail(ctx, inputUser.Email)
	if err != nil {
		return domain.User{}, httpCode, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(inputUser.Password), []byte(dbUser.Password))
	if err != nil {
		return domain.User{}, http.StatusBadRequest, err
	}
	return dbUser, http.StatusOK, nil
}

func (svc *UserServiceV1) SearchOrCreateUserByPhoneNumber(ctx context.Context, phoneNumber string) (user domain.User, httpCode int, err error) {
	user, ok, httpCode, err := svc.repo.SearchUserByPhoneNumber(ctx, phoneNumber)
	switch {
	case err != nil:
		return domain.User{}, httpCode, err
	case ok:
		return user, http.StatusOK, nil
	}

	user, httpCode, err = svc.repo.Create(ctx, domain.User{
		PhoneNumber: phoneNumber,
	})
	if err != nil {
		return domain.User{}, httpCode, err
	}
	return user, httpCode, nil
}
