package service

import (
	"GoServer/Entities"
	"GoServer/internal/repository"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/Eugene-Usachev/fastbytes"
	"github.com/Eugene-Usachev/fst"
	"hash"
	"os"
	"sync"
)

type AuthService struct {
	repository       repository.Authorization
	accessConverter  *fst.Converter
	refreshConverter *fst.Converter
}

type AuthServiceConfig struct {
	repository       repository.Authorization
	accessConverter  *fst.Converter
	refreshConverter *fst.Converter
}

func NewAuthService(cfg *AuthServiceConfig) *AuthService {
	return &AuthService{
		repository:       cfg.repository,
		accessConverter:  cfg.accessConverter,
		refreshConverter: cfg.refreshConverter,
	}
}

func (services *AuthService) CreateUser(ctx context.Context, user Entities.UserDTO) (uint, error, Entities.AllTokenResponse) {
	user.Password = generatePasswordHash(user.Password)
	id, err := services.repository.CreateUser(ctx, user)
	if err != nil {
		return 0, err, Entities.AllTokenResponse{}
	}
	var tokens Entities.AllTokenResponse
	tokens.AccessToken = services.accessConverter.NewToken(fastbytes.U2B(id))
	tokens.RefreshToken = services.refreshConverter.NewToken(fastbytes.S2B(user.Password))

	return id, nil, tokens
}

func (services *AuthService) SignIn(ctx context.Context, input Entities.SignInDTO) (Entities.SignInReturnDTO, Entities.AllTokenResponse, error) {

	input.Password = generatePasswordHash(input.Password)
	user, err := services.repository.SignInUser(ctx, input)
	if err != nil {
		return Entities.SignInReturnDTO{}, Entities.AllTokenResponse{}, err
	}
	var tokens Entities.AllTokenResponse
	tokens.AccessToken = services.accessConverter.NewToken(fastbytes.U2B(user.ID))
	tokens.RefreshToken = services.refreshConverter.NewToken(fastbytes.S2B(input.Password))
	if err != nil {
		return Entities.SignInReturnDTO{}, Entities.AllTokenResponse{}, err
	}

	return user, tokens, nil

}

func (services *AuthService) Check(ctx context.Context, email, login string) (isEmailNotBusy, isLoginNotBusy bool) {
	return services.repository.Check(ctx, email, login)
}

var (
	invalidTokenError = errors.New("invalid token")
)

func (services *AuthService) Refresh(ctx context.Context, id uint, refreshToken string) (Entities.RefreshResponseDTO, error) {
	passwordHash, err := services.refreshConverter.ParseToken(refreshToken)
	if err != nil {
		return Entities.RefreshResponseDTO{}, err
	}
	if len(passwordHash) == 0 {
		return Entities.RefreshResponseDTO{}, invalidTokenError
	}
	var dto Entities.RefreshResponseDTO
	dto, err = services.repository.Refresh(ctx, id, fastbytes.B2S(passwordHash))
	if err != nil {
		return Entities.RefreshResponseDTO{}, err
	}
	dto.AccessToken = services.accessConverter.NewToken(fastbytes.U2B(id))
	dto.RefreshToken = services.refreshConverter.NewToken(passwordHash)
	return dto, nil
}

func (services *AuthService) RefreshTokens(ctx context.Context, id uint, refreshToken string) (Entities.AllTokenResponse, error) {
	passwordHash, err := services.refreshConverter.ParseToken(refreshToken)
	if err != nil {
		return Entities.AllTokenResponse{}, err
	}
	if len(passwordHash) == 0 {
		return Entities.AllTokenResponse{}, invalidTokenError
	}
	err = services.repository.CheckPassword(ctx, id, fastbytes.B2S(passwordHash))
	if err != nil {
		return Entities.AllTokenResponse{}, err
	}
	var dto Entities.AllTokenResponse
	dto.AccessToken = services.accessConverter.NewToken(fastbytes.U2B(id))
	dto.RefreshToken = services.refreshConverter.NewToken(passwordHash)
	return dto, nil
}

var (
	SALT     = fastbytes.S2B(os.Getenv("SALT"))
	hashPool = sync.Pool{
		New: func() interface{} {
			return sha256.New()
		},
	}
)

func generatePasswordHash(password string) string {
	sha := hashPool.Get().(hash.Hash)
	sha.Reset()
	defer func() {
		hashPool.Put(sha)
	}()
	sha.Write(fastbytes.S2B(password))
	return fmt.Sprintf("%x", sha.Sum(nil))
}
