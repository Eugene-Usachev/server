package Entities

// ID is a primary key
type Message struct {
	//Primary key
	ID int64 `json:"id"`
	//Index
	ParentChatID    uint     `json:"parent_chat_id"`
	ParentUserID    uint     `json:"parent_user_id"`
	Data            string   `json:"data"`
	Date            string   `json:"date"`
	Files           []string `json:"files"`
	MessageParentID uint     `json:"message_parent_id"`
}

type MessageDTO struct {
	ID              int64    `json:"id"`
	ParentChatID    uint     `json:"parent_chat_id" binding:"required"`
	ParentUserID    uint     `json:"parent_user_id" binding:"required"`
	Data            string   `json:"data" binding:"required"`
	Files           []string `json:"files" binding:"required"`
	MessageParentID uint     `json:"message_parent_id"`
	Date            string   `json:"date"`
}
