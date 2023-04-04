package Entities

import (
	"database/sql"
)

// ID is a primary key
type Comment struct {
	//Primary key
	ID int `json:"id"`
	//Index
	ParentPostID    int           `json:"parent_post_id"`
	Data            string        `json:"data"`
	Date            string        `json:"date"`
	ParentUserId    int           `json:"parent_user_id"`
	Likes           int           `json:"likes"`
	LikedBy         []int32       `json:"liked_by"`
	Dislikes        int           `json:"dislikes"`
	DislikedBy      []int32       `json:"disliked_by"`
	Files           []string      `json:"files"`
	ParentCommentId sql.NullInt32 `json:"parent_comment_id"`
}

type CommentDTO struct {
	ParentCommentId uint     `json:"parent_comment_id"`
	ParentPostID    uint     `json:"parent_post_id"`
	Data            string   `json:"data"  binding:"required"`
	Files           []string `json:"files" binding:"required"`
}

type CommentUpdateDTO struct {
	Data  string   `json:"data" binding:"required"`
	Files []string `json:"files" binding:"required"`
}
