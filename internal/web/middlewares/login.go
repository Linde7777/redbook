package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

type LoginMiddlewareBuilder struct {
	ignoredPaths map[string]struct{}
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (b *LoginMiddlewareBuilder) IgnorePath(paths ...string) {
	for i := 0; i < len(paths); i++ {
		b.ignoredPaths[paths[i]] = struct{}{}
	}
}

const (
	JWTSecret = "95osj3fUD7fo0mlYdDbncXz4VD2igvf0"

	// KeyBackendJWTHeader 是后端返回给前端JWT Token时，要存放的位置
	KeyBackendJWTHeader = "X-JWT-Token"
)

func (b *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if _, ok := b.ignoredPaths[path]; ok {
			return
		}

		token, claim, err := getTokenAndCustomClaims(c)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !validateJWT(token, claim) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

func getTokenAndCustomClaims(c *gin.Context) (*jwt.Token, customClaims, error) {
	authStr := c.GetHeader("Authorization")
	if authStr == "" {
		return nil, customClaims{}, errors.New("no Authorization header")
	}

	tempArr := strings.Split(authStr, " ")
	if len(tempArr) != 2 {
		return nil, customClaims{}, errors.New("authorization header format error")
	}
	tokenStr := tempArr[1]
	var claims customClaims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSecret), nil
	})
	return token, claims, err
}

func validateJWT(token *jwt.Token, claims customClaims) bool {
	return token == nil || !token.Valid || claims.ExpiresAt.Before(time.Now())
}

type customClaims struct {
	jwt.RegisteredClaims
	UserID uint64 `json:"user_id"`
}

func SetJWT(c *gin.Context, userID uint64) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, customClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})
	tokenStr, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		return err
	}
	c.Header(KeyBackendJWTHeader, tokenStr)
	return nil
}

func GetUserID(c *gin.Context) uint64 {
	_, claim, _ := getTokenAndCustomClaims(c)
	return claim.UserID
}
