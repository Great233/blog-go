package utils

import (
	"blog/config"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"time"
)

type User struct {
	Username string `json:"username"`
}

type Claims struct {
	User *User
	jwt.StandardClaims
}

func GenerateJsonWebToken(user *User) (string, error) {
	now := time.Now().Unix()
	id := uuid.NewV4()

	claims := Claims{
		user,
		jwt.StandardClaims{
			Audience:  "Go-gin-blog",
			ExpiresAt: now + 30 * 60,
			Id:        id.String(),
			IssuedAt:  now,
			Issuer:    "blog",
			NotBefore: now,
			Subject:   "Session",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	key, err := base64.StdEncoding.DecodeString(config.App.Auth.Key)
	if err != nil {
		return "", err
	}
	return token.SignedString(key)
}

func ParseJsonWebToken(token string) (*User, error) {
	claims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return base64.StdEncoding.DecodeString(config.App.Auth.Key)
	})

	if err != nil {
		return nil, err
	}

	var parsed *Claims
	var ok bool
	if parsed, ok = claims.Claims.(*Claims); !ok {
		return nil, errors.New("token parse error")
	}
	err = parsed.Valid()
	if err != nil {
		return nil, err
	}

	return parsed.User, nil
}
