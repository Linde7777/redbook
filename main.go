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

	corsConfig := cors.DefaultConfig()
	router.Use(cors.New(corsConfig))

	return router
}

func initUserRoutes(db *gorm.DB, router *gin.Engine) {
	ud := dao.NewUserDAO(db)
	ur := repository.NewUserRepository(ud)
	us := service.NewUserService(ur)
	uh := web.NewUserHandler(us)
	uh.RegisterRoutes(router)

}
