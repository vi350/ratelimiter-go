package floodcontrol

import (
	"context"
	"floodcontrol/internal/storage"
	"strconv"
	"time"
)

type RateLimiter struct {
	storage storage.Storage
	limit   int64
	period  time.Duration
}

func NewRateLimiter(storage storage.Storage, limit int64, period time.Duration) *RateLimiter {
	return &RateLimiter{storage: storage, limit: limit, period: period}
}

func (r *RateLimiter) Check(ctx context.Context, userID int64) (bool, error) {
	res, err := r.storage.Get(ctx, strconv.FormatInt(userID, 10), r.limit, r.period)
	if err != nil {
		return false, err
	}
	return !res.Reached, nil
}
