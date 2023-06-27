package service

import (
	"GoServer/Entities"
	"GoServer/internal/repository"
	"context"
	"github.com/Eugene-Usachev/fst"
	"github.com/gofiber/fiber/v2"
	"mime/multipart"
)

type Authorization interface {
	CreateUser(ctx context.Context, dto Entities.UserDTO) (uint, error, Entities.AllTokenResponse)
	SignIn(ctx context.Context, input Entities.SignInDTO) (Entities.SignInReturnDTO, Entities.AllTokenResponse, error)
	Check(ctx context.Context, email, login string) (isEmailNotBusy, isLoginNotBusy bool)
	Refresh(ctx context.Context, id uint, refreshToken string) (Entities.RefreshResponseDTO, error)
	RefreshTokens(ctx context.Context, id uint, refreshToken string) (Entities.AllTokenResponse, error)
}

type User interface {
	GetUserById(ctx context.Context, id uint) (Entities.GetUserDTO, error)
	GetUserSubsIds(ctx context.Context, id uint) ([]uint, error)
	GetFriendsAndSubs(ctx context.Context, clientId, userId uint) (Entities.GetFriendsAndSubsDTO, error)
	UpdateUser(ctx context.Context, id uint, UpdateUserDTO Entities.UpdateUserDTO) error
	ChangeAvatar(ctx *fiber.Ctx, id uint) (string, error)
	AddToFriends(ctx context.Context, id, body uint) error
	DeleteFromFriends(ctx context.Context, id, body uint) error
	AddToSubs(ctx context.Context, id, body uint) error
	DeleteFromSubs(ctx context.Context, id, body uint) error
	DeleteUser(ctx context.Context, id uint) error
	GetUsers(ctx context.Context, idOfUsers string) ([]Entities.MiniUser, error)
	GetUsersForFriendsPage(ctx context.Context, idOfUsers string) ([]Entities.FriendUser, error)

	GetOnlineUsers(ctx context.Context, slice []string) ([]int, error)
	SubscribeOnUsers(ctx context.Context, slice []string, clientId string) error
}

type Post interface {

	/*region post*/

	CreatePost(ctx *fiber.Ctx, id uint, postDTO Entities.CreateAPostDTO, surveyDTO Entities.CreateASurveyDTO, files []*multipart.FileHeader) error
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
	SaveMessage(ctx context.Context, userId uint, messageDTO Entities.MessageDTO) (uint, []uint, string, error)
	UpdateMessage(ctx context.Context, messageId uint, userId uint, newData string) ([]uint, error)
	DeleteMessage(ctx context.Context, messageId uint, userId uint) ([]uint, error)
	GetLastMessages(ctx context.Context, userId uint, chatsId string) ([]Entities.Message, error)
	GetMessages(ctx context.Context, chatId, offset uint) ([20]Entities.Message, error)
}

type Chat interface {
	CreateChat(ctx context.Context, userId uint, chatDTO Entities.ChatDTO) (uint, error)
	UpdateChat(ctx context.Context, userId, chatId uint, chatDTO Entities.ChatUpdateDTO) error
	DeleteChat(ctx context.Context, userId uint, chatId uint) ([]uint, error)
	GetChats(ctx context.Context, userId uint) (string, string, string, []uint, []uint, string, []Entities.Chat, error)
	UpdateChatLists(ctx context.Context, id uint, newChatLists string) error
}

type Music interface {
	GetMusics(ctx context.Context, name string, offset uint) ([]Entities.Music, error)
	GetMusic(ctx context.Context, id uint) (string, string, error)
	AddMusic(ctx *fiber.Ctx, id uint, music Entities.CreateMusicDTO) error
}

type Service struct {
	Authorization
	User
	Post
	Message
	Music
	Chat
}

type ServiceConfig struct {
	Repository       *repository.Repository
	AccessConverter  *fst.Converter
	RefreshConverter *fst.Converter
}

func NewService(cfg *ServiceConfig) *Service {
	return &Service{
		Authorization: NewAuthService(&AuthServiceConfig{
			repository:       cfg.Repository.Authorization,
			accessConverter:  cfg.AccessConverter,
			refreshConverter: cfg.RefreshConverter,
		}),
		User:    NewUserService(cfg.Repository.User),
		Music:   NewMusicService(cfg.Repository.Music),
		Post:    NewPostService(cfg.Repository.Post),
		Message: NewMessageService(cfg.Repository.Message),
		Chat:    NewChatService(cfg.Repository.Chat),
	}
}
