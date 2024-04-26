package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"testing"
)

func TestGORMUserDAO_Insert(t *testing.T) {
	testCases := []struct {
		name             string
		mock             func(t *testing.T) *sql.DB
		ctx              context.Context
		inputUser        *User
		expectedHTTPCode int
		expectedErr      error
	}{
		{
			name: "成功",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mockResult := sqlmock.NewResult(1, 1)
				mock.ExpectExec("INSERT INTO .*").WillReturnResult(mockResult)
				return db
			},
			ctx: context.Background(),
			inputUser: &User{
				Email:       sql.NullString{},
				Password:    "",
				PhoneNumber: sql.NullString{},
			},
			expectedHTTPCode: http.StatusOK,
			expectedErr:      nil,
		},
		{
			name: "失败",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectExec("INSERT INTO .*").WillReturnError(errors.New("test error"))
				return db
			},
			ctx: context.Background(),
			inputUser: &User{
				Email:       sql.NullString{},
				Password:    "",
				PhoneNumber: sql.NullString{},
			},
			expectedHTTPCode: http.StatusInternalServerError,
			expectedErr:      errors.New("test error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqlDB := tc.mock(t)
			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqlDB,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				DisableAutomaticPing:   true,
				SkipDefaultTransaction: true,
			})
			assert.NoError(t, err)
			mockDAO := NewGORMUserDAO(db)
			httpCode, err := mockDAO.Insert(tc.ctx, tc.inputUser)
			assert.Equal(t, tc.expectedHTTPCode, httpCode)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
