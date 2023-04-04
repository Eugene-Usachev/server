package repository

import (
	"GoServer/Entities"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthPostgres struct {
	database *pgxpool.Pool
}

func NewAuthPostgres(db *pgxpool.Pool) *AuthPostgres {
	return &AuthPostgres{
		database: db,
	}
}

func (repository *AuthPostgres) CreateUser(ctx context.Context, user Entities.UserDTO) (uint, error) {
	var id uint
	row := repository.database.QueryRow(ctx, `INSERT INTO users (login, email, password, name, surname) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		user.Login, user.Email, user.Password, user.Name, user.Surname)
	if err := row.Scan(&id); err != nil {
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
		row = repository.database.QueryRow(ctx, `SELECT email, login, id FROM users WHERE email = $1 AND password =$2`,
			input.Email, input.Password)
	} else {
		row = repository.database.QueryRow(ctx, `SELECT email, login, id FROM users WHERE login = $1 AND password =$2`,
			input.Login, input.Password)
	}
	if err = row.Scan(&user.Email, &user.Login, &user.ID); err != nil {
		return Entities.SignInReturnDTO{}, err
	}
	return user, nil
}

func (repository *AuthPostgres) RefreshTokens(ctx context.Context, email, password string) (uint, error) {
	var id uint

	row := repository.database.QueryRow(ctx, `SELECT id FROM users WHERE password = $1 AND email = $2`, password, email)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
