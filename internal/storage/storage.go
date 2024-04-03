package storage

import (
	"context"
	"time"
)

type Storage interface {
	// Get увеличивает счётчик на 1, и возвращает лимит для указанного идентификатора.
	Get(ctx context.Context, key string, limit int64, period time.Duration) (Result, error)
	// Peek возвращает лимит для указанного идентификатора.
	Peek(ctx context.Context, key string, limit int64) (Result, error)
	// Reset сбрасывает счётчик для указанного идентификатора.
	Reset(ctx context.Context, key string, limit int64) (Result, error)
}

type StorageOptions struct {
	// Prefix это значение перед идентификатором.
	Prefix string
	// CleanUpInterval это интервал очистки кэша.
	CleanUpInterval time.Duration
}

type Result struct {
	Remaining int64
	Reached   bool
}
