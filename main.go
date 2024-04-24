package main

func main() {
	router := InitWebServer()
	router.Run(":8080")
}

//func initDB() *gorm.DB {
//	db, err := gorm.Open(mysql.Open("root:123456@tcp(localhost:3306)/redbook?charset=utf8mb4&parseTime=True&loc=Local"))
//	if err != nil {
//		panic(err)
//	}
//	return db
//}
//
//func initRedis() redis.Cmdable {
//	client := redis.NewClient(&redis.Options{
//		Addr: "localhost:6379",
//	})
//	return client
//}
//
//func initWebServer() *gin.Engine {
//	router := gin.Default()
//
//	corsConfig := cors.Config{
//		AllowHeaders:  []string{"Content-Type", "Authorization"},
//		ExposeHeaders: []string{middlewares.KeyBackendJWTHeader},
//		MaxAge:        12 * time.Hour,
//		AllowOriginFunc: func(origin string) bool {
//			if strings.HasPrefix(origin, "http://localhost") {
//				return true
//			}
//			return strings.Contains(origin, "mycompany.com")
//		},
//	}
//	router.Use(cors.New(corsConfig))
//
//	//v1 := e.Group("/v1")
//	//v1UserGroup := v1.Group("/user")
//	//v1UserGroup.POST("/signup", h.Signup)
//	//v1UserGroup.POST("/login-by-password", h.LoginByPassword)
//	//v1UserGroup.POST("/send-login-sms-auth-code", h.SendLoginSMSAuthCode)
//	//v1UserGroup.POST("/login-by-sms-auth-code", h.LoginBySMSAuthCode)
//	builder := middlewares.NewLoginMiddlewareBuilder()
//	builder.IgnorePath("/v1/user/signup",
//		"/v1/user/login-by-password", "v1/user/send-login-sms-auth-code",
//		"v1/user/login-by-sms-auth-code")
//	router.Use(builder.CheckLogin())
//
//	return router
//}
//
//func initUserRoutes(db *gorm.DB, redisCMD redis.Cmdable, router *gin.Engine) {
//	userDAO := dao.NewUserDAO(db)
//	userCache := cache.NewRedisUserCache(redisCMD, 15*time.Minute)
//	userRepo := repository.NewUserRepository(userDAO, userCache)
//	userService := service.NewUserService(userRepo)
//
//	authCodeCache := cache.NewRedisAuthCodeCache(redisCMD)
//	authCodeRepo := repository.NewAuthCodeRepository(authCodeCache)
//	// todo: 解引用问题
//	authCodeService := service.NewAuthCodeService(*authCodeRepo, nil)
//
//	userHandler := web.NewUserHandler(userService, authCodeService)
//	userHandler.RegisterRoutes(router)
//}
