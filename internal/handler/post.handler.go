package handler

import (
	"GoServer/Entities"
	utils "GoServer/pkg/fasthttp_utils"
	"github.com/Eugene-Usachev/fastbytes"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

/*region post*/
func getPostAndUserID(c *fiber.Ctx) (uint, uint) {
	userID := c.Locals("userId")
	if userID == "" {
		NewErrorResponse(c, fiber.StatusUnauthorized, "invalid auth token")
		return 0, 0
	}

	postID, err := strconv.ParseUint(c.Params("postId"), 10, 64)
	if err != nil || uint(postID) < 1 {
		NewErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
		return 0, 0
	}

	return userID.(uint), uint(postID)
}

func (handler *Handler) createPost(c *fiber.Ctx) error {
	userID, exist := c.Locals("userId").(uint)
	if !exist || userID < 1 {
		return NewErrorResponse(c, fiber.StatusUnauthorized, "invalid auth token")
	}

	form, err := c.MultipartForm()
	if err != nil {
		handler.Logger.Error("Create post error:", err.Error())
		return NewErrorResponse(c, fiber.StatusBadRequest, "Error processing form data")
	}

	post := form.Value["post"]
	if post == nil {
		handler.Logger.Error("Create post error:", err.Error())
		return NewErrorResponse(c, fiber.StatusBadRequest, "Error: Missing post value")
	}

	var postDTO Entities.CreatePostDTO
	err = json.Unmarshal([]byte(post[0]), &postDTO)
	if err != nil {
		handler.Logger.Error("Create post error:", err.Error())
		return NewErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	survey := form.Value["survey"]
	var surveyDTO Entities.CreateSurveyDTO
	if survey != nil && survey[0] != "" {
		err = json.Unmarshal([]byte(survey[0]), &surveyDTO)
		if surveyDTO.Background > 13 {
			surveyDTO.Background = 0
		}
		if err != nil {
			return NewErrorResponse(c, fiber.StatusBadRequest, err.Error())
		}
	} else {
		surveyDTO = Entities.CreateSurveyDTO{}
	}

	files := form.File["files"]
	if len(files) > 10 {
		files = files[:10]
	}

	var id uint

	id, err = handler.services.Post.CreatePost(c, userID, postDTO, surveyDTO, files)
	if err != nil {
		handler.Logger.Error("create post error: " + err.Error())
		return NewErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"id": id})
}

func (handler *Handler) getPostsByUserID(c *fiber.Ctx) error {
	authorID, err := strconv.ParseUint(c.Params("authorId"), 10, 64)
	if err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	offset, err := strconv.ParseUint(c.Query("offset"), 10, 32)
	if err != nil {
		offset = 0
	}

	if uint(authorID) < 1 {
		return NewErrorResponse(c, fiber.StatusBadRequest, "nothing to get")
	}

	var id uint = 0
	auth := utils.GetAuthorizationHeader(c.Context())
	if auth != nil {
		idB, err := handler.accessConverter.ParseToken(fastbytes.B2S(auth))
		if err == nil {
			id = fastbytes.B2U(idB)
		}
	}

	posts, surveys, err := handler.services.Post.GetPostsByUserID(c.Context(), uint(authorID), uint(offset), id)
	if err != nil {
		handler.Logger.Error("get posts by user id error: " + err.Error())
		return NewErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"posts":   posts,
		"surveys": surveys,
	})
}

