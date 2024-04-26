package repository

import (
	"context"
	"database/sql"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"main/internal/domain"
	"main/internal/repository/cache"
	"main/internal/repository/dao"
	"net/http"
	"testing"
)

func TestUserRepositoryWithCache_SearchUserByEmail(t *testing.T) {
	testCases := []struct {
		name             string
		mock             func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDAO)
		inputEmail       string
		expectedUser     domain.User
		expectedHTTPCode int
	}{
		{
			name: "缓存没数据",
			mock: func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDAO) {
				mockCache := cache.NewMockUserCache(ctrl)
				mockCache.EXPECT().GetUserByEmail(gomock.Any(), "abc@gmail.com").Return(domain.User{
					ID:    1,
					Email: "abc@gmail.com",
				}, http.StatusOK, nil)
				return mockCache, nil
			},
			inputEmail: "abc@gmail.com",
			expectedUser: domain.User{
				ID:    1,
				Email: "abc@gmail.com",
			},
			expectedHTTPCode: http.StatusOK,
		},
		{
			name: "缓存没数据，数据库有数据",
			mock: func(ctrl *gomock.Controller) (cache.UserCache, dao.UserDAO) {
				mockCache := cache.NewMockUserCache(ctrl)
				mockCache.EXPECT().GetUserByEmail(gomock.Any(), "abc@gmail.com").
					Return(domain.User{}, http.StatusNotFound, redis.Nil)

				mockDAO := dao.NewMockUserDAO(ctrl)
				mockDAO.EXPECT().SearchUserByEmail(gomock.Any(), "abc@gmail.com").
					Return(dao.User{
						ID: 1,
						Email: sql.NullString{
							String: "abc@gmail.com",
							Valid:  true,
						},
					}, true, http.StatusOK, nil)

				mockCache.EXPECT().SetUserByEmail(gomock.Any(), domain.User{
					ID:    1,
					Email: "abc@gmail.com",
				}).Return(http.StatusOK, nil)

				return mockCache, mockDAO
			},
			inputEmail: "abc@gmail.com",
			expectedUser: domain.User{
				ID:    1,
				Email: "abc@gmail.com",
			},
			expectedHTTPCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockCache, mockDAO := tc.mock(ctrl)
			repo := NewUserRepositoryWithCache(mockDAO, mockCache)
			user, ok, httpCode, err := repo.SearchUserByEmail(context.Background(), tc.inputEmail)
			assert.True(t, ok)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, httpCode)
			assert.Equal(t, tc.inputEmail, user.Email)
			assert.Equal(t, tc.expectedUser, user)
			assert.Equal(t, tc.expectedHTTPCode, httpCode)
		})
	}
}
