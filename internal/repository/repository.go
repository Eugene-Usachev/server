package repository

import (
	"GoServer/Entities"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Authorization interface {
	CreateUser(ctx context.Context, user Entities.UserDTO) (uint, error)
	SignInUser(ctx context.Context, input Entities.SignInDTO) (Entities.SignInReturnDTO, error)
	RefreshTokens(ctx context.Context, email, password string) (uint, error)
}

type User interface {
	GetUserById(ctx context.Context, id uint, requestOwnerId uint) (Entities.GetUserDTO, []int64, error)
	GetFriendsAndSubs(ctx context.Context, clientId, userId uint) (Entities.GetFriendsAndSubsDTO, error)
	UpdateUser(ctx context.Context, id uint, UpdateUserDTO Entities.UpdateUserDTO) error
	ChangeAvatar(ctx context.Context, id uint, fileName string) error
	AddToFriends(ctx context.Context, id, body uint) error
	DeleteFromFriends(ctx context.Context, id, body uint) error
	AddToSubs(ctx context.Context, id, body uint) error
	DeleteFromSubs(ctx context.Context, id, body uint) error
	DeleteUser(ctx context.Context, id uint) error
	GetUsersForFriendsPage(ctx context.Context, idOfUsers string) ([]Entities.FriendUser, error)
	GetUsers(ctx context.Context, idOfUsers string) ([]Entities.MiniUser, error)
}

type Post interface {
	/*region post*/
	CreateAPost(ctx context.Context, id uint, postDTO Entities.CreateAPostDTO, surveyDTO Entities.CreateASurveyDTO, date string) error
	GetPostsByUserID(ctx context.Context, userID uint, offset uint) ([]Entities.Post, []Entities.Survey, error)
	LikePost(ctx context.Context, userId, postId uint) error
	UnlikePost(ctx context.Context, userId, postId uint) error
	DislikePost(ctx context.Context, userId, postId uint) error
	UndislikePost(ctx context.Context, userId, postId uint) error
	DeletePost(ctx context.Context, postId, userId uint) error
	/*endregion*/

	/*region comments*/
	GetCommentsByPostId(ctx context.Context, postId uint, offset uint) ([]Entities.Comment, error)
	CreateComment(ctx context.Context, userId uint, comment Entities.CommentDTO, date string) (uint, error)
	LikeComment(ctx context.Context, userID uint, commentID uint) error
	UnlikeComment(ctx context.Context, userID uint, commentID uint) error
	DislikeComment(ctx context.Context, userID uint, commentID uint) error
	UndislikeComment(ctx context.Context, userID uint, commentID uint) error
	UpdateComment(ctx context.Context, userID uint, commentID uint, updateDTO Entities.CommentUpdateDTO) error
	DeleteComment(ctx context.Context, userID uint, commentID uint) error
	/*endregion*/

	/*region survey*/
	VoteInSurvey(ctx context.Context, userId uint, surveyId uint, votedFor []uint8) error
	/*endregion*/
}

type Message interface {
	SaveMessage(ctx context.Context, userId int64, messageDTO Entities.MessageDTO) (int64, []int64, error)
	UpdateMessage(ctx context.Context, messageId, userId int64, newData string) ([]int64, error)
	DeleteMessage(ctx context.Context, messageId, userId int64) ([]int64, error)
	GetLastMessages(ctx context.Context, userId uint, chatsId string) ([]Entities.Message, error)
	GetMessages(ctx context.Context, chatId, offset uint) ([20]Entities.Message, error)
}

type Chat interface {
	CreateChat(ctx context.Context, chatDTO Entities.ChatDTO) (int64, error)
	UpdateChat(ctx context.Context, userId, chatId int64, chatDTO Entities.ChatUpdateDTO) error
	DeleteChat(ctx context.Context, userId, chatId int64) ([]int64, error)
	GetChats(ctx context.Context, userId int64) (string, string, string, []int64, []int64, string, []Entities.Chat, error)
	UpdateChatLists(ctx context.Context, id int64, newChatLists string) error
}

type Music interface {
	GetMusics(ctx context.Context, name string, offset uint) ([]Entities.Music, error)
	GetMusic(ctx context.Context, id uint) (uint, string, error)
	AddMusic(ctx context.Context, id uint, music Entities.CreateMusicDTO) (uint, error)
}

type Repository struct {
	Authorization
	User
	Post
	Message
	Music
	Chat
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(pool),
		User:          NewUserPostgres(pool),
		Music:         NewMusicPostgres(pool),
		Post:          NewPostPostgres(pool),
		Chat:          NewChatPostgres(pool),
		Message:       NewMessagePostgres(pool),
	}
}
