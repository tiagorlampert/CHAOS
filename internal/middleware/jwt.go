package middleware

import (
	"bytes"
	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	jwtUtil "github.com/tiagorlampert/CHAOS/internal/utils/jwt"
	"github.com/tiagorlampert/CHAOS/services/auth"
	"github.com/tiagorlampert/CHAOS/services/user"
	"net/http"
	"time"
)

const (
	nameToDisplay   = "chaos"
	tokenLookup     = "cookie:jwt"
	tokenHeaderName = "Bearer"
	authorizedKey   = "authorized"
)

type JWT struct {
	*jwt.GinJWTMiddleware
}

func NewJwtMiddleware(
	authService auth.Service,
	userService user.Service,
) *JWT {
	secret, err := authService.GetSecret()
	if err != nil {
		panic(err)
	}

	authHandler := newAuthHandler(userService)

	m, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:           nameToDisplay,
		Key:             bytes.NewBufferString(secret).Bytes(),
		Timeout:         time.Hour,
		MaxRefresh:      time.Hour,
		IdentityKey:     jwtUtil.IdentityKey,
		TokenLookup:     tokenLookup,
		TokenHeadName:   tokenHeaderName,
		SendCookie:      true,
		TimeFunc:        time.Now,
		PayloadFunc:     authHandler.payloadFuncHandler,
		IdentityHandler: authHandler.identityHandler,
		Authenticator:   authHandler.authenticatorHandler,
		Unauthorized:    authHandler.unauthorizedHandler,
		LogoutResponse:  authHandler.logoutResponseHandler,
	})
	if err != nil {
		panic(err)
	}
	return &JWT{m}
}

func (j *JWT) AuthAdmin(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	if claims[jwtUtil.IdentityKey] != jwtUtil.IdentityAdminUser {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
}
