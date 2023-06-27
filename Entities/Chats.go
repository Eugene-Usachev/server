package Entities

// ID is a primary key
type Chat struct {
	//Primary key
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Avatar  string `json:"avatar"`
	Members []uint `json:"members"`
}

type ChatDTO struct {
	ID      uint   `json:"id"`
	Name    string `json:"name" binding:"required"`
	Avatar  string `json:"avatar" binding:"required"`
	Members []uint `json:"members" binding:"required"`
}

type ChatUpdateDTO struct {
	Name    string `json:"name" binding:"required"`
	Avatar  string `json:"avatar" binding:"required"`
	Members []uint `json:"members" binding:"required"`
}
