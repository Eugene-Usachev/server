package main

import (
	"GoServer/internal/handler"
	"GoServer/internal/repository"
	"GoServer/internal/server"
	"GoServer/internal/service"
	"GoServer/internal/websocket"
	"GoServer/pkg/redis"
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/Eugene-Usachev/fastbytes"
	"github.com/Eugene-Usachev/fst"
	"log"
	"os"
	"time"
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

	redisClient, err := redis.NewClient([]string{os.Getenv("REDIS_ADDRESS")}, os.Getenv("REDIS_PASSWORD"))
	if err != nil {
		log.Fatalf("error in connection to redis: %s", err)
	}
	log.Println("connected to redis")

	repositoryImpl := repository.NewRepository(repository.NewDataBases(pool, redisClient))

	accessConverter := fst.NewConverter(&fst.ConverterConfig{
		SecretKey:          fastbytes.S2B(os.Getenv("JWT_SECRET_KEY")),
		Postfix:            nil,
		ExpirationTime:     15 * time.Minute,
		HashType:           sha256.New,
		WithExpirationTime: true,
	})
	refreshConverter := fst.NewConverter(&fst.ConverterConfig{
		SecretKey:          fastbytes.S2B(os.Getenv("JWT_SECRET_KEY_FOR_REFRESH_TOKEN")),
		Postfix:            nil,
		ExpirationTime:     31 * time.Hour * 24,
		HashType:           sha256.New,
		WithExpirationTime: true,
	})

	serviceImpl := service.NewService(&service.ServiceConfig{
		Repository:       repositoryImpl,
		AccessConverter:  accessConverter,
		RefreshConverter: refreshConverter,
	})
	websocketClient, err := websocket.NewWebsocketClient(serviceImpl, redisClient, accessConverter)
	if err != nil {
		log.Fatalf("error occured while creating websocket client: %s", err.Error())
	}
	handlerImpl := handler.NewHandler(&handler.HandlerConfig{
		Services:         serviceImpl,
		AccessConverter:  accessConverter,
		RefreshConverter: refreshConverter,
	})

	serverImpl := new(server.Server)

	if err = serverImpl.Run(fmt.Sprintf(":%s", os.Getenv("PORT")), handlerImpl, websocketClient); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}
