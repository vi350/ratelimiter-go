package redis

import (
	"github.com/redis/go-redis/v9"
)

type Client struct {
	*redis.Client
}

func NewClient(host, port, password string) *Client {
	return &Client{
		redis.NewClient(&redis.Options{
			Addr:     host + ":" + port,
			Password: password,
		}),
	}
}
