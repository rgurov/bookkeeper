package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWT struct {
	secret []byte
}

func NewJWT(secret string) *JWT {
	return &JWT{
		secret: []byte(secret),
	}
}

func (j *JWT) Sign(m jwt.MapClaims, exp time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	m["exp"] = time.Now().Add(exp).Unix()
	token.Claims = m
	tokenString, err := token.SignedString(j.secret)

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *JWT) Parse(strToken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(strToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("error while parsing token")
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("couldn't parse jwt claims")
	}

	exp := claims["exp"].(float64)
	if int64(exp) < time.Now().Local().Unix() {
		return nil, errors.New("token expired")
	}

	return claims, nil
}
