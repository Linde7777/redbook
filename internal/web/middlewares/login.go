package middlewares

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"main/internal/web"
	"net/http"
)

type LoginMiddlewareBuilder struct {
}

func (b *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sess := sessions.Default(ctx)
		if sess.Get(web.KeyUserID) == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
