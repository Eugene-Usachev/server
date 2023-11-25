package repository

import (
	"GoServer/Entities"
	"context"
)

type MusicPostgres struct {
	dataBases *DataBases
}

func NewMusicPostgres(dataBases *DataBases) *MusicPostgres {
	return &MusicPostgres{
		dataBases: dataBases,
	}
}

func (repository *MusicPostgres) GetMusics(ctx context.Context, name string, offset uint) ([]Entities.Music, error) {
	var (
		music  = Entities.Music{}
		musics []Entities.Music
	)

	name = "%" + name + "%"

	rows, err := repository.dataBases.Postgres.pool.Query(ctx, `SELECT * FROM musics WHERE author LIKE $1 OR title LIKE $1
                     ORDER BY number_of_eavesdroppers DESC LIMIT 20 OFFSET $2`, name, offset)
	defer rows.Close()
	if err != nil {
		return musics, err
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		_ = rows.Scan(&music.ID, &music.ParentUserID, &music.Author, &music.Title, &music.NumberOfEavesdroppers)
		musics = append(musics, music)
	}

	return musics, nil
}

func (repository *MusicPostgres) GetMusic(ctx context.Context, id uint) (uint, string, error) {
	var (
		title        string
		parentUserID uint
	)
	err := repository.dataBases.Postgres.pool.QueryRow(ctx, `UPDATE musics SET number_of_eavesdroppers = number_of_eavesdroppers + 1 WHERE id = $1 RETURNING title, parent_user_id`,
		id).Scan(&title, &parentUserID)
	if err != nil {
		return 0, "", err
	}

	return parentUserID, title, nil
}

func (repository *MusicPostgres) AddMusic(ctx context.Context, id uint, music Entities.CreateMusicDTO) (uint, error) {
	var musicId uint
	row := repository.dataBases.Postgres.pool.QueryRow(ctx, `INSERT INTO musics (author, title, parent_user_id) VALUES ($1, $2, $3) RETURNING id`, music.Author, music.Title, id)

	err := row.Scan(&musicId)
	if err != nil {
		return 0, err
	}

	return musicId, nil
}
