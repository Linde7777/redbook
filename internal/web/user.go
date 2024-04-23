package web

import (
	"github.com/gin-gonic/gin"
	"main/internal/domain"
	"main/internal/service"
	"main/internal/web/middlewares"
	"net/http"
)

type UserHandler struct {
	userService     *service.UserService
	authCodeService *service.AuthCodeService
}

func NewUserHandler(userService *service.UserService, authCodeService *service.AuthCodeService) *UserHandler {
	return &UserHandler{
		userService:     userService,
		authCodeService: authCodeService,
	}
}

func (h *UserHandler) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/user")
	group.POST("/signup", h.Signup)
	group.POST("/login-by-password", h.LoginByPassword)
	group.POST("/send-login-sms-auth-code", h.SendLoginSMSAuthCode)
	group.POST("/login-by-sms-auth-code", h.LoginBySMSAuthCode)
}

func (h *UserHandler) Signup(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
		Name  string `json:"name" binding:"required"`

		// bcrypt算法最大长度72字节
		Password string `json:"password" binding:"required,min=8,max=72"`

		ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=Password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	user, httpCode, err := h.userService.Signup(c, &domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.String(httpCode, err.Error())
	}
	c.JSON(httpCode, user)
}

// const KeyUserID = "userID"

func (h *UserHandler) LoginByPassword(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	user, httpCode, err := h.userService.LoginByPassword(c, &domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.String(httpCode, err.Error())
	}

	err = middlewares.SetJWT(c, user.ID, c.GetHeader("User-Agent"))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	c.String(httpCode, "login success")
}

func (h *UserHandler) SendLoginSMSAuthCode(c *gin.Context) {
	var req struct {
		// 关于tag e164: https://github.com/go-playground/validator#:~:text=Datetime-,e164,-e164%20formatted%20phone
		PhoneNumber string `json:"phone_number" binding:"required,e164"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	const businessName = "login"
	httpCode, err := h.authCodeService.SendAuthCode(c, businessName, req.PhoneNumber)
	if err != nil {
		c.String(httpCode, err.Error())
	}
	c.String(httpCode, "send auth code success")
}

func (h *UserHandler) LoginBySMSAuthCode(c *gin.Context) {
	var req struct {
		PhoneNumber string `json:"phone_number" binding:"required,e164"`
		AuthCode    string `json:"auth_code" binding:"required,len=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	const businessName = "login"
	httpCode, err := h.authCodeService.VerifyAuthCode(c, businessName, req.PhoneNumber, req.AuthCode)
	if err != nil {
		c.String(httpCode, err.Error())
	}

	user, httpCode, err := h.userService.SearchOrCreateUserByPhoneNumber(c, req.PhoneNumber)
	if err != nil {
		c.String(httpCode, err.Error())
	}

	err = middlewares.SetJWT(c, user.ID, c.GetHeader("User-Agent"))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}

	c.String(httpCode, "login success")
}
