package service

import (
	"GoServer/Entities"
	"GoServer/internal/repository"
	filesLib "GoServer/internal/service/files"
	"GoServer/pkg/customTime"
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"mime/multipart"
	"path/filepath"
)

type PostService struct {
	repository repository.Post
}

func NewPostService(repository repository.Post) *PostService {
	return &PostService{
		repository: repository,
	}
}

func NewDate() int64 {
	return customTime.Now.Load()
}

/*region post*/

func (service *PostService) CreatePost(ctx *fiber.Ctx, id uint, postDTO Entities.CreatePostDTO, surveyDTO Entities.CreateSurveyDTO, files []*multipart.FileHeader) (uint, error) {
	path := fmt.Sprintf("./static/UserFiles/%d/", id)
	var (
		postFiles []string
	)
	for _, file := range files {
		switch ext := filepath.Ext(file.Filename); ext {
		case ".jpg", ".jpeg", ".png":
			if image, e := filesLib.UploadFile(ctx, file, path+"Images/"); e == nil {
				postFiles = append(postFiles, image)
			} else {
				return 0, errors.New("failed to upload files")
			}
		case ".mp3", ".wav":
			if music, e := filesLib.UploadFile(ctx, file, path+"Musics/"); e == nil {
				postFiles = append(postFiles, music)
			} else {
				return 0, errors.New("failed to upload files")
			}
		case ".mp4", ".avi":
			if video, e := filesLib.UploadFile(ctx, file, path+"Videos/"); e == nil {
				postFiles = append(postFiles, video)
			} else {
				return 0, errors.New("failed to upload files")
			}
		default:
			if other, e := filesLib.UploadFile(ctx, file, path+"Others/"); e == nil {
				postFiles = append(postFiles, other)
			} else {
				return 0, errors.New("failed to upload files")
			}
		}
	}

	postDTO.Files = postFiles
	return service.repository.CreatePost(ctx.Context(), id, postDTO, surveyDTO, NewDate())
}

func (service *PostService) GetPostsByUserID(ctx context.Context, userID uint, offset uint, clientId uint) ([]Entities.GetPostDTO, []Entities.GetSurveyDTO, error) {
	return service.repository.GetPostsByUserID(ctx, userID, offset, clientId)
}

func (service *PostService) LikePost(ctx context.Context, userId, postId uint) error {
	return service.repository.LikePost(ctx, userId, postId)
}
func (service *PostService) UnlikePost(ctx context.Context, userId, postId uint) error {
	return service.repository.UnlikePost(ctx, userId, postId)
}
func (service *PostService) DislikePost(ctx context.Context, userId, postId uint) error {
	return service.repository.DislikePost(ctx, userId, postId)
}
func (service *PostService) UndislikePost(ctx context.Context, userId, postId uint) error {
	return service.repository.UndislikePost(ctx, userId, postId)
}
func (service *PostService) DeletePost(ctx context.Context, postId, userId uint) error {
	return service.repository.DeletePost(ctx, postId, userId)
}

/*endregion*/

/*region comment*/

func (service *PostService) GetCommentsByPostId(ctx context.Context, postId uint, offset uint) ([]Entities.Comment, error) {
	return service.repository.GetCommentsByPostId(ctx, postId, offset)
}
func (service *PostService) CreateComment(ctx context.Context, userId uint, postId uint, comment Entities.CommentDTO) (uint, error) {
	if postId == 0 {
		return 0, errors.New("invalid request body")
	}
	comment.ParentPostID = postId
	return service.repository.CreateComment(ctx, userId, comment, NewDate())
}
func (service *PostService) LikeComment(ctx context.Context, userID uint, commentID uint) error {
	return service.repository.LikeComment(ctx, userID, commentID)
}
func (service *PostService) UnlikeComment(ctx context.Context, userID uint, commentID uint) error {
	return service.repository.UnlikeComment(ctx, userID, commentID)
}
func (service *PostService) DislikeComment(ctx context.Context, userID uint, commentID uint) error {
	return service.repository.DislikeComment(ctx, userID, commentID)
}
func (service *PostService) UndislikeComment(ctx context.Context, userID uint, commentID uint) error {
	return service.repository.UndislikeComment(ctx, userID, commentID)
}
func (service *PostService) UpdateComment(ctx context.Context, userID uint, commentID uint, updateDTO Entities.CommentUpdateDTO) error {
	return service.repository.UpdateComment(ctx, userID, commentID, updateDTO)
}
func (service *PostService) DeleteComment(ctx context.Context, userID uint, commentID uint) error {
	return service.repository.DeleteComment(ctx, userID, commentID)
}

/*endregion*/

/*region Survey*/

func (service *PostService) VoteInSurvey(ctx context.Context, userId uint, surveyId uint, votedFor uint16) error {
	return service.repository.VoteInSurvey(ctx, userId, surveyId, votedFor)
}

/*endregion*/
