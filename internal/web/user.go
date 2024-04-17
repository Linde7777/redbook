package web

import (
	"github.com/gin-gonic/gin"
	"main/internal/domain"
	"main/internal/service"
	"net/http"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}

func (h *UserHandler) RegisterRoutes(e *gin.Engine) {
	ug := e.Group("/user")
	ug.POST("/signup", h.Signup)
}

func (h *UserHandler) Signup(c *gin.Context) {
	type ReqSignup struct {
		Email string `json:"email" binding:"required,email"`
		Name  string `json:"name" binding:"required"`

		// bcrypt算法最大长度72字节
		Password string `json:"password" binding:"required,min=8,max=72"`

		ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=Password"`
	}

	var req ReqSignup
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.svc.Signup(c, &domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.String(http.StatusConflict, err.Error())
	}

	c.String(http.StatusOK, "sign up success")
}
