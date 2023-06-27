package repository

import (
	"GoServer/Entities"
	"context"
)

type ChatPostgres struct {
	dataBases *DataBases
}

func NewChatPostgres(dataBases *DataBases) *ChatPostgres {
	return &ChatPostgres{
		dataBases: dataBases,
	}
}

func (repository *ChatPostgres) CreateChat(ctx context.Context, chatDTO Entities.ChatDTO) (uint, error) {
	tx, err := repository.dataBases.Postgres.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var chatId uint
	row := tx.QueryRow(ctx, "INSERT INTO chats (name, avatar, members) VALUES ($1,$2,$3) RETURNING id", chatDTO.Name, chatDTO.Avatar, chatDTO.Members)
	if err = row.Scan(&chatId); err != nil {
		return 0, err
	}

	if _, err = tx.Exec(ctx, `UPDATE users SET all_chats = array_append(all_chats, $1) WHERE id = ANY($2)`, chatId, chatDTO.Members); err != nil {
		return 0, err
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}

	return chatId, nil
}

func (repository *ChatPostgres) UpdateChat(ctx context.Context, userId, chatId uint, chatDTO Entities.ChatUpdateDTO) error {
	var isOk bool
	row := repository.dataBases.Postgres.QueryRow(ctx, `UPDATE chats SET name = $1, avatar = $2, members = $3
             WHERE id = $4 AND $5=ANY(members) RETURNING TRUE`, chatDTO.Name, chatDTO.Avatar, chatDTO.Members, chatId, userId)

	if err := row.Scan(&isOk); err != nil {
		return err
	}

	return nil
}

func (repository *ChatPostgres) DeleteChat(ctx context.Context, userId, chatId uint) ([]uint, error) {
	var members []uint
	row := repository.dataBases.Postgres.QueryRow(ctx, `DELETE FROM chats WHERE id = $1 AND $2 = ANY(members) RETURNING members`, chatId, userId)
	if err := row.Scan(&members); err != nil {
		return []uint{}, err
	}
	return members, nil
}

//TODO time of time we need to check exists chats and remove deleted

func (repository *ChatPostgres) GetChats(ctx context.Context, userId uint) (string, string, string, []uint, []uint, string, []Entities.Chat, error) {
	var (
		err         error
		chatLists   string
		avatar      string
		name        string
		surname     string
		friends     = []uint{}
		subscribers = []uint{}
		allChats    = []uint{}
	)
	err = repository.dataBases.Postgres.QueryRow(ctx, `
		SELECT chat_lists, all_chats, avatar, name, surname, subscribers, friends
		FROM users
		WHERE id = $1
	`, userId).Scan(&chatLists, &allChats, &avatar, &name, &surname, &subscribers, &friends)
	if err != nil {
		return avatar, "", "", nil, nil, "", nil, err
	}

	if len(allChats) == 0 {
		return avatar, name, surname, friends, subscribers, chatLists, []Entities.Chat{}, err
	}

	chats := make([]Entities.Chat, 0, len(allChats))

	rows, err := repository.dataBases.Postgres.Query(ctx, `
		SELECT id, name, avatar, members
		FROM chats
		WHERE id = ANY ($1)
	`, allChats)
	if err != nil {
		return avatar, name, surname, friends, subscribers, chatLists, nil, err
	}

	for rows.Next() {
		var chat Entities.Chat
		err = rows.Scan(&chat.ID, &chat.Name, &chat.Avatar, &chat.Members)
		if err != nil {
			continue
		}
		chats = append(chats, chat)
	}

	return avatar, name, surname, friends, subscribers, chatLists, chats, err
}

func (repository *ChatPostgres) UpdateChatLists(ctx context.Context, id uint, newChatLists string) error {
	_, err := repository.dataBases.Postgres.Exec(ctx, `UPDATE users SET chat_lists=$2 WHERE id=$1`, id, newChatLists)
	return err
}
