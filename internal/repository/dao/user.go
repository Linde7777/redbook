package dao

import (
	"context"
	"gorm.io/gorm"
)

type UserDAO struct {
	db *gorm.DB
}

type User struct {
	ID         int64  `gorm:"primaryKey"`
	Email      string `gorm:"unique"`
	Password   string
	CreateTime int64
	UpdateTime int64
}

func (dao *UserDAO) Insert(ctx context.Context, user *User) error {
	return dao.db.WithContext(ctx).Create(user).Error
}
