package ioc

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"main/internal/web"
	"main/internal/web/middlewares"
	"strings"
	"time"
)

// 不要和wire.go里面的InitWebServer重复
func InitGinWebServer(middlewares []gin.HandlerFunc, userHandler *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(middlewares...)
	userHandler.RegisterRoutes(server)
	return server
}

func InitGinMiddlewares(redisCmd redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		func(context *gin.Context) {
			keyFunc := func(ctx *gin.Context) string {
				return fmt.Sprintf("ip-limiter:%s", ctx.ClientIP())
			}
			builder := middlewares.
				NewRedisSlidingWindowsLimiterBuilder(
					redisCmd, 10, 60, keyFunc)
			builder.Build()(context)
		},

		cors.New(cors.Config{
			AllowHeaders:  []string{"Content-Type", "Authorization"},
			ExposeHeaders: []string{middlewares.KeyBackendJWTHeader},
			MaxAge:        12 * time.Hour,
			AllowOriginFunc: func(origin string) bool {
				if strings.HasPrefix(origin, "http://localhost") {
					return true
				}
				return strings.Contains(origin, "mycompany.com")
			},
		}),

		func(context *gin.Context) {
			builder := middlewares.NewLoginMiddlewareBuilder()

			// builder.Build()返回的是一个gin.HandlerFunc
			// gin.HandlerFunc接受一个*gin.Context参数
			// 等效于gin.HandlerFunc(context)
			builder.Build()(context)

			builder.IgnorePath("/v1/user/signup",
				"/v1/user/login-by-password", "v1/user/send-login-sms-auth-code",
				"v1/user/login-by-sms-auth-code")
		},
	}
}
