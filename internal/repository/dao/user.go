package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"net/http"
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

type User struct {
	ID          uint64
	UUID        string
	Email       string
	Password    string
	PhoneNumber string

	// 理论上unix时间应该用uint64，但是time.UnixSec()返回的是int64
	CreateTime int64
	UpdateTime int64
}

func (dao *UserDAO) Insert(ctx context.Context, user *User) (httpCode int, err error) {
	err = dao.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (dao *UserDAO) SearchUserByEmail(ctx context.Context, email string) (user User, ok bool, httpCode int, err error) {
	err = dao.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	switch {
	case err != nil && errors.Is(err, gorm.ErrRecordNotFound):
		return user, false, http.StatusNotFound, nil
	case err != nil:
		return user, false, http.StatusInternalServerError, err
	default:
		return user, true, http.StatusOK, nil
	}
}

func (dao *UserDAO) SearchUserByPhoneNumber(ctx context.Context, number string) (user User, ok bool, httpCode int, err error) {
	err = dao.db.WithContext(ctx).Where("phone_number = ?", number).First(&user).Error
	switch {
	case err != nil && errors.Is(err, gorm.ErrRecordNotFound):
		return user, false, http.StatusNotFound, nil
	case err != nil:
		return user, false, http.StatusInternalServerError, err
	default:
		return user, true, http.StatusOK, nil
	}
}
