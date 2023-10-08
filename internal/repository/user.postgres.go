package repository

import (
	"GoServer/Entities"
	"context"
	"errors"
	"fmt"
)

type UserPostgres struct {
	dataBases *DataBases
}

func NewUserPostgres(dataBases *DataBases) *UserPostgres {
	return &UserPostgres{
		dataBases: dataBases,
	}
}

func (repository *UserPostgres) GetUserById(ctx context.Context, id uint) (Entities.GetUserDTO, error) {
	var (
		user Entities.GetUserDTO
	)
	row := repository.dataBases.Postgres.QueryRow(ctx, `
		SELECT name, surname, avatar, birthday, attitude_to_alcohol, attitude_to_smocking, attitude_to_sport,
			   family_status, friends, users.subscribers,
				favourites_books, favourites_films, favourites_games, favourites_meals, description, dreams,
				place_of_residence  FROM users WHERE id = $1
		`, id)
	if err := row.Scan(&user.Name, &user.Surname, &user.Avatar, &user.Birthday, &user.Attitude_to_alcohol, &user.Attitude_to_smocking,
		&user.Attitude_to_sport, &user.Family_status, &user.Friends, &user.Subscribers, &user.Favourites_books, &user.Favourites_films,
		&user.Favourites_games, &user.Favourites_meals, &user.Description, &user.Dreams, &user.Place_of_residence); err != nil {
		return user, err
	}
	return user, nil
}

func (repository *UserPostgres) GetUserSubsIds(ctx context.Context, id uint) ([]uint, error) {
	var ids = []uint{}

	row := repository.dataBases.Postgres.QueryRow(ctx, `SELECT subscribers FROM users WHERE id = $1`, id)
	if err := row.Scan(&ids); err != nil {
		return ids, err
	}

	return ids, nil
}

func (repository *UserPostgres) GetFriendsAndSubs(ctx context.Context, clientId, userId uint) (Entities.GetFriendsAndSubsDTO, error) {
	var DTO Entities.GetFriendsAndSubsDTO
	row := repository.dataBases.Postgres.QueryRow(ctx, `SELECT name, surname, avatar, friends, subscribers FROM users WHERE id = $1`, userId)
	if err := row.Scan(&DTO.User.Name, &DTO.User.Surname, &DTO.User.Avatar, &DTO.User.Friends, &DTO.User.Subscribers); err != nil {
		return Entities.GetFriendsAndSubsDTO{}, err
	}

	if clientId == 0 {
		return DTO, nil
	}
	row = repository.dataBases.Postgres.QueryRow(ctx, `SELECT name, surname, avatar, friends, subscribers FROM users WHERE id = $1`, clientId)
	if err := row.Scan(&DTO.Client.Name, &DTO.Client.Surname, &DTO.Client.Avatar, &DTO.Client.Friends, &DTO.Client.Subscribers); err != nil {
		return Entities.GetFriendsAndSubsDTO{}, err
	}

	return DTO, nil
}

func (repository *UserPostgres) GetUsersForFriendsPage(ctx context.Context, idOfUsers string) ([]Entities.FriendUser, error) {
	var miniUsers = []Entities.FriendUser{}
	str := fmt.Sprintf(`SELECT id, name, surname, avatar, subscribers FROM users WHERE id in %s`, idOfUsers)
	rows, err := repository.dataBases.Postgres.Query(ctx, str)
	for rows.Next() {
		var miniUser Entities.FriendUser
		if err = rows.Scan(&miniUser.ID, &miniUser.Name, &miniUser.Surname, &miniUser.Avatar, &miniUser.Subscribers); err == nil {
			miniUsers = append(miniUsers, miniUser)
		} else {
			continue
		}
	}
	if err != nil {
		return nil, err
	}
	return miniUsers, nil
}

func (repository *UserPostgres) GetUsers(ctx context.Context, idOfUsers string) ([]Entities.MiniUser, error) {
	var miniUsers []Entities.MiniUser = []Entities.MiniUser{}
	str := fmt.Sprintf(`SELECT id, name, surname, avatar FROM users WHERE id in %s`, idOfUsers)
	rows, err := repository.dataBases.Postgres.Query(ctx, str)
	for rows.Next() {
		var miniUser Entities.MiniUser
		if err = rows.Scan(&miniUser.ID, &miniUser.Name, &miniUser.Surname, &miniUser.Avatar); err == nil {
			miniUsers = append(miniUsers, miniUser)
		} else {
			continue
		}
	}
	if err != nil {
		return nil, err
	}
	return miniUsers, nil
}

