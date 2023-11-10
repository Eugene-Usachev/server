package Entities

// Id is a primary key
type Comment struct {
	//Primary key
	Id uint `json:"id"`
	//Index
	ParentPostId uint   `json:"parent_post_id"`
	Data         string `json:"data"`
	Date         int64  `json:"date"`
	ParentUserId uint   `json:"parent_user_id"`
	Likes        uint   `json:"likes"`
	Dislikes     uint   `json:"dislikes"`
	// -1 - disliked, 0 - none, 1 - liked
	LikesStatus     int8     `json:"likes_status"`
	Files           []string `json:"files"`
	ParentCommentId int32    `json:"parent_comment_id"`
}

type CommentDTO struct {
	ParentCommentId uint     `json:"parent_comment_id"`
	ParentPostId    uint     `json:"parent_post_id"`
	Data            string   `json:"data"  binding:"required"`
	Files           []string `json:"files" binding:"required"`
}

type CommentUpdateDTO struct {
	Data  string   `json:"data" binding:"required"`
	Files []string `json:"files" binding:"required"`
}
