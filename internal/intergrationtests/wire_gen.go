// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package intergrationtests

import (
	"github.com/gin-gonic/gin"
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

// Injectors from wire.go:

func InitWebServer() *gin.Engine {
	v := ioc.InitGinMiddlewares()
	db := InitDB()
	userDAO := dao.NewUserDAO(db)
	cmdable := InitRedis()
	userCache := cache.NewRedisUserCache(cmdable)
	userRepository := repository.NewUserRepositoryWithCache(userDAO, userCache)
	userService := service.NewUserServiceV1(userRepository)
	authCodeCache := cache.NewRedisAuthCodeCache(cmdable)
	authCodeRepository := repository.NewAuthCodeRepositoryWithCache(authCodeCache)
	smsService := ioc.InitSMSService()
	authCodeService := service.NewAuthCodeServiceV1(authCodeRepository, smsService)
	userHandler := web.NewUserHandler(userService, authCodeService)
	engine := ioc.InitGinWebServer(v, userHandler)
	return engine
}

// wire.go:

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