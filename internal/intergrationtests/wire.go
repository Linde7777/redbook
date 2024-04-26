//go:build wireinject

package intergrationtests

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"main/internal/repository"
	"main/internal/repository/cache"
	"main/internal/repository/dao"
	"main/internal/service"
	"main/internal/web"
	"main/ioc"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:123456@tcp(localhost:3306)/redbook?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(err)
	}
	return db
}

func InitRedis() redis.Cmdable {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return client
}

func InitWebServer() *gin.Engine {
	wire.Build(InitDB, InitRedis,
		dao.NewUserDAO,
		cache.NewRedisUserCache, cache.NewRedisAuthCodeCache,
		repository.NewUserRepositoryWithCache, repository.NewAuthCodeRepositoryWithCache,
		ioc.InitSMSService, service.NewUserServiceV1, service.NewAuthCodeServiceV1,
		web.NewUserHandler,
		ioc.InitGinMiddlewares, ioc.InitGinWebServer)
	return gin.Default()
}
