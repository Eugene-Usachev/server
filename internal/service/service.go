package service

import (
	"GoServer/Entities"
	"GoServer/internal/repository"
	"GoServer/pkg/jwt"
	"context"
	"github.com/gin-gonic/gin"
	"mime/multipart"
)

type Authorization interface {
	CreateUser(ctx context.Context, dto Entities.UserDTO) (uint, error, jwt.LongliveAndAccessTokens)
	SignIn(ctx context.Context, input Entities.SignInDTO) (uint, string, string, string, string, error)
	RefreshToken(ctx context.Context, longLiveToken, email string) (uint, string, string, error)
}

type User interface {
	GetUserById(ctx context.Context, id uint, requestOwnerId uint) (Entities.GetUserDTO, []int64, error)
	GetFriendsAndSubs(ctx context.Context, clientId, userId uint) (Entities.GetFriendsAndSubsDTO, error)
	UpdateUser(ctx context.Context, id uint, UpdateUserDTO Entities.UpdateUserDTO) error
	ChangeAvatar(ctx *gin.Context, id uint) (string, error)
	AddToFriends(ctx context.Context, id, body uint) error
	DeleteFromFriends(ctx context.Context, id, body uint) error
	AddToSubs(ctx context.Context, id, body uint) error
	DeleteFromSubs(ctx context.Context, id, body uint) error
	DeleteUser(ctx context.Context, id uint) error
	GetUsers(ctx context.Context, idOfUsers string) ([]Entities.MiniUser, error)
	GetUsersForFriendsPage(ctx context.Context, idOfUsers string) ([]Entities.FriendUser, error)
}

type Post interface {

	/*region post*/

	CreateAPost(ctx *gin.Context, id uint, postDTO Entities.CreateAPostDTO, surveyDTO Entities.CreateASurveyDTO, files []*multipart.FileHeader) error
	GetPostsByUserID(ctx context.Context, userID uint, offset uint) ([]Entities.Post, []Entities.Survey, error)
	LikePost(ctx context.Context, userId, postId uint) error
	UnlikePost(ctx context.Context, userId, postId uint) error
	DislikePost(ctx context.Context, userId, postId uint) error
	UndislikePost(ctx context.Context, userId, postId uint) error
	DeletePost(ctx context.Context, postId, userId uint) error

	/*endregion*/

	/*region comment*/

	GetCommentsByPostId(ctx context.Context, postId uint, offset uint) ([]Entities.Comment, error)
	CreateComment(ctx context.Context, userId uint, postId uint, comment Entities.CommentDTO) (uint, error)
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
	SaveMessage(ctx context.Context, userId int64, messageDTO Entities.MessageDTO) (int64, []int64, string, error)
	UpdateMessage(ctx context.Context, messageId int64, userId int64, newData string) ([]int64, error)
	DeleteMessage(ctx context.Context, messageId int64, userId int64) ([]int64, error)
	GetLastMessages(ctx context.Context, userId uint, chatsId string) ([]Entities.Message, error)
	GetMessages(ctx context.Context, chatId, offset uint) ([20]Entities.Message, error)
}

type Chat interface {
	CreateChat(ctx context.Context, userId int64, chatDTO Entities.ChatDTO) (int64, error)
	UpdateChat(ctx context.Context, userId, chatId int64, chatDTO Entities.ChatUpdateDTO) error
	DeleteChat(ctx context.Context, userId int64, chatId int64) ([]int64, error)
	GetChats(ctx context.Context, userId int64) (string, string, string, []int64, []int64, string, []Entities.Chat, error)
	UpdateChatLists(ctx context.Context, id int64, newChatLists string) error
}

type Music interface {
	GetMusics(ctx context.Context, name string, offset uint) ([]Entities.Music, error)
	GetMusic(ctx context.Context, id uint) (string, string, error)
	AddMusic(ctx *gin.Context, id uint, music Entities.CreateMusicDTO) error
}

type Service struct {
	Authorization
	User
	Post
	Message
	Music
	Chat
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repository.Authorization),
		User:          NewUserService(repository.User),
		Music:         NewMusicService(repository.Music),
		Post:          NewPostService(repository.Post),
		Message:       NewMessageService(repository.Message),
		Chat:          NewChatService(repository.Chat),
	}
}
