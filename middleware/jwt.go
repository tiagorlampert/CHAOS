package middleware

import (
	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/tiagorlampert/CHAOS/entities"
	jwtUtil "github.com/tiagorlampert/CHAOS/internal/utilities/jwt"
	"github.com/tiagorlampert/CHAOS/services"
	"net/http"
	"time"
)

type JWT struct {
	*jwt.GinJWTMiddleware
}

func NewJWTMiddleware(secretKey string, userService services.User) (*JWT, error) {
	middleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:         "chaos",
		Key:           []byte(secretKey),
		Timeout:       time.Hour,
		MaxRefresh:    time.Hour,
		IdentityKey:   jwtUtil.IdentityKey,
		SendCookie:    true,
		TokenLookup:   "cookie:jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*entities.User); ok {
				return jwt.MapClaims{
					"authorized":        true,
					jwtUtil.IdentityKey: v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &entities.User{
				Username: claims[jwtUtil.IdentityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var user entities.User
			if err := c.ShouldBind(&user); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			if userService.Login(user.Username, user.Password) {
				return &entities.User{
					Username: user.Username,
				}, nil
			}
			return nil, jwt.ErrFailedAuthentication
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"unauthorized": true})
			return
		},
		LogoutResponse: func(c *gin.Context, code int) {
			c.HTML(http.StatusOK, "login.html", gin.H{"unauthorized": true})
			return
		},
	})
	if err != nil {
		return nil, err
	}
	return &JWT{middleware}, nil
}

func (j *JWT) AuthAdmin(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	if claims[jwtUtil.IdentityKey] != jwtUtil.IdentityAdminUser {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
}
