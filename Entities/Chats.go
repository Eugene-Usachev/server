package Entities

// ID is a primary key
type Chat struct {
	//Primary key
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	//TODO maybe int32?
	Members []int64 `json:"members"`
}

type ChatDTO struct {
	ID      int64   `json:"id"`
	Name    string  `json:"name" binding:"required"`
	Avatar  string  `json:"avatar" binding:"required"`
	Members []int64 `json:"members" binding:"required"`
}

type ChatUpdateDTO struct {
	Name    string  `json:"name" binding:"required"`
	Avatar  string  `json:"avatar" binding:"required"`
	Members []int64 `json:"members" binding:"required"`
}
