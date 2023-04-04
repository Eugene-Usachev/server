package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"os"
)

func ParseAccessToken(tokenString string) (uint, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", errors.New("invalid signing method")
		}

		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return 0, err
	}

	tokenClaims, ok := token.Claims.(*AccessClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	return tokenClaims.UserID, nil
}

func ParseLongliveToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &LongliveClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", errors.New("invalid signing method")
		}

		return []byte(os.Getenv("JWT_SECRET_KEY_FOR_LONGLIVE_TOKEN")), nil
	})
	if err != nil {
		return "", err
	}

	tokenClaims, ok := token.Claims.(*LongliveClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	return tokenClaims.Password, nil
}
