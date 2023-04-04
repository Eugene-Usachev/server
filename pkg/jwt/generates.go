package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

const (
	AccessTokenLive   = time.Minute * 15
	RefreshTokenLive  = time.Hour * 24 * 31
	LongliveTokenLive = time.Hour * 24 * 31
)

func NewRefreshToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(RefreshTokenLive).Unix(),
		IssuedAt:  time.Now().Unix(),
	})
	RefreshToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return RefreshToken, err
}

func NewAccessToken(id uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &AccessClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(AccessTokenLive).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		id,
	})
	AccessToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return AccessToken, err
}

func NewLongliveToken(passwordHash string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &LongliveClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(LongliveTokenLive).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		passwordHash,
	})
	AccessToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY_FOR_LONGLIVE_TOKEN")))
	if err != nil {
		return "", err
	}
	return AccessToken, err
}

func GenerateAllTokens(id uint, RefreshToken string, LongliveToken string) (AllTokenResponse, error) {

	AccessToken, err := NewAccessToken(id)
	if err != nil {
		return AllTokenResponse{}, err
	}

	if RefreshToken == "" {
		RefreshToken, err = NewRefreshToken()
		if err != nil {
			return AllTokenResponse{}, err
		}
	}

	return AllTokenResponse{
		AccessToken:   AccessToken,
		LongliveToken: LongliveToken,
		RefreshToken:  RefreshToken,
	}, nil
}
