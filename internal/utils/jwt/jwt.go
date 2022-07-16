package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

const (
	IdentityKey         = "user"
	IdentityAdminUser   = "admin"
	IdentityDefaultUser = "default"
)

func NewToken(key string, identity ...string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims[IdentityKey] = IdentityAdminUser
	claims["exp"] = time.Now().AddDate(1, 0, 0).Unix()
	if len(identity) != 0 {
		claims[IdentityKey] = identity[0]
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
