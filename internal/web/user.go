package web

import (
	"github.com/gin-gonic/gin"
	"main/internal/domain"
	"main/internal/service"
	"main/internal/web/middlewares"
	"net/http"
)

// UserHandler 仅仅是为了注册路由，不会被其他人调用，所以这里目前不需要接口
type UserHandler struct {
	userService     service.UserService
	authCodeService service.AuthCodeService
}

func NewUserHandler(userService service.UserService, authCodeService service.AuthCodeService) *UserHandler {
	return &UserHandler{
		userService:     userService,
		authCodeService: authCodeService,
	}
}

func (h *UserHandler) RegisterRoutes(e *gin.Engine) {
	v1 := e.Group("/v1")
	v1UserGroup := v1.Group("/user")
	v1UserGroup.POST("/signup", h.Signup)
	v1UserGroup.POST("/login-by-password", h.LoginByPassword)
	v1UserGroup.POST("/send-login-sms-auth-code", h.SendLoginSMSAuthCode)
	v1UserGroup.POST("/login-by-sms-auth-code", h.LoginBySMSAuthCode)
}

type ReqSignup struct {
	Email string `json:"email" binding:"required,email"`

	// bcrypt算法最大长度72字节
	Password string `json:"password" binding:"required,min=8,max=72"`

	ConfirmPassword string `json:"confirmPassword" binding:"required,eqfield=Password"`
}

func (h *UserHandler) Signup(c *gin.Context) {
	var req ReqSignup
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	_, httpCode, err := h.userService.Signup(c, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.String(httpCode, err.Error())
	}
	c.String(httpCode, "signup success")
}

// const KeyUserID = "userID"

func (h *UserHandler) LoginByPassword(c *gin.Context) {
	type ReqLoginByPassword struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	var req ReqLoginByPassword
	if err := c.ShouldBindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	user, httpCode, err := h.userService.LoginByPassword(c, domain.User{
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
	type ReqSendLoginSMSAuthCode struct {
		// 关于tag e164: https://github.com/go-playground/validator#:~:text=Datetime-,e164,-e164%20formatted%20phone
		PhoneNumber string `json:"phone_number" binding:"required,e164"`
	}

	var req ReqSendLoginSMSAuthCode
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
	type ReqLoginBySMSAuthCode struct {
		PhoneNumber string `json:"phone_number" binding:"required,e164"`
		AuthCode    string `json:"auth_code" binding:"required,len=6"`
	}

	var req ReqLoginBySMSAuthCode
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