func (repository *UserPostgres) UpdateUser(ctx context.Context, id uint, UpdateUserDTO Entities.UpdateUserDTO) error {
	var err error
	_, err = repository.dataBases.Postgres.Exec(ctx, `UPDATE users SET favourites_films=$2, favourites_books=$3,
		 favourites_games=$4, dreams = $5,attitude_to_sport =$6, attitude_to_alcohol =$7, attitude_to_smocking =$8 ,
		 place_of_residence =$9, family_status =$10,name =$11, surname=$12, birthday=$13, favourites_meals=$14, description=$15 WHERE id = $1`,
		id, UpdateUserDTO.Favourites_films, UpdateUserDTO.Favourites_books, UpdateUserDTO.Favourites_games,
		UpdateUserDTO.Dreams, UpdateUserDTO.Attitude_to_sport, UpdateUserDTO.Attitude_to_alcohol, UpdateUserDTO.Attitude_to_smocking,
		UpdateUserDTO.Place_of_residence, UpdateUserDTO.Family_status, UpdateUserDTO.Name, UpdateUserDTO.Surname, UpdateUserDTO.Birthday,
		UpdateUserDTO.Favourites_meals, UpdateUserDTO.Description)
	return err
}

func (repository *UserPostgres) ChangeAvatar(ctx context.Context, id uint, fileName string) error {
	var err error
	_, err = repository.dataBases.Postgres.Exec(ctx, `UPDATE users SET avatar=$1 WHERE id = $2`, fileName, id)
	if err != nil {
		return err
	}
	return nil
}

func (repository *UserPostgres) AddToFriends(ctx context.Context, id, body uint) error {
	tx, err := repository.dataBases.Postgres.Begin(ctx)
	if err != nil {
		return err
	}

	row := tx.QueryRow(ctx, `WITH updated AS (
		UPDATE users 
		SET friends = array_append(friends, $2), subscribers = array_remove(subscribers, $2)
		WHERE id = $1 AND $2 = ANY(subscribers) AND array_position(friends , $2) IS NULL 
		RETURNING TRUE
	)
	UPDATE users 
	SET friends = array_append(friends, $1) 
	WHERE id = $2
	AND (SELECT true FROM updated) RETURNING TRUE`, id, body)

	var result bool
	if err = row.Scan(&result); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	if !result {
		_ = tx.Rollback(ctx)
		return errors.New("not subscribed")
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (repository *UserPostgres) DeleteFromFriends(ctx context.Context, id, body uint) error {
	tx, err := repository.dataBases.Postgres.Begin(ctx)
	if err != nil {
		return err
	}

	row := tx.QueryRow(ctx, `WITH updated AS (
		UPDATE users 
		SET friends = array_remove(friends, $2), subscribers = array_append(subscribers, $2)
		WHERE id = $1 AND $2 = ANY(friends)
		RETURNING TRUE
	)
	UPDATE users 
	SET friends = array_remove(friends, $1) 
	WHERE id = $2
	AND (SELECT true FROM updated) RETURNING TRUE`, id, body)

	var result bool
	if err = row.Scan(&result); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	if !result {
		_ = tx.Rollback(ctx)
		return errors.New("not subscribed")
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (repository *UserPostgres) AddToSubs(ctx context.Context, id, body uint) error {
	var result bool

	row := repository.dataBases.Postgres.QueryRow(ctx, `
		UPDATE users SET subscribers = array_append(subscribers, $1) WHERE id = $2 AND array_position(subscribers , $1) IS NULL RETURNING TRUE;
	`, id, body)
	if err := row.Scan(&result); err != nil {
		return err
	}
	if !result {
		return errors.New("are subscribed")
	}

	return nil
}

func (repository *UserPostgres) DeleteFromSubs(ctx context.Context, id, body uint) error {

	_, err := repository.dataBases.Postgres.Exec(ctx, `
		UPDATE users SET subscribers = array_remove(subscribers, $1) WHERE id = $2;
	`, id, body)
	if err != nil {
		return err
	}

	return nil
}

func (repository *UserPostgres) DeleteUser(ctx context.Context, id uint) error {
	_, err := repository.dataBases.Postgres.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}
