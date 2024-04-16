package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"main/internal/web"
)

func main() {
	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	router.Use(cors.New(corsConfig))

	uh := &web.UserHandler{}
	uh.RegisterRoutes(router)

	router.Use()
}
