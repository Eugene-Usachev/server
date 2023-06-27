package websocket

import (
	"GoServer/internal/service"
	"context"
	"github.com/Eugene-Usachev/fst"
	"github.com/redis/rueidis"
)

type WebsocketClient struct {
	redis   rueidis.Client
	handler *Handler
	hub     *Hub
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewWebsocketClient(service *service.Service, redisClient *rueidis.Client, accessConverter *fst.Converter) (*WebsocketClient, error) {
	hub := NewHub(accessConverter)
	ctx, cancel := context.WithCancel(context.Background())
	return &WebsocketClient{
		redis:   *redisClient,
		hub:     hub,
		handler: newHandler(service),
		ctx:     ctx,
		cancel:  cancel,
	}, nil
}

func (client *WebsocketClient) Close() {
	client.cancel()
	client.redis.Close()
}
