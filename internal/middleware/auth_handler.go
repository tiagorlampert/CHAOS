package middleware

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/tiagorlampert/CHAOS/entities"
	jwtUtil "github.com/tiagorlampert/CHAOS/internal/utils/jwt"
	"github.com/tiagorlampert/CHAOS/services/user"
	"net/http"
)

type authHandler struct {
	UserService user.Service
}

func newAuthHandler(userService user.Service) *authHandler {
	return &authHandler{UserService: userService}
}

func (a authHandler) payloadFuncHandler(data interface{}) jwt.MapClaims {
	if v, ok := data.(*entities.User); ok {
		return jwt.MapClaims{
			authorizedKey:       true,
			jwtUtil.IdentityKey: v.Username,
		}
	}
	return jwt.MapClaims{}
}

func (a authHandler) identityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return &entities.User{
		Username: claims[jwtUtil.IdentityKey].(string),
	}
}

func (a authHandler) authenticatorHandler(c *gin.Context) (interface{}, error) {
	var user entities.User
	if err := c.ShouldBind(&user); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	if a.UserService.Login(user.Username, user.Password) {
		return &entities.User{
			Username: user.Username,
		}, nil
	}
	return nil, jwt.ErrFailedAuthentication
}

func (a authHandler) unauthorizedHandler(c *gin.Context, code int, message string) {
	c.HTML(http.StatusUnauthorized, "login.html", gin.H{"unauthorized": true})
	return
}

func (a authHandler) logoutResponseHandler(c *gin.Context, code int) {
	c.HTML(http.StatusOK, "login.html", gin.H{"unauthorized": true})
	return
}
