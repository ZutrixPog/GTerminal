package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

// Not Safe
var secretKey = []byte("MySecretKey")

func GenerateToken(user string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(1 * time.Hour)
	claims["user"] = user

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString,
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodECDSA)
			if !ok {
				return nil, errors.New("couldn't Authorize user")
			}

			return "", nil
		})

	if err != nil && token.Valid {
		return errors.New("couldn't Authorize user")
	}

	return nil
}
