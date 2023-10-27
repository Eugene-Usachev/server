package repository

import (
	"GoServer/Entities"
	"context"
	"github.com/jackc/pgx/v5"
)

type MessagePostgres struct {
	dataBases *DataBases
}

func NewMessagePostgres(dataBases *DataBases) *MessagePostgres {
	return &MessagePostgres{
		dataBases: dataBases,
	}
}

func (repository *MessagePostgres) SaveMessage(ctx context.Context, userId uint, messageDTO Entities.MessageDTO) (uint, []uint, error) {
	tx, err := repository.dataBases.Postgres.pool.Begin(ctx)
	if err != nil {
		return 0, nil, err
	}
	defer tx.Rollback(ctx)

	row := tx.QueryRow(ctx, `
		WITH chat_data AS (
		  SELECT members
		  FROM chats
		  WHERE id = $1 AND $2 = ANY(members)
		),
		new_message AS (
		  INSERT INTO messages (parent_chat_id, parent_user_id, data, date, files, message_parent_id)
		  SELECT $1, $2, $3, $6, $4, $5
		  FROM chat_data
		  RETURNING id
		)
		SELECT new_message.id, chat_data.members
		FROM chat_data, new_message
	`,
		messageDTO.ParentChatID, userId, messageDTO.Data, messageDTO.Files, messageDTO.MessageParentID, messageDTO.Date)

	var (
		chatMembers []uint
		messageId   uint
	)
	if err = row.Scan(&messageId, &chatMembers); err != nil {
		return 0, chatMembers, err
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, chatMembers, err
	}

	return messageId, chatMembers, nil
}

func (repository *MessagePostgres) UpdateMessage(ctx context.Context, messageId, userId uint, newData int64) ([]uint, error) {
	var members []uint
	row := repository.dataBases.Postgres.pool.QueryRow(ctx, `
		WITH new_message AS (
			UPDATE messages SET data = $3
				WHERE id = $2 AND parent_user_id = $1 RETURNING parent_chat_id
		)
			SELECT members FROM chats, new_message WHERE id=new_message.parent_chat_id
	`, userId, messageId, newData)

	if err := row.Scan(&members); err != nil {
		return members, err
	}

	return members, nil
}

func (repository *MessagePostgres) DeleteMessage(ctx context.Context, messageId, userId uint) ([]uint, error) {
	row := repository.dataBases.Postgres.pool.QueryRow(ctx, `
		WITH chat_id AS (
		    DELETE FROM messages WHERE id = $1 AND parent_user_id = $2 RETURNING parent_chat_id
		)
		SELECT members FROM chats, chat_id WHERE id=chat_id.parent_chat_id
	`, messageId, userId)
	var members []uint
	if err := row.Scan(&members); err != nil {
		return members, err
	}
	return members, nil
}

func (repository *MessagePostgres) GetLastMessages(ctx context.Context, userId uint, chatsId string) ([]Entities.Message, error) {
	var array = []Entities.Message{}

	rows, err := repository.dataBases.Postgres.pool.Query(ctx, `
		SELECT id FROM chats WHERE $1 = ANY(members) AND id IN `+chatsId+`
	`, userId)
	if err != nil {
		return array, err
	}
	defer rows.Close()

	var chatsWhereUserIsMember []int64
	for rows.Next() {
		var chatId int64
		if err := rows.Scan(&chatId); err != nil {
			continue
		}
		chatsWhereUserIsMember = append(chatsWhereUserIsMember, chatId)
	}

	rows, err = repository.dataBases.Postgres.pool.Query(ctx, `
		SELECT t1.* FROM messages t1
		INNER JOIN (
			SELECT parent_chat_id, max(id) as lastID 
			FROM messages
			WHERE parent_chat_id = ANY($1)
			GROUP BY parent_chat_id
		) t2 ON t1.parent_chat_id = t2.parent_chat_id AND t1.id = t2.lastID
	`, chatsWhereUserIsMember)
	if err != nil {
		return array, err
	}
	defer rows.Close()

	for rows.Next() {
		var message Entities.Message
		if err := rows.Scan(&message.ID, &message.ParentChatID, &message.ParentUserID, &message.Data, &message.Date, &message.Files, &message.MessageParentID); err != nil {
			continue
		}
		array = append(array, message)
	}

	return array, nil
}

func (repository *MessagePostgres) GetMessages(ctx context.Context, chatId, offset uint) ([20]Entities.Message, error) {
	var (
		err   error
		array = [20]Entities.Message{}
		rows  pgx.Rows
		i     uint8 = 0
	)

	rows, err = repository.dataBases.Postgres.pool.Query(ctx, `
		SELECT id, parent_chat_id, parent_user_id, files, data, date, message_parent_id 
			FROM messages
			WHERE parent_chat_id = $2
			ORDER BY id DESC
			OFFSET $1
			LIMIT 20
	`, offset, chatId)
	if err != nil {
		return array, err
	}
	defer rows.Close()
	for rows.Next() {
		var message Entities.Message
		if err = rows.Scan(&message.ID, &message.ParentChatID, &message.ParentUserID, &message.Files, &message.Data, &message.Date, &message.MessageParentID); err != nil {
			continue
		}
		array[i] = message
		i++
	}

	return array, err
}
