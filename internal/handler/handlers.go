package handler

import (
	"GoServer/internal/metrics"
	"GoServer/internal/service"
	"GoServer/internal/websocket"
	"github.com/Eugene-Usachev/fastbytes"
	"github.com/Eugene-Usachev/fst"
	loggerLib "github.com/Eugene-Usachev/logger"
	fiberWebsocket "github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"os"
	"time"
)

var (
	MethodOptionsBytes = fastbytes.S2B("OPTIONS")
)

type Handler struct {
	services         *service.Service
	Logger           *loggerLib.FastLogger
	accessConverter  *fst.Converter
	refreshConverter *fst.Converter
}

type HandlerConfig struct {
	Services         *service.Service
	Logger           *loggerLib.FastLogger
	AccessConverter  *fst.Converter
	RefreshConverter *fst.Converter
}

func NewHandler(cfg *HandlerConfig) *Handler {
	h := &Handler{
		services:         cfg.Services,
		Logger:           cfg.Logger,
		accessConverter:  cfg.AccessConverter,
		refreshConverter: cfg.RefreshConverter,
	}
	return h
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

	disableMetrics := os.Getenv("DISABLE_METRICS")
	if disableMetrics == "true" {
		app.Use(func(c *fiber.Ctx) error {
			startTime := time.Now()
			err := c.Next()
			end := time.Since(startTime).Microseconds()
			method := c.Method()
			path := c.Path()
			statusCode := c.Response().StatusCode()
			if err != nil {
				handler.Logger.FormatInfo("request | %-7s | %-42s | %d | %-12d microseconds |\n", method, path, 500, end)
				return err
			}
			handler.Logger.FormatInfo("request | %-7s | %-42s | %d | %-12d microseconds |\n", method, path, statusCode, end)
			return nil
		})
	} else {
		app.Use(func(c *fiber.Ctx) error {
			defer func() {
				if err := recover(); err != nil {
					handler.Logger.Error("Handled panic in http handler, reason: ", err)
				}
			}()
			startTime := time.Now()
			err := c.Next()
			if err != nil {
				handler.Logger.Error("Handled error in http handler, reason: ", err)
			}
			end := time.Since(startTime).Microseconds()
			method := c.Method()
			path := c.Path()
			statusCode := c.Response().StatusCode()
			handler.Logger.FormatInfo("request | %-7s | %-42s | %d | %-12d microseconds |\n", method, path, statusCode, end)
			metrics.ObserveRequest(float64(end), method, path, statusCode)
			return nil
		})
	}
}

func (handler *Handler) InitRoutes(app *fiber.App, websocketHub *websocket.Hub) {

	app.Static("/UserFiles", "../static/UserFiles")
	app.Static("/pages", "../static/pages/")

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
			userWithoutAuth.Get("/many/", handler.getUsers)
			userWithoutAuth.Get("/friends_users/:clientId", handler.getUsersForFriendPage)
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
			chat.Get("/list/", handler.getChatsList)
			chat.Patch("/list/", handler.UpdateChatsLists)
		}
		message := api.Group("/message", handler.CheckAuth)
		{
			message.Get("/last", handler.getLastMessages)
			message.Get("/:id", handler.getMessages)
		}
	}

	app.Use("/ws", func(c *fiber.Ctx) error {
		if fiberWebsocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws", fiberWebsocket.New(func(conn *fiberWebsocket.Conn) {
		websocketHub.ServeWs(conn)
	}, *websocket.Config))

	metricsHandler := fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler())

	if os.Getenv("DISABLE_METRICS") != "true" {
		app.Get("/metrics", func(ctx *fiber.Ctx) error {
			metricsHandler(ctx.Context())
			return nil
		})
	}
}
