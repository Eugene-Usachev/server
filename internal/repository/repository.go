// TODO we not a SOLId! We use raw Postgres in repository (not an interface)!

// TODO add name for all returning fields for all methods!
package repository

import (
	"GoServer/Entities"
	"context"
	"github.com/Eugene-Usachev/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/rueidis"
)

type Authorization interface {
	CreateUser(ctx context.Context, user Entities.UserDTO) (uint, error)
	SignInUser(ctx context.Context, input Entities.SignInDTO) (Entities.SignInReturnDTO, error)
	RefreshTokens(ctx context.Context, email, password string) (uint, error)
	Refresh(ctx context.Context, id uint, password string) (dto Entities.RefreshResponseDTO, err error)
	CheckPassword(ctx context.Context, id uint, password string) error
}

type User interface {
	GetUserById(ctx context.Context, id uint) (Entities.GetUserDTO, error)
	GetUserSubsIds(ctx context.Context, id uint) ([]uint, error)
	GetFriendsAndSubs(ctx context.Context, clientId, userId uint) (Entities.GetFriendsAndSubsDTO, error)
	UpdateUser(ctx context.Context, id uint, UpdateUserDTO Entities.UpdateUserDTO) error
	ChangeAvatar(ctx context.Context, id uint, fileName string) error
	AddToFriends(ctx context.Context, id, body uint) error
	DeleteFromFriends(ctx context.Context, id, body uint) error
	AddToSubs(ctx context.Context, id, body uint) error
	DeleteFromSubs(ctx context.Context, id, body uint) error
	DeleteUser(ctx context.Context, id uint) error
	GetUsersForFriendsPage(ctx context.Context, idOfUsers string, clientId uint) ([]Entities.FriendUser, error)
	GetUsers(ctx context.Context, idOfUsers string) ([]Entities.MiniUser, error)
}

type Post interface {
	/*region post*/
	CreatePost(ctx context.Context, id uint, postDTO Entities.CreatePostDTO, surveyDTO Entities.CreateSurveyDTO, date int64) (uint, error)
	GetPostsByUserId(ctx context.Context, authorId uint, offset uint, clientId uint) ([]Entities.GetPostDTO, []Entities.GetSurveyDTO, error)
	LikePost(ctx context.Context, userId, postId uint) error
	UnlikePost(ctx context.Context, userId, postId uint) error
	DislikePost(ctx context.Context, userId, postId uint) error
	UndislikePost(ctx context.Context, userId, postId uint) error
	DeletePost(ctx context.Context, postId, userId uint) error
	/*endregion*/

	/*region comments*/
	GetCommentsByPostId(ctx context.Context, postId uint, offset uint, clientId uint) ([]Entities.Comment, error)
	CreateComment(ctx context.Context, userId uint, comment Entities.CommentDTO, date int64) (uint, error)
	LikeComment(ctx context.Context, userId uint, commentId uint) error
	UnlikeComment(ctx context.Context, userId uint, commentId uint) error
	DislikeComment(ctx context.Context, userId uint, commentId uint) error
	UndislikeComment(ctx context.Context, userId uint, commentId uint) error
	UpdateComment(ctx context.Context, userId uint, commentId uint, updateDTO Entities.CommentUpdateDTO) error
	DeleteComment(ctx context.Context, userId uint, commentId uint) error
	/*endregion*/

	/*region survey*/
	VoteInSurvey(ctx context.Context, userId uint, surveyId uint, votedFor uint16) error
	/*endregion*/
}

type Message interface {
	SaveMessage(ctx context.Context, userId uint, messageDTO Entities.MessageDTO) (uint, []uint, error)
	UpdateMessage(ctx context.Context, messageId, userId uint, newData string) ([]uint, error)
	DeleteMessage(ctx context.Context, messageId, userId uint) ([]uint, error)
	GetLastMessages(ctx context.Context, userId uint, chatsId string) ([]Entities.Message, error)
	GetMessages(ctx context.Context, chatId, offset uint) ([]Entities.Message, error)
}

type Chat interface {
	CreateChat(ctx context.Context, chatDTO Entities.ChatDTO) (uint, error)
	UpdateChat(ctx context.Context, userId, chatId uint, chatDTO Entities.ChatUpdateDTO) error
	DeleteChat(ctx context.Context, userId, chatId uint) ([]uint, error)
	GetChatsListAndInfoForUser(ctx context.Context, userId uint) (friends []uint, chatLists string, rawChats []uint, err error)
	GetChats(ctx context.Context, userId uint, chatsId string) ([]Entities.Chat, error)
	UpdateChatLists(ctx context.Context, id uint, newChatLists string, isSetRawChatsToEmpty bool) error
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

type Postgres struct {
	pool   *pgxpool.Pool
	logger *logger.FastLogger
}

type Redis struct {
	client rueidis.Client
	logger *logger.FastLogger
}

type DataBases struct {
	Postgres *Postgres
	Redis    *Redis
}

func NewDataBases(pool *pgxpool.Pool, postgresLogger *logger.FastLogger, redis rueidis.Client, redisLogger *logger.FastLogger) *DataBases {
	return &DataBases{
		Postgres: &Postgres{pool: pool, logger: postgresLogger},
		Redis:    &Redis{client: redis, logger: redisLogger},
	}
}

func NewRepository(dataBases *DataBases) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(dataBases),
		User:          NewUserPostgres(dataBases),
		Music:         NewMusicPostgres(dataBases),
		Post:          NewPostPostgres(dataBases),
		Chat:          NewChatPostgres(dataBases),
		Message:       NewMessagePostgres(dataBases),
	}
}
