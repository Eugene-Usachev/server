package main

import (
	"GoServer/internal/handler"
	"GoServer/internal/repository"
	"GoServer/internal/server"
	"GoServer/internal/service"
	"GoServer/internal/websocket"
	"GoServer/pkg/customTime"
	"GoServer/pkg/redis"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/Eugene-Usachev/fastbytes"
	"github.com/Eugene-Usachev/fst"
	"github.com/Eugene-Usachev/logger"
	"os"
	"path/filepath"
	"time"
)

func main() {
	customTime.Start()

	serverLogger, websocketLogger, postgresFastLogger := initLogs()

	pool, err := repository.NewPostgresDB(context.Background(), 12,
		repository.Config{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			UserName: os.Getenv("DB_USERNAME"),
			UserPass: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("SSL_MODE"),
		}, repository.NewPostgresLogger(postgresFastLogger))
	if err != nil {
		serverLogger.FormatFatal("error in connection to database: %s", err)
	}

	redisClient, err := redis.NewClient([]string{os.Getenv("REDIS_ADDRESS")}, os.Getenv("REDIS_PASSWORD"))
	if err != nil {
		serverLogger.FormatFatal("error in connection to redis: %s", err)
	}
	serverLogger.Info("connected to redis")

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
		Logger:           serverLogger,
		AccessConverter:  accessConverter,
		RefreshConverter: refreshConverter,
	})
	websocketClient, err := websocket.NewWebsocketClient(serviceImpl, redisClient, accessConverter, websocketLogger)
	if err != nil {
		serverLogger.FormatFatal("error occured while creating websocket client: %s", err.Error())
	}
	handlerImpl := handler.NewHandler(&handler.HandlerConfig{
		Services:         serviceImpl,
		Logger:           serverLogger,
		AccessConverter:  accessConverter,
		RefreshConverter: refreshConverter,
	})

	serverImpl := new(server.Server)

	if err = serverImpl.Run(fmt.Sprintf(":%s", os.Getenv("PORT")), handlerImpl, websocketClient); err != nil {
		serverLogger.FormatFatal("error occured while running http server: %s", err.Error())
	}
}

