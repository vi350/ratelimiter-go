package redis

import (
	"context"
	redisLibrary "github.com/redis/go-redis/v9"
	"strconv"
	"sync"
	"task/internal/storage"
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
	return NewWithOptions(client, &storage.StorageOptions{
		Prefix: "defpref",
	})
}

func NewWithOptions(client Client, options *storage.StorageOptions) *Storage {
	return &Storage{
		client: client,
		Prefix: options.Prefix,
	}
}

func (s *Storage) Get(ctx context.Context, key string, limit int64, period time.Duration) (storage.Result, error) {
	s.Lock()
	incr := s.client.Incr(ctx, s.Prefix+key).Val()
	if incr == 1 {
		s.client.SetEX(ctx, s.Prefix+key, 1, period)
	}
	s.Unlock() // можно было использовать defer, но в данном случае лок нужен только на одну строку
	reached := incr >= limit
	return storage.Result{
		Remaining: limit - incr,
		Reached:   reached,
	}, nil
}

func (s *Storage) Peek(ctx context.Context, key string, limit int64) (storage.Result, error) {
	s.RLock()
	val := s.client.Get(ctx, s.Prefix+key)
	s.RUnlock() // можно было использовать defer, но в данном случае лок нужен только на одну строку
	incr, err := strconv.ParseInt(val.Val(), 10, 64)
	if err != nil {
		return storage.Result{}, err
	}
	reached := incr >= limit
	return storage.Result{
		Remaining: limit - incr,
		Reached:   reached,
	}, nil
}

func (s *Storage) Reset(ctx context.Context, key string, limit int64) (storage.Result, error) {
	s.Lock()
	_, err := s.client.Del(ctx, s.Prefix+key).Result()
	s.Unlock() // можно было использовать defer, но в данном случае лок нужен только на одну строку
	if err != nil {
		return storage.Result{}, err
	}
	return storage.Result{Remaining: limit}, nil
}
