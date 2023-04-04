package handler

import (
	"GoServer/Entities"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

/*region post*/
func getPostAndUserId(ctx *gin.Context) (uint, uint) {
	userId, exist := ctx.Get("userId")
	if !exist || userId.(uint) < 1 {
		NewErrorResponse(ctx, http.StatusUnauthorized, "invalid auth token")
		return 0, 0
	}

	postId, err := strconv.ParseUint(ctx.Param("postId"), 10, 64)
	if err != nil || uint(postId) < 1 {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return 0, 0
	}

	return userId.(uint), uint(postId)
}

func (handler *Handler) createAPost(ctx *gin.Context) {
	userId, exist := ctx.Get("userId")
	if !exist || userId.(uint) < 1 {
		NewErrorResponse(ctx, http.StatusUnauthorized, "invalid auth token")
		return
	}
	form, err := ctx.MultipartForm()
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "Error processing form data")
		return
	}

	post := form.Value["post"]
	if post == nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "Error: Missing post value")
		return
	}

	var PostDTO Entities.CreateAPostDTO
	err = json.Unmarshal([]byte(post[0]), &PostDTO)
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	survey := form.Value["survey"]
	var SurveyDTO Entities.CreateASurveyDTO
	if survey != nil && survey[0] != "" {
		err = json.Unmarshal([]byte(survey[0]), &SurveyDTO)
		if err != nil {
			NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
			return
		}
	} else {
		SurveyDTO = Entities.CreateASurveyDTO{}
	}

	files := form.File["files"]
	if len(files) > 10 {
		files = files[:10]
	}

	err = handler.services.CreateAPost(ctx, userId.(uint), PostDTO, SurveyDTO, files)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Post and files created successfully"})
}

func (handler *Handler) getPostsByUserID(ctx *gin.Context) {

	authorId, err := strconv.ParseUint(ctx.Param("authorId"), 10, 64)
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	var offset uint64
	offset, err = strconv.ParseUint(ctx.Query("offset"), 10, 32)

	if uint(authorId) < 1 {
		NewErrorResponse(ctx, http.StatusBadRequest, "nothing to get")
		return
	}

	posts, surveys, err := handler.services.GetPostsByUserID(ctx.Request.Context(), uint(authorId), uint(offset))
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"posts":   posts,
		"surveys": surveys,
	})
}

func (handler *Handler) likePost(ctx *gin.Context) {
	ctx2 := ctx.Request.Context()
	userId, postId := getPostAndUserId(ctx)
	err := handler.services.LikePost(ctx2, userId, postId)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (handler *Handler) unlikePost(ctx *gin.Context) {
	ctx2 := ctx.Request.Context()
	userId, postId := getPostAndUserId(ctx)
	err := handler.services.UnlikePost(ctx2, userId, postId)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (handler *Handler) dislikePost(ctx *gin.Context) {
	ctx2 := ctx.Request.Context()
	userId, postId := getPostAndUserId(ctx)
	err := handler.services.DislikePost(ctx2, userId, postId)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (handler *Handler) undislikePost(ctx *gin.Context) {
	ctx2 := ctx.Request.Context()
	userId, postId := getPostAndUserId(ctx)
	err := handler.services.UndislikePost(ctx2, userId, postId)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
}

func (handler *Handler) deletePost(ctx *gin.Context) {
	ctx2 := ctx.Request.Context()
	userId, postId := getPostAndUserId(ctx)
	err := handler.services.DeletePost(ctx2, userId, postId)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
}

/*endregion*/

/*region comment*/
func getCommentAndUserId(ctx *gin.Context) (uint, uint) {
	userId, exist := ctx.Get("userId")
	if !exist || userId.(uint) < 1 {
		NewErrorResponse(ctx, http.StatusUnauthorized, "invalid auth token")
		return 0, 0
	}

	commentId, err := strconv.ParseUint(ctx.Param("commentId"), 10, 64)
	if err != nil || uint(commentId) < 1 {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
		return 0, 0
	}

	return userId.(uint), uint(commentId)
}

func (handler *Handler) getCommentsByPostId(ctx *gin.Context) {

	postId, err := strconv.ParseUint(ctx.Param("postId"), 10, 64)
	if uint(postId) < 1 {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid post id")
		return
	}

	offset, err := strconv.ParseUint(ctx.Query("offset"), 10, 64)
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid post id")
	}

	comments, err := handler.services.GetCommentsByPostId(ctx.Request.Context(), uint(postId), uint(offset))
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"comments": comments,
	})
}
func (handler *Handler) createComment(ctx *gin.Context) {

	userId, postId := getPostAndUserId(ctx)
	var commentDTO Entities.CommentDTO

	if err := ctx.BindJSON(&commentDTO); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	commentId, err := handler.services.CreateComment(ctx.Request.Context(), userId, postId, commentDTO)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"commentId": commentId,
	})
}
func (handler *Handler) likeComment(ctx *gin.Context) {
	userId, commentId := getCommentAndUserId(ctx)

	err := handler.services.LikeComment(ctx.Request.Context(), userId, commentId)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}
func (handler *Handler) unlikeComment(ctx *gin.Context) {
	userId, commentId := getCommentAndUserId(ctx)

	err := handler.services.UnlikeComment(ctx.Request.Context(), userId, commentId)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}
func (handler *Handler) dislikeComment(ctx *gin.Context) {
	userId, commentId := getCommentAndUserId(ctx)

	err := handler.services.DislikeComment(ctx.Request.Context(), userId, commentId)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}
func (handler *Handler) undislikeComment(ctx *gin.Context) {
	userId, commentId := getCommentAndUserId(ctx)

	err := handler.services.UndislikeComment(ctx.Request.Context(), userId, commentId)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}
func (handler *Handler) updateComment(ctx *gin.Context) {
	userId, commentId := getCommentAndUserId(ctx)
	var commentUpdateDTO Entities.CommentUpdateDTO

	if err := ctx.BindJSON(&commentUpdateDTO); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err := handler.services.UpdateComment(ctx.Request.Context(), userId, commentId, commentUpdateDTO)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}
func (handler *Handler) deleteComment(ctx *gin.Context) {
	userId, commentId := getCommentAndUserId(ctx)

	err := handler.services.DeleteComment(ctx.Request.Context(), userId, commentId)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{})
}

/*endregion*/

/*region survey*/

func getSurveyAndUserId(ctx *gin.Context) (uint, uint) {
	userId, exist := ctx.Get("userId")
	if !exist || userId.(uint) < 1 {
		NewErrorResponse(ctx, http.StatusBadRequest, "user id is required")
		return 0, 0
	}

	surveyId, err := strconv.ParseUint(ctx.Param("postId"), 10, 64)
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "post id is required")
		return 0, 0
	}

	return userId.(uint), uint(surveyId)
}

func (handler *Handler) voteInSurvey(ctx *gin.Context) {
	userId, surveyId := getSurveyAndUserId(ctx)

	var votedFor struct {
		VotedFor []uint8 `json:"voted_for" binding:"required"`
	}

	if err := ctx.BindJSON(&votedFor); err != nil || len(votedFor.VotedFor) == 0 {
		NewErrorResponse(ctx, http.StatusBadRequest, `invalid voted for`)
		return
	}

	err := handler.services.VoteInSurvey(ctx.Request.Context(), userId, surveyId, votedFor.VotedFor)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{})
}

/*endregion*/
