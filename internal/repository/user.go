package repository

import (
	"context"
	"main/internal/domain"
	"main/internal/repository/dao"
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (repo *UserRepository) Create(ctx context.Context, user *domain.User) error {
	return repo.dao.Insert(ctx, &dao.User{
		Email:    user.Email,
		Password: user.Password,
	})
}
