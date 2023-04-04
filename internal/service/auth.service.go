package service

import (
	"GoServer/Entities"
	"GoServer/internal/repository"
	"GoServer/pkg/jwt"
	"context"
	"crypto/sha1"
	"fmt"
	"os"
)

type AuthService struct {
	repository repository.Authorization
}

func NewAuthService(repository repository.Authorization) *AuthService {
	return &AuthService{
		repository: repository,
	}
}

func (services *AuthService) CreateUser(ctx context.Context, user Entities.UserDTO) (uint, error, jwt.LongliveAndAccessTokens) {
	user.Password = generatePasswordHash(user.Password)
	longliveToken, err := jwt.NewLongliveToken(user.Password)
	if err != nil {
		return 0, err, jwt.LongliveAndAccessTokens{}
	}
	id, err := services.repository.CreateUser(ctx, user)
	if err != nil {
		return 0, err, jwt.LongliveAndAccessTokens{}
	}
	accessToken, err := jwt.NewAccessToken(id)
	if err != nil {
		return 0, err, jwt.LongliveAndAccessTokens{}
	}

	return id, nil, jwt.LongliveAndAccessTokens{
		AccessToken:   accessToken,
		LongliveToken: longliveToken,
	}
}

func (services *AuthService) SignIn(ctx context.Context, input Entities.SignInDTO) (uint, string, string, string, string, error) {

	input.Password = generatePasswordHash(input.Password)
	user, err := services.repository.SignInUser(ctx, input)
	if err != nil {
		return 0, "", "", "", "", err
	}
	accessToken, err := jwt.NewAccessToken(user.ID)
	if err != nil {
		return 0, "", "", "", "", err
	}

	longliveToken, err := jwt.NewLongliveToken(input.Password)
	if err != nil {
		return 0, "", "", "", "", err
	}

	return user.ID, user.Email, user.Login, accessToken, longliveToken, nil

}

func (services *AuthService) RefreshToken(ctx context.Context, password, email string) (uint, string, string, error) {

	id, err := services.repository.RefreshTokens(ctx, email, password)
	if err != nil {
		return 0, "", "", err
	}

	accessToken, err := jwt.NewAccessToken(id)
	if err != nil {
		return 0, "", "", err
	}
	longliveToken, err := jwt.NewLongliveToken(password)
	if err != nil {
		return 0, "", "", err
	}

	return id, accessToken, longliveToken, err
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(os.Getenv("SALT"))))
}
