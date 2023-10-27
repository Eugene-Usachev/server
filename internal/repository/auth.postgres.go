package repository

import (
	"GoServer/Entities"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
)

type AuthPostgres struct {
	dataBases *DataBases
}

func NewAuthPostgres(dataBases *DataBases) *AuthPostgres {
	return &AuthPostgres{
		dataBases: dataBases,
	}
}

var EmailBusy = errors.New("email is busy")
var LoginBusy = errors.New("login is busy")

func (repository *AuthPostgres) CreateUser(ctx context.Context, user Entities.UserDTO) (uint, error) {
	var err error

	if err = repository.dataBases.Postgres.pool.QueryRow(ctx, "SELECT id FROM users WHERE login=$1", user.Login).Scan(); err != pgx.ErrNoRows {
		return 0, LoginBusy
	}

	if err = repository.dataBases.Postgres.pool.QueryRow(ctx, "SELECT id FROM users WHERE email=$1", user.Email).Scan(); err != pgx.ErrNoRows {
		return 0, EmailBusy
	}

	var id uint
	row := repository.dataBases.Postgres.pool.QueryRow(ctx, `INSERT INTO users (login, email, password, name, surname) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		user.Login, user.Email, user.Password, user.Name, user.Surname)
	if err = row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (repository *AuthPostgres) SignInUser(ctx context.Context, input Entities.SignInDTO) (Entities.SignInReturnDTO, error) {
	var (
		user Entities.SignInReturnDTO
		row  pgx.Row
		err  error
	)

	if input.Login == "" {
		row = repository.dataBases.Postgres.pool.QueryRow(ctx, `SELECT email, login, id, name, surname, avatar FROM users WHERE email = $1 AND password =$2`,
			input.Email, input.Password)
	} else {
		row = repository.dataBases.Postgres.pool.QueryRow(ctx, `SELECT email, login, id, name, surname, avatar FROM users WHERE login = $1 AND password =$2`,
			input.Login, input.Password)
	}
	if err = row.Scan(&user.Email, &user.Login, &user.ID, &user.Name, &user.Surname, &user.Avatar); err != nil {
		return Entities.SignInReturnDTO{}, err
	}
	return user, nil
}

func (repository *AuthPostgres) RefreshTokens(ctx context.Context, email, password string) (uint, error) {
	var id uint

	row := repository.dataBases.Postgres.pool.QueryRow(ctx, `SELECT id FROM users WHERE password = $1 AND email = $2`, password, email)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (repository *AuthPostgres) Refresh(ctx context.Context, id uint, password string) (dto Entities.RefreshResponseDTO, err error) {
	err = repository.dataBases.Postgres.pool.QueryRow(ctx, `SELECT avatar, name, surname FROM users WHERE id = $1 AND password = $2`, id, password).Scan(&dto.Avatar, &dto.Name, &dto.Surname)
	return
}

func (repository *AuthPostgres) CheckPassword(ctx context.Context, id uint, password string) error {
	return repository.dataBases.Postgres.pool.QueryRow(ctx, `SELECT 1 FROM users WHERE id = $1 AND password = $2`, id, password).Scan(nil)
}
