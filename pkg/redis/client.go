package redis

import (
	redisDriver "github.com/redis/rueidis"
	"time"
)

func NewClient(addr []string, pass string) (*redisDriver.Client, error) {
	redisClient, err := redisDriver.NewClient(redisDriver.ClientOption{
		InitAddress:   addr,
		Password:      pass,
		SelectDB:      0,
		MaxFlushDelay: 100 * time.Microsecond,
	})
	return &redisClient, err
}
