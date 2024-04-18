package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"main/internal/repository"
	"main/internal/repository/dao"
	"main/internal/service"
	"main/internal/web"
	"main/internal/web/middlewares"
	"strings"
	"time"
)

func main() {
	db := initDB()
	router := initWebServer()

	initUserRoutes(db, router)

	router.Run(":8080")
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:123456@tcp(localhost:3306)/redbook?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(err)
	}
	return db
}

func initWebServer() *gin.Engine {
	router := gin.Default()

	corsConfig := cors.Config{
		AllowHeaders:  []string{"Content-Type", "Authorization"},
		ExposeHeaders: []string{middlewares.KeyBackendJWTHeader},
		MaxAge:        12 * time.Hour,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "mycompany.com")
		},
	}
	router.Use(cors.New(corsConfig))

	builder := middlewares.NewLoginMiddlewareBuilder()
	builder.IgnorePath("/v1/user/signup", "/v1/user/login")
	router.Use(builder.CheckLogin())

	return router
}

func initUserRoutes(db *gorm.DB, router *gin.Engine) {
	ud := dao.NewUserDAO(db)
	ur := repository.NewUserRepository(ud)
	us := service.NewUserService(ur)
	uh := web.NewUserHandler(us)
	uh.RegisterRoutes(router)
}
