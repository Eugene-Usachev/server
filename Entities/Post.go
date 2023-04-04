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
	Date         string   `json:"date"`
	Files        []string `json:"files"`
	HaveASurvey  bool     `json:"have_survey"`
}

type CreateAPostDTO struct {
	Data        string   `json:"data" binding:"required"`
	Files       []string `json:"files" binding:"required"`
	HaveASurvey bool     `json:"have_survey" binding:"required"`
}
