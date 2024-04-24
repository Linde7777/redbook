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

func InitWebServer() *gin.Engine {
	wire.Build(ioc.InitDB, ioc.InitRedis,
		dao.NewUserDAO,
		cache.NewRedisUserCache, cache.NewRedisAuthCodeCache,
		repository.NewUserRepository, repository.NewAuthCodeRepository,
		ioc.InitSMSService, service.NewUserService, service.NewAuthCodeService,
		web.NewUserHandler,
		ioc.InitGinMiddlewares, ioc.InitGinWebServer)
	return gin.Default()
}
