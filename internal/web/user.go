package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
}

func (uh *UserHandler) RegisterRoutes(e *gin.Engine) {
	ug := e.Group("/user")
	ug.POST("/signup", uh.Signup)
}

func (uh *UserHandler) Signup(c *gin.Context) {
	type ReqSignup struct {
		Email           string `json:"email" binding:"required,email"`
		Password        string `json:"password" binding:"required,min=8,max=32"`
		ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=Password"`
	}

	var req ReqSignup
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}
