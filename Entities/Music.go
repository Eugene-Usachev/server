package Entities

// Music SRC is ParentUserID + "/Music/" +  ID + "." + file extension. e.g. "1/Music/1.mp3". File extension gets from Title field.
type Music struct {
	//Primary key. SRC is ParentUserID + "/Music/" +  ID + "." + file extension. e.g. "1/Music/1.mp3". File extension gets from Title field.
	ID                    int    `json:"id"`
	Author                string `json:"author"`
	ParentUserID          int    `json:"parent_user_id"`
	Title                 string `json:"title"`
	NumberOfEavesdroppers int    `json:"number_of_eavesdroppers"`
}

type CreateMusicDTO struct {
	Title  string `json:"title" form:"title"`
	Author string `json:"author" form:"author"`
}
