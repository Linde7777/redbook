package dao

import (
	"context"
	"database/sql"
	"errors"
	"gorm.io/gorm"
	"net/http"
	"sync"
)

type UserDAO interface {
	Insert(ctx context.Context, user *User) (httpCode int, err error)
	SearchUserByEmail(ctx context.Context, email string) (user User, ok bool, httpCode int, err error)
	SearchUserByPhoneNumber(ctx context.Context, number string) (user User, ok bool, httpCode int, err error)
}

var once sync.Once
var gormUserDAO *GORMUserDAO

type GORMUserDAO struct {
	db *gorm.DB
}

// NewGORMUserDAO 为了适配wire，只能返回接口，而不是返回具体实现
func NewGORMUserDAO(db *gorm.DB) *GORMUserDAO {
	once.Do(func() {
		gormUserDAO = &GORMUserDAO{
			db: db,
		}
	})
	return gormUserDAO
}

type User struct {
	ID          uint64
	Email       sql.NullString //用手机注册的用户，刚注册时没填邮箱
	Password    string
	PhoneNumber sql.NullString //用邮箱注册的用户，刚注册时没填手机号

	// 理论上unix时间应该用uint64，但是time.UnixSec()返回的是int64
	CreateTime int64
	UpdateTime int64
}

func (dao *GORMUserDAO) Insert(ctx context.Context, user *User) (httpCode int, err error) {
	err = dao.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (dao *GORMUserDAO) SearchUserByEmail(ctx context.Context, email string) (user User, ok bool, httpCode int, err error) {
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

func (dao *GORMUserDAO) SearchUserByPhoneNumber(ctx context.Context, number string) (user User, ok bool, httpCode int, err error) {
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
