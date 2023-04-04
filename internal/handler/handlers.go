package handler

import (
	"GoServer/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (handler *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
	newHub := NewHub()
	go newHub.Run(handler)

	router.Static("/UserFiles", "./static/UserFiles")
	router.Static("/pages", "./static/pages/")

	router.MaxMultipartMemory = 1 << 20

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", handler.signUp)
		auth.POST("/sign-in", handler.signIn)
		auth.POST("/logout", handler.logout)
	}
	api := router.Group("/api")
	{
		userWithoutAuth := api.Group("/user", handler.SetTokensInFirst)
		{
			userWithoutAuth.GET("/:id", handler.getUser)
			userWithoutAuth.GET("/friends_and_subs/:userId", handler.getFriendsAndSubs)
			userWithoutAuth.GET("/many", handler.getUsers)
			userWithoutAuth.GET("/friends_users", handler.getUsersForFriendPage)
		}
		user := api.Group("/user", handler.CheckAuth)
		{
			user.PATCH("/", handler.updateUser)
			user.POST("/avatar/", handler.changeAvatar)
			user.PATCH("/friends/add/:id", handler.addToFriends)
			user.PATCH("/friends/delete/:id", handler.deleteFromFriends)
			user.PATCH("/subscribers/add/:id", handler.addToSubs)
			user.PATCH("/subscribers/delete/:id", handler.deleteFromSubs)
			user.DELETE("/", handler.deleteUser)
		}

		music := api.Group("/music")
		{
			music.GET("/musics", handler.getMusics)
			music.GET("/:id", handler.getMusic)
			music.POST("/", handler.CheckAuth, handler.addMusic)
			//music.PATCH("/:id", handler.updateMusic)
			//music.DELETE("/:id", handler.deleteMusic)
			//TODO playlists
		}

		file := api.Group("/file")
		{
			file.POST("/upload", handler.uploadFile)
		}

		postWithoutAuth := api.Group("/post")
		{
			postWithoutAuth.GET("/:authorId", handler.getPostsByUserID)
		}
		post := api.Group("/post", handler.CheckAuth)
		{
			post.POST("/", handler.createAPost)
			//post.PATCH("/:id", handler.updatePost)
			post.PATCH("/:postId/survey", handler.voteInSurvey)
			post.PATCH("/:postId/likes/like", handler.likePost)
			post.PATCH("/:postId/likes/unlike", handler.unlikePost)
			post.PATCH("/:postId/dislikes/dislike", handler.dislikePost)
			post.PATCH("/:postId/dislikes/undislike", handler.undislikePost)
			post.DELETE("/:postId", handler.deletePost)
		}

		commentWithoutAuth := api.Group("/comment")
		{
			commentWithoutAuth.GET("/:postId", handler.getCommentsByPostId)
		}
		comment := api.Group("/comment", handler.CheckAuth)
		{
			comment.POST("/:postId", handler.createComment)
			comment.PATCH("/:commentId", handler.updateComment)
			comment.PATCH("/:commentId/likes/like", handler.likeComment)
			comment.PATCH("/:commentId/likes/unlike", handler.unlikeComment)
			comment.PATCH("/:commentId/dislikes/dislike", handler.dislikeComment)
			comment.PATCH("/:commentId/dislikes/undislike", handler.undislikeComment)
			comment.DELETE("/:commentId", handler.deleteComment)
		}

		chat := api.Group("/chat", handler.CheckAuth)
		{
			chat.GET("/", handler.getChats)
			chat.PATCH("/chat-list/", handler.UpdateChatLists)
		}
		message := api.Group("/message", handler.CheckAuth)
		{
			message.GET("/last", handler.getLastMessages)
			message.GET("/:id", handler.getMessages)
		}
	}

	router.GET("/ws", func(ctx *gin.Context) {
		newHub.ServeWs(ctx)
	})

	return router
}
