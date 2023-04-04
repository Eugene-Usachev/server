package jwt

import "github.com/dgrijalva/jwt-go"

type AllTokenResponse struct {
	AccessToken   string `json:"access_token"`
	RefreshToken  string `json:"refresh_token"`
	LongliveToken string `json:"longlive_token"`
}

type LongliveAndAccessTokens struct {
	AccessToken   string `json:"access_token"`
	LongliveToken string `json:"longlive_token"`
}

type AccessClaims struct {
	jwt.StandardClaims
	UserID uint `json:"user_id"`
}

type LongliveClaims struct {
	jwt.StandardClaims
	Password string `json:"password"`
}
