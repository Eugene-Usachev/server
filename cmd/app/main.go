package main

import (
	"GoServer/internal/handler"
	"GoServer/internal/repository"
	"GoServer/internal/server"
	"GoServer/internal/service"
	"context"
	"log"
	"os"
)

func main() {

	pool, err := repository.NewPostgresDB(context.Background(), 12,
		repository.Config{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			UserName: os.Getenv("DB_USERNAME"),
			UserPass: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("SSL_MODE"),
		})
	if err != nil {
		log.Fatalf("error in connection to database: %s", err)
	}

	repositoryImpl := repository.NewRepository(pool)
	serviceImpl := service.NewService(repositoryImpl)
	handlerImpl := handler.NewHandler(serviceImpl)

	serverImpl := new(server.Server)

	if err = serverImpl.Run(os.Getenv("PORT"), handlerImpl.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}
