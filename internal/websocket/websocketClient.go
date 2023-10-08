package websocket

import (
	"GoServer/internal/service"
	"context"
	"github.com/Eugene-Usachev/fst"
	loggerLib "github.com/Eugene-Usachev/logger"
	"github.com/redis/rueidis"
)

type WebsocketClient struct {
	redis   rueidis.Client
	handler *Handler
	hub     *Hub
	logger  *loggerLib.FastLogger
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewWebsocketClient(service *service.Service, redisClient *rueidis.Client, accessConverter *fst.Converter, logger *loggerLib.FastLogger) (*WebsocketClient, error) {
	hub := NewHub(accessConverter)
	ctx, cancel := context.WithCancel(context.Background())
	return &WebsocketClient{
		redis:   *redisClient,
		hub:     hub,
		logger:  logger,
		handler: newHandler(service),
		ctx:     ctx,
		cancel:  cancel,
	}, nil
}

func (client *WebsocketClient) Close() {
	client.cancel()
	client.redis.Close()
}
