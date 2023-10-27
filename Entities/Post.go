package Entities

// ID is a primary key
type Post struct {
	//Primary key
	ID uint `json:"id"`
	//Index
	ParentUserID uint     `json:"parent_user_id"`
	Likes        uint     `json:"likes"`
	LikedBy      []int32  `json:"liked_by"`
	Dislikes     uint     `json:"dislikes"`
	DislikedBy   []int32  `json:"disliked_by"`
	Data         string   `json:"data"`
	Date         int64    `json:"date"`
	Files        []string `json:"files"`
	HaveSurvey   bool     `json:"have_survey"`
	// TODO
	//IsVisibleForAll bool `json:"is_visible_for_all"`
}

type GetPostDTO struct {
	ID       uint `json:"id"`
	Likes    uint `json:"likes"`
	Dislikes uint `json:"dislikes"`
	// -1 - disliked, 0 - none, 1 - liked
	LikeStatus int8     `json:"likes_status"`
	Data       string   `json:"data"`
	Date       int64    `json:"date"`
	Files      []string `json:"files"`
	HaveSurvey bool     `json:"have_survey"`
	// TODO
	//IsVisibleForAll bool `json:"is_visible_for_all"`
}

type CreatePostDTO struct {
	Data        string   `json:"data" binding:"required"`
	Files       []string `json:"files" binding:"required"`
	HaveASurvey bool     `json:"have_survey" binding:"required"`
}
