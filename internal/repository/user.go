package repository

import (
	"context"
	"main/internal/domain"
	"main/internal/repository/cache"
	"main/internal/repository/dao"
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache cache.UserCache
}

func NewUserRepository(dao *dao.UserDAO, cache cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: cache,
	}
}

func (repo *UserRepository) Create(ctx context.Context, user *domain.User) (httpCode int, err error) {
	return repo.dao.Insert(ctx, &dao.User{
		Email:    user.Email,
		Password: user.Password,
	})
}

func (repo *UserRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	domainUser, err := repo.cache.GetUserByEmail(ctx, email)
	if err == nil {
		return domainUser, nil
	}

	daoUser, err := repo.dao.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	domainUser = toDomainUser(daoUser)
	err = repo.cache.SetUserByEmail(ctx, domainUser)

	return domainUser, nil
}

func toDomainUser(daoUser dao.User) domain.User {
	return domain.User{
		ID:       daoUser.ID,
		Email:    daoUser.Email,
		Password: daoUser.Password,
	}
}
