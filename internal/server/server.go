package server

import (
	"GoServer/internal/handler"
	"GoServer/internal/service/files"
	"GoServer/internal/websocket"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"runtime"
	"time"
)

type Server struct {
	httpServer *fiber.App
	handler    *handler.Handler
}

type ServerInterface interface {
	Run(port string, handler *handler.Handler) error
	ShutDown() error
}

func (s *Server) Run(port string, handler *handler.Handler, websocketClient *websocket.WebsocketClient) error {
	s.httpServer = fiber.New(fiber.Config{
		Prefork:           false,
		StrictRouting:     true,
		CaseSensitive:     true,
		BodyLimit:         files.PostFilesMaxSize + 10*files.KB,
		Concurrency:       runtime.NumCPU(),
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       5 * time.Second,
		ReadBufferSize:    4096,
		WriteBufferSize:   4096,
		ErrorHandler:      nil,
		DisableKeepalive:  false,
		AppName:           "Hey go server",
		StreamRequestBody: true,
		ReduceMemoryUsage: false,
		JSONEncoder:       json.Marshal,
		JSONDecoder:       json.Unmarshal,
		EnablePrintRoutes: true,
	})
	s.handler = handler
	s.handler.InitMiddlewares(s.httpServer)
	s.handler.InitRoutes(s.httpServer, websocketClient)
	handler.Logger.Info("routers have been initialized")
	go websocketClient.Run()
	err := s.httpServer.Listen(port)
	if err != nil {
		return err
	}
	s.handler.Logger.FormatSuccess("Server is running on port %s", port)
	return nil
}

func (s *Server) ShutDown() error {
	return s.httpServer.Shutdown()
}
