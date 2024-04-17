package web

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
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

	store := cookie.NewStore([]byte("secret"))
	builder := middlewares.LoginMiddlewareBuilder{}
	e.Use(sessions.Sessions("ssid", store), builder.CheckLogin())

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

	err := h.svc.Signup(c, &domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.String(http.StatusConflict, err.Error())
	}

	c.String(http.StatusOK, "sign up success")
}

const KeyUserID = "userID"

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

	user, err := h.svc.Login(c, &domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.String(http.StatusConflict, err.Error())
	}

	sess := sessions.Default(c)
	sess.Set(KeyUserID, user.ID)
	// todo: MaxAge从配置文件读取
	sess.Options(sessions.Options{
		MaxAge: 15 * 60,
	})
	err = sess.Save()
	if err != nil {
		c.String(http.StatusConflict, err.Error())
		return
	}

	c.String(http.StatusOK, "login success")
}
