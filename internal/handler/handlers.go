package handler

import (
	"GoServer/internal/service"
	"GoServer/internal/websocket"
	"github.com/Eugene-Usachev/fastbytes"
	"github.com/Eugene-Usachev/fst"
	"github.com/gofiber/fiber/v2"
	"log"
)

var (
	MethodOptionsBytes = fastbytes.S2B("OPTIONS")
)

type Handler struct {
	services         *service.Service
	accessConverter  *fst.Converter
	refreshConverter *fst.Converter
}

type HandlerConfig struct {
	Services         *service.Service
	AccessConverter  *fst.Converter
	RefreshConverter *fst.Converter
}

func NewHandler(cfg *HandlerConfig) *Handler {
	return &Handler{
		services:         cfg.Services,
		accessConverter:  cfg.AccessConverter,
		refreshConverter: cfg.RefreshConverter,
	}
}

func (handler *Handler) InitMiddlewares(app *fiber.App) {
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Set("Access-Control-Allow-Credentials", "true")
		c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		if fastbytes.Equal(c.Request().Header.Method(), MethodOptionsBytes) {
			return c.SendStatus(204)
		}

		return c.Next()
	})

	//app.Use(recover.New())
}

func (handler *Handler) InitRoutes(app *fiber.App, websocketClient *websocket.WebsocketClient) {

	app.Static("/UserFiles", "./static/UserFiles")
	app.Static("/pages", "./static/pages/")

	auth := app.Group("/auth")
	{
		auth.Post("/sign-up", handler.signUp)
		auth.Post("/sign-in", handler.signIn)
		auth.Post("/refresh", handler.refresh)
		auth.Get("/refresh-tokens/:id", handler.refreshTokens)
	}
	api := app.Group("/api")
	{
		userWithoutAuth := api.Group("/user")
		{
			userWithoutAuth.Get("/:id", handler.getUser)
			userWithoutAuth.Get("/friends_and_subs/:userId", handler.getFriendsAndSubs)
			userWithoutAuth.Get("/many", handler.getUsers)
			userWithoutAuth.Get("/friends_users", handler.getUsersForFriendPage)
		}
		user := api.Group("/user", handler.CheckAuth)
		{
			user.Get("/subs/", handler.getUserSubsIds)
			user.Patch("/", handler.updateUser)
			user.Post("/avatar/", handler.changeAvatar)
			user.Patch("/friends/add/:id", handler.addToFriends)
			user.Patch("/friends/delete/:id", handler.deleteFromFriends)
			user.Patch("/subscribers/add/:id", handler.addToSubs)
			user.Patch("/subscribers/delete/:id", handler.deleteFromSubs)
			user.Delete("/", handler.deleteUser)
		}

		music := api.Group("/music")
		{
			music.Get("/musics", handler.getMusics)
			music.Get("/:id", handler.getMusic)
			music.Post("/", handler.CheckAuth, handler.addMusic)
			//music.Patch("/:id", handler.updateMusic)
			//music.Delete("/:id", handler.deleteMusic)
			//TODO playlists
		}

		file := api.Group("/file")
		{
			file.Post("/upload", handler.uploadFile)
		}

		postWithoutAuth := api.Group("/post")
		{
			postWithoutAuth.Get("/:authorId", handler.getPostsByUserID)
		}
		post := api.Group("/post", handler.CheckAuth)
		{
			post.Post("/", handler.createPost)
			//post.Patch("/:id", handler.updatePost)
			post.Patch("/:postId/survey", handler.voteInSurvey)
			post.Patch("/:postId/likes/like", handler.likePost)
			post.Patch("/:postId/likes/unlike", handler.unlikePost)
			post.Patch("/:postId/dislikes/dislike", handler.dislikePost)
			post.Patch("/:postId/dislikes/undislike", handler.undislikePost)
			post.Delete("/:postId", handler.deletePost)
		}

		commentWithoutAuth := api.Group("/comment")
		{
			commentWithoutAuth.Get("/:postId", handler.getCommentsByPostId)
		}
		comment := api.Group("/comment", handler.CheckAuth)
		{
			comment.Post("/:postId", handler.createComment)
			comment.Patch("/:commentId", handler.updateComment)
			comment.Patch("/:commentId/likes/like", handler.likeComment)
			comment.Patch("/:commentId/likes/unlike", handler.unlikeComment)
			comment.Patch("/:commentId/dislikes/dislike", handler.dislikeComment)
			comment.Patch("/:commentId/dislikes/undislike", handler.undislikeComment)
			comment.Delete("/:commentId", handler.deleteComment)
		}

		chat := api.Group("/chat", handler.CheckAuth)
		{
			chat.Get("/", handler.getChats)
			chat.Patch("/chat-list/", handler.UpdateChatLists)
		}
		message := api.Group("/message", handler.CheckAuth)
		{
			message.Get("/last", handler.getLastMessages)
			message.Get("/:id", handler.getMessages)
		}
	}

	go websocketClient.Run()

	app.Get("/ws", func(ctx *fiber.Ctx) error {
		return websocketClient.ServeWs(ctx.Context())
	})

	log.Println("routers have been initialized")
}
