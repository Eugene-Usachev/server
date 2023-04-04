package repository

import (
	"GoServer/Entities"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MusicPostgres struct {
	database *pgxpool.Pool
}

func NewMusicPostgres(db *pgxpool.Pool) *MusicPostgres {
	return &MusicPostgres{
		database: db,
	}
}

func (repository *MusicPostgres) GetMusics(ctx context.Context, name string, offset uint) ([]Entities.Music, error) {
	var (
		music  = Entities.Music{}
		musics []Entities.Music
	)

	name = "%" + name + "%"

	rows, err := repository.database.Query(ctx, `SELECT * FROM musics WHERE author LIKE $1 OR title LIKE $1
                     ORDER BY number_of_eavesdroppers DESC LIMIT 20 OFFSET $2`, name, offset)
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
	err := repository.database.QueryRow(ctx, `UPDATE musics SET number_of_eavesdroppers = number_of_eavesdroppers + 1 WHERE id = $1 RETURNING title, parent_user_id`,
		id).Scan(&title, &parentUserID)
	if err != nil {
		return 0, "", err
	}

	return parentUserID, title, nil
}

func (repository *MusicPostgres) AddMusic(ctx context.Context, id uint, music Entities.CreateMusicDTO) (uint, error) {
	var musicId uint
	row := repository.database.QueryRow(ctx, `INSERT INTO musics (author, title, parent_user_id) VALUES ($1, $2, $3) RETURNING id`, music.Author, music.Title, id)

	err := row.Scan(&musicId)
	if err != nil {
		return 0, err
	}

	return musicId, nil
}
