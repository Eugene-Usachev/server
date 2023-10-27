package Entities

/*
ParentPostID is a primary key
SL is a "survey line", V - "voices"
*/
type Survey struct {
	//Primary key
	ParentPostID  int      `json:"parent_post_id"`
	Data          []string `json:"data"`
	SL0V          int      `json:"sl0v"`
	SL0VBY        []int32  `json:"sl0vby"`
	SL1V          int      `json:"sl1v"`
	SL1VBY        []int32  `json:"sl1vby"`
	SL2V          int      `json:"sl2v"`
	SL2VBY        []int32  `json:"sl2vby"`
	SL3V          int      `json:"sl3v"`
	SL3VBY        []int32  `json:"sl3vby"`
	SL4V          int      `json:"sl4v"`
	SL4VBY        []int32  `json:"sl4vby"`
	SL5V          int      `json:"sl5v"`
	SL5VBY        []int32  `json:"sl5vby"`
	SL6V          int      `json:"sl6v"`
	SL6VBY        []int32  `json:"sl6vby"`
	SL7V          int      `json:"sl7v"`
	SL7VBY        []int32  `json:"sl7vby"`
	SL8V          int      `json:"sl8v"`
	SL8VBY        []int32  `json:"sl8vby"`
	SL9V          int      `json:"sl9v"`
	SL9VBY        []int32  `json:"sl9vby"`
	VotedBy       []int32  `json:"voted_by"`
	Background    uint8    `json:"background"`
	IsMultiVoices bool     `json:"is_multi_voices"`
}

type GetSurveyDTO struct {
	//Primary key
	ParentPostID uint     `json:"parent_post_id"`
	Data         []string `json:"data"`
	SL0V         uint     `json:"sl0v"`
	SL1V         uint     `json:"sl1v"`
	SL2V         uint     `json:"sl2v"`
	SL3V         uint     `json:"sl3v"`
	SL4V         uint     `json:"sl4v"`
	SL5V         uint     `json:"sl5v"`
	SL6V         uint     `json:"sl6v"`
	SL7V         uint     `json:"sl7v"`
	SL8V         uint     `json:"sl8v"`
	SL9V         uint     `json:"sl9v"`
	// Bits flags. Like 1000000000000000 means that user has voted for first line.
	VotedFor      uint16 `json:"voted_for"`
	Background    uint8  `json:"background"`
	IsMultiVoices bool   `json:"is_multi_voices"`
}

type CreateSurveyDTO struct {
	Data          []string `json:"data"`
	Background    uint8    `json:"background"`
	IsMultiVoices bool     `json:"is_multi_voices"`
}
