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

func (repo *UserRepository) Create(ctx context.Context, user *domain.User) (httpCode int, err error) {
	return repo.dao.Insert(ctx, &dao.User{
		Email:    user.Email,
		Password: user.Password,
	})
}

func (repo *UserRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	daoUser, err := repo.dao.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return toDomainUser(daoUser), nil
}

func toDomainUser(daoUser dao.User) domain.User {
	return domain.User{
		ID:       daoUser.ID,
		Email:    daoUser.Email,
		Password: daoUser.Password,
	}
}
