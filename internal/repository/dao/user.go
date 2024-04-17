package dao

import (
	"context"
	"gorm.io/gorm"
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
	ID       uint64 `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string

	// 理论上unix时间应该用uint64，但是time.UnixSec()返回的是int64
	CreateTime int64
	UpdateTime int64
}

func (dao *UserDAO) Insert(ctx context.Context, user *User) error {
	return dao.db.WithContext(ctx).Create(user).Error
}

func (dao *UserDAO) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return &user, err
}
