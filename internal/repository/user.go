package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/redis/go-redis/v9"
	"main/internal/domain"
	"main/internal/repository/cache"
	"main/internal/repository/dao"
	"net/http"
)

type UserRepository interface {
	Create(ctx context.Context, inputDomainUser *domain.User) (user domain.User, httpCode int, err error)
	SearchUserByEmail(ctx context.Context, email string) (user domain.User, httpCode int, err error)
	SearchUserByPhoneNumber(ctx context.Context, phoneNumber string) (user domain.User, ok bool, httpCode int, err error)
	toDomainUser(daoUser dao.User) domain.User
	toDaoUser(domainUser domain.User) dao.User
}

type UserRepositoryWithCache struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

// NewUserRepositoryWithCache 为了适配wire，只能返回接口，而不是返回具体实现
func NewUserRepositoryWithCache(dao dao.UserDAO, cache cache.UserCache) UserRepository {
	return &UserRepositoryWithCache{
		dao:   dao,
		cache: cache,
	}
}

func (repo *UserRepositoryWithCache) Create(ctx context.Context, inputDomainUser *domain.User) (user domain.User, httpCode int, err error) {
	daoUser := repo.toDaoUser(*inputDomainUser)
	httpCode, err = repo.dao.Insert(ctx, &daoUser)
	if err != nil {
		return domain.User{}, httpCode, err
	}
	return repo.toDomainUser(daoUser), http.StatusOK, nil
}

func (repo *UserRepositoryWithCache) SearchUserByEmail(ctx context.Context, email string) (user domain.User, httpCode int, err error) {
	domainUser, httpCode, err := repo.cache.GetUserByEmail(ctx, email)
	switch {
	case err == nil:
		return domainUser, httpCode, nil
		// todo: 缓存穿透
	case errors.Is(err, redis.Nil):
		break
	}

	daoUser, ok, httpCode, err := repo.dao.SearchUserByEmail(ctx, email)
	if err != nil || !ok {
		return domain.User{}, httpCode, err
	}
	domainUser = repo.toDomainUser(daoUser)
	httpCode, err = repo.cache.SetUserByEmail(ctx, domainUser)
	return domainUser, httpCode, err
}

func (repo *UserRepositoryWithCache) SearchUserByPhoneNumber(ctx context.Context,
	phoneNumber string) (user domain.User, ok bool, httpCode int, err error) {

	daoUser, ok, httpCode, err := repo.dao.SearchUserByPhoneNumber(ctx, phoneNumber)
	if err != nil || !ok {
		return domain.User{}, false, httpCode, err
	}
	return repo.toDomainUser(daoUser), true, httpCode, nil
}

func (repo *UserRepositoryWithCache) toDomainUser(daoUser dao.User) domain.User {
	return domain.User{
		ID:          daoUser.ID,
		Email:       daoUser.Email.String,
		PhoneNumber: daoUser.PhoneNumber.String,
		Password:    daoUser.Password,
	}
}

func (repo *UserRepositoryWithCache) toDaoUser(domainUser domain.User) dao.User {
	return dao.User{
		ID: domainUser.ID,
		Email: sql.NullString{
			String: domainUser.Email,
			Valid:  domainUser.Email != "",
		},
		Password: domainUser.Password,
		PhoneNumber: sql.NullString{
			String: domainUser.PhoneNumber,
			Valid:  domainUser.PhoneNumber != "",
		},
	}
}