func initLogs() (*logger.FastLogger, *logger.FastLogger, *logger.FastLogger) {
	serverLogsDir := filepath.Join("../logs", "server_logs")
	err := os.Mkdir(serverLogsDir, 0777)
	if err != nil && !errors.Is(err, os.ErrExist) {
		panic(fmt.Sprintf("failed to create logs dir: %s", err))
	}
	filePathInfo := filepath.Join(serverLogsDir, "info.txt")
	infoFile, _ := os.OpenFile(filePathInfo, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	filePathError := filepath.Join(serverLogsDir, "error.txt")
	errorFile, _ := os.OpenFile(filePathError, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	filePathWarning := filepath.Join(serverLogsDir, "warning.txt")
	warningFile, _ := os.OpenFile(filePathWarning, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	filePathSuccess := filepath.Join(serverLogsDir, "success.txt")
	successFile, _ := os.OpenFile(filePathSuccess, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	filePathRecord := filepath.Join(serverLogsDir, "record.txt")
	recordFile, _ := os.OpenFile(filePathRecord, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	filePathRaw := filepath.Join(serverLogsDir, "raw.txt")
	rawFile, _ := os.OpenFile(filePathRaw, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	filePathFatal := filepath.Join(serverLogsDir, "fatal.txt")
	fatalFile, _ := os.OpenFile(filePathFatal, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	cfg := &logger.StandardLoggerConfig{
		IsWritingToTheConsole: true,
		ErrorWriter:           errorFile,
		WarningWriter:         warningFile,
		InfoWriter:            infoFile,
		SuccessWriter:         successFile,
		FatalWriter:           fatalFile,
		RecordWriter:          recordFile,
		RawWriter:             rawFile,
		ShowDate:              true,
	}
	handlerLogger := logger.NewFastLogger(&logger.FastLoggerConfig{
		StandardLoggerConfig: *cfg,
		FlushInterval:        200 * time.Millisecond,
		FatalFunc:            nil,
	})

	websocketLogsDir := filepath.Join("../logs", "websocket_logs")
	err = os.Mkdir(websocketLogsDir, 0777)
	if err != nil && !errors.Is(err, os.ErrExist) {
		panic(fmt.Sprintf("failed to create logs dir: %s", err))
	}
	filePathInfo = filepath.Join(websocketLogsDir, "info.txt")
	infoFile, _ = os.OpenFile(filePathInfo, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	filePathError = filepath.Join(websocketLogsDir, "error.txt")
	errorFile, _ = os.OpenFile(filePathError, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	filePathWarning = filepath.Join(websocketLogsDir, "warning.txt")
	warningFile, _ = os.OpenFile(filePathWarning, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	filePathSuccess = filepath.Join(websocketLogsDir, "success.txt")
	successFile, _ = os.OpenFile(filePathSuccess, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	filePathRecord = filepath.Join(websocketLogsDir, "record.txt")
	recordFile, _ = os.OpenFile(filePathRecord, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	filePathRaw = filepath.Join(websocketLogsDir, "raw.txt")
	rawFile, _ = os.OpenFile(filePathRaw, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	filePathFatal = filepath.Join(websocketLogsDir, "fatal.txt")
	fatalFile, _ = os.OpenFile(filePathFatal, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	cfg = &logger.StandardLoggerConfig{
		IsWritingToTheConsole: true,
		ErrorWriter:           errorFile,
		WarningWriter:         warningFile,
		InfoWriter:            infoFile,
		SuccessWriter:         successFile,
		FatalWriter:           fatalFile,
		RecordWriter:          recordFile,
		RawWriter:             rawFile,
		ShowDate:              true,
	}
	websocketLogger := logger.NewFastLogger(&logger.FastLoggerConfig{
		StandardLoggerConfig: *cfg,
		FlushInterval:        1 * time.Second,
		FatalFunc:            nil,
	})

	postgresLogsDir := filepath.Join("../logs", "postgres_logs")
	err = os.Mkdir(postgresLogsDir, 0777)
	if err != nil && !errors.Is(err, os.ErrExist) {
		panic(fmt.Sprintf("failed to create logs dir: %s", err))
	}
	filePathInfo = filepath.Join(postgresLogsDir, "info.txt")
	infoFile, _ = os.OpenFile(filePathInfo, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	filePathError = filepath.Join(postgresLogsDir, "error.txt")
	errorFile, _ = os.OpenFile(filePathError, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	filePathWarning = filepath.Join(postgresLogsDir, "warning.txt")
	warningFile, _ = os.OpenFile(filePathWarning, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	filePathSuccess = filepath.Join(postgresLogsDir, "success.txt")
	successFile, _ = os.OpenFile(filePathSuccess, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	filePathRecord = filepath.Join(postgresLogsDir, "record.txt")
	recordFile, _ = os.OpenFile(filePathRecord, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	filePathRaw = filepath.Join(postgresLogsDir, "raw.txt")
	rawFile, _ = os.OpenFile(filePathRaw, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	filePathFatal = filepath.Join(postgresLogsDir, "fatal.txt")
	fatalFile, _ = os.OpenFile(filePathFatal, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	cfg = &logger.StandardLoggerConfig{
		IsWritingToTheConsole: false,
		ErrorWriter:           errorFile,
		WarningWriter:         warningFile,
		InfoWriter:            infoFile,
		SuccessWriter:         successFile,
		FatalWriter:           fatalFile,
		RecordWriter:          recordFile,
		RawWriter:             rawFile,
		ShowDate:              true,
	}
	postgresLogger := logger.NewFastLogger(&logger.FastLoggerConfig{
		StandardLoggerConfig: *cfg,
		FlushInterval:        1 * time.Second,
		FatalFunc:            nil,
	})

	return handlerLogger, websocketLogger, postgresLogger
}
