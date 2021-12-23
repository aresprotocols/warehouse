package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type AdminClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

const TokenExpireDuration = time.Hour * 2

func GenToken(username string, signSecret []byte) (string, error) {

	c := AdminClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    "get-price",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	return token.SignedString(signSecret)
}

func ParseToken(tokenString string, signSecret []byte) (*AdminClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &AdminClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return signSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*AdminClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