func (handler *Handler) likePost(c *fiber.Ctx) error {
	ctx2 := c.Context()
	userID, postID := getPostAndUserID(c)
	err := handler.services.Post.LikePost(ctx2, userID, postID)
	if err != nil {
		handler.Logger.Error("like post error: " + err.Error())
		return NewErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (handler *Handler) unlikePost(c *fiber.Ctx) error {
	ctx2 := c.Context()
	userID, postID := getPostAndUserID(c)
	err := handler.services.Post.UnlikePost(ctx2, userID, postID)
	if err != nil {
		handler.Logger.Error("unlike post error: " + err.Error())
		return NewErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (handler *Handler) dislikePost(c *fiber.Ctx) error {
	ctx2 := c.Context()
	userID, postID := getPostAndUserID(c)
	err := handler.services.Post.DislikePost(ctx2, userID, postID)
	if err != nil {
		handler.Logger.Error("dislike post error: " + err.Error())
		return NewErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (handler *Handler) undislikePost(c *fiber.Ctx) error {
	ctx2 := c.Context()
	userID, postID := getPostAndUserID(c)
	err := handler.services.Post.UndislikePost(ctx2, userID, postID)
	if err != nil {
		handler.Logger.Error("undislike post error: " + err.Error())
		return NewErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (handler *Handler) deletePost(c *fiber.Ctx) error {
	ctx2 := c.Context()
	userID, postID := getPostAndUserID(c)
	err := handler.services.Post.DeletePost(ctx2, userID, postID)
	if err != nil {
		handler.Logger.Error("delete post error: " + err.Error())
		return NewErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

/*endregion*/

/*region comment*/
func getCommentAndUserId(c *fiber.Ctx) (uint, uint) {
	userId := c.Locals("userId").(uint)
	if userId < 1 {
		NewErrorResponse(c, fiber.StatusUnauthorized, "invalid auth token")
		return 0, 0
	}

	commentId, err := strconv.ParseUint(c.Params("commentId"), 10, 64)
	if err != nil || uint(commentId) < 1 {
		NewErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
		return 0, 0
	}

	return userId, uint(commentId)
}

func (handler *Handler) getCommentsByPostId(c *fiber.Ctx) error {
	postId, err := strconv.ParseUint(c.Params("postId"), 10, 64)
	if uint(postId) < 1 {
		return NewErrorResponse(c, fiber.StatusBadRequest, "invalid post id")
	}

	offset, err := strconv.ParseUint(c.Query("offset"), 10, 64)
	if err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, "invalid post id")
	}

	comments, err := handler.services.Post.GetCommentsByPostId(c.Context(), uint(postId), uint(offset))
	if err != nil {
		handler.Logger.Error("get comments by post id error: " + err.Error())
		return NewErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"comments": comments,
	})
}

func (handler *Handler) createComment(c *fiber.Ctx) error {
	userId, postId := getPostAndUserID(c)
	var commentDTO Entities.CommentDTO

	if err := c.BodyParser(&commentDTO); err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	commentId, err := handler.services.Post.CreateComment(c.Context(), userId, postId, commentDTO)
	if err != nil {
		handler.Logger.Error("create comment error: " + err.Error())
		return NewErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"commentId": commentId,
	})
}

func (handler *Handler) likeComment(c *fiber.Ctx) error {
	userId, commentId := getCommentAndUserId(c)

	err := handler.services.Post.LikeComment(c.Context(), userId, commentId)
	if err != nil {
		handler.Logger.Error("like comment error: " + err.Error())
		return NewErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (handler *Handler) unlikeComment(c *fiber.Ctx) error {
	userId, commentId := getCommentAndUserId(c)

	err := handler.services.Post.UnlikeComment(c.Context(), userId, commentId)
	if err != nil {
		handler.Logger.Error("unlike comment error: " + err.Error())
		return NewErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (handler *Handler) dislikeComment(c *fiber.Ctx) error {
	userId, commentId := getCommentAndUserId(c)

	err := handler.services.Post.DislikeComment(c.Context(), userId, commentId)
	if err != nil {
		handler.Logger.Error("dislike comment error: " + err.Error())
		return NewErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (handler *Handler) undislikeComment(c *fiber.Ctx) error {
	userId, commentId := getCommentAndUserId(c)

	err := handler.services.Post.UndislikeComment(c.Context(), userId, commentId)
	if err != nil {
		handler.Logger.Error("undislike comment error: " + err.Error())
		return NewErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (handler *Handler) updateComment(c *fiber.Ctx) error {
	userId, commentId := getCommentAndUserId(c)
	var commentUpdateDTO Entities.CommentUpdateDTO

	if err := c.BodyParser(&commentUpdateDTO); err != nil {
		return NewErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	err := handler.services.Post.UpdateComment(c.Context(), userId, commentId, commentUpdateDTO)
	if err != nil {
		handler.Logger.Error("update comment error: " + err.Error())
		return NewErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (handler *Handler) deleteComment(c *fiber.Ctx) error {
	userId, commentId := getCommentAndUserId(c)

	err := handler.services.Post.DeleteComment(c.Context(), userId, commentId)
	if err != nil {
		handler.Logger.Error("delete comment error: " + err.Error())
		return NewErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

/*endregion*/

/*region survey*/

func getSurveyAndUserId(ctx *fiber.Ctx) (uint, uint) {
	userId := ctx.Locals("userId").(uint)
	if userId < 1 {
		NewErrorResponse(ctx, fiber.StatusBadRequest, "user id is required")
		return 0, 0
	}

	surveyId, err := strconv.ParseUint(ctx.Params("postId"), 10, 64)
	if err != nil {
		NewErrorResponse(ctx, fiber.StatusBadRequest, "post id is required")
		return 0, 0
	}

	return userId, uint(surveyId)
}

func (handler *Handler) voteInSurvey(ctx *fiber.Ctx) error {
	userId, surveyId := getSurveyAndUserId(ctx)

	var votedFor struct {
		VotedFor uint16 `json:"voted_for" binding:"required"`
	}

	err := json.Unmarshal(ctx.Body(), &votedFor)
	if err != nil {
		return NewErrorResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	if votedFor.VotedFor == 0 {
		return ctx.SendStatus(fiber.StatusNoContent)
	}

	err = handler.services.Post.VoteInSurvey(ctx.Context(), userId, surveyId, votedFor.VotedFor)
	if err != nil {
		handler.Logger.Error("vote in survey error: " + err.Error())
		return NewErrorResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}

/*endregion*/
