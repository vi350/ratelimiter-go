package redis

import (
	"context"
	redisLibrary "github.com/redis/go-redis/v9"
	"strconv"
	"sync"
	"task/internal/floodcontrol"
	"time"
)

// Client интерфейс для поддержки всех клиентов редиса (кластерный и соло).
type Client interface {
	Get(ctx context.Context, key string) *redisLibrary.StringCmd
	Incr(ctx context.Context, key string) *redisLibrary.IntCmd
	Del(ctx context.Context, keys ...string) *redisLibrary.IntCmd
	SetEX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redisLibrary.StatusCmd
}

type Storage struct {
	Prefix string
	client Client
	sync.RWMutex
}

func New(client Client) *Storage {
	return NewWithOptions(client, &floodcontrol.StorageOptions{
		Prefix: "defpref",
	})
}

func NewWithOptions(client Client, options *floodcontrol.StorageOptions) *Storage {
	storage := &Storage{
		client: client,
		Prefix: options.Prefix,
	}
	return storage
}

func (s *Storage) Get(ctx context.Context, key string, limit int64, period time.Duration) (floodcontrol.Result, error) {
	s.Lock()
	incr := s.client.Incr(ctx, s.Prefix+key).Val()
	if incr == 1 {
		s.client.SetEX(ctx, s.Prefix+key, 1, period)
	}
	s.Unlock() // можно было использовать defer, но в данном случае лок нужен только на одну строку
	reached := incr >= limit
	return floodcontrol.Result{
		Remaining: limit - incr,
		Reached:   reached,
	}, nil
}

func (s *Storage) Peek(ctx context.Context, key string, limit int64) (floodcontrol.Result, error) {
	s.RLock()
	val := s.client.Get(ctx, s.Prefix+key)
	s.RUnlock() // можно было использовать defer, но в данном случае лок нужен только на одну строку
	incr, err := strconv.ParseInt(val.Val(), 10, 64)
	if err != nil {
		return floodcontrol.Result{}, err
	}
	reached := incr >= limit
	return floodcontrol.Result{
		Remaining: limit - incr,
		Reached:   reached,
	}, nil
}

func (s *Storage) Reset(ctx context.Context, key string, limit int64) (floodcontrol.Result, error) {
	s.Lock()
	_, err := s.client.Del(ctx, s.Prefix+key).Result()
	s.Unlock() // можно было использовать defer, но в данном случае лок нужен только на одну строку
	if err != nil {
		return floodcontrol.Result{}, err
	}
	return floodcontrol.Result{Remaining: limit}, nil
}
