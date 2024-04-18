package web

import (
	"github.com/gin-gonic/gin"
	"main/internal/domain"
	"main/internal/service"
	"main/internal/web/middlewares"
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
	ug.POST("/login", h.Login)
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
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	httpCode, err := h.svc.Signup(c, &domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.String(httpCode, err.Error())
	}
	c.String(httpCode, "sign up success")
}

// const KeyUserID = "userID"

func (h *UserHandler) Login(c *gin.Context) {
	type ReqLogin struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	var req ReqLogin
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	user, httpCode, err := h.svc.Login(c, &domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.String(httpCode, err.Error())
	}

	err = middlewares.SetJWT(c, user.ID)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	c.String(httpCode, "login success")
}
