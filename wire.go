//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"main/internal/repository"
	"main/internal/repository/cache"
	"main/internal/repository/dao"
	"main/internal/service"
	"main/internal/web"
	"main/ioc"
)

var providerSet wire.ProviderSet = wire.NewSet(
	ioc.InitDB, ioc.InitRedis,
	dao.NewGORMUserDAO,

	wire.Bind(new(dao.UserDAO), new(*dao.GORMUserDAO)))

func InitWebServer() *gin.Engine {
	wire.Build(ioc.InitDB, ioc.InitRedis,
		dao.NewGORMUserDAO,
		cache.NewRedisUserCache, cache.NewRedisAuthCodeCache,
		repository.NewUserRepositoryWithCache, repository.NewAuthCodeRepositoryWithCache,
		ioc.InitTencentSMSService, service.NewUserServiceV1, service.NewAuthCodeServiceV1,
		web.NewUserHandler,
		ioc.InitGinMiddlewares, ioc.InitGinWebServer)
	//wire.Build(providerSet)
	return gin.Default()
}
