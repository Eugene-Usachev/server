package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
)

type Server struct {
	httpServer *http.Server
}

type ServerInterface interface {
	Run(port string, handler http.Handler) error
	ShutDown(ctx context.Context) error
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, //1MB
		WriteTimeout:   10 * time.Second,
		ReadTimeout:    10 * time.Second,
	}
	err := s.httpServer.ListenAndServe()
	if err != nil {
		return err
	}
	log.Printf("Server is running on port %s", os.Getenv("PORT"))
	return nil
}

func (s *Server) ShutDown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
