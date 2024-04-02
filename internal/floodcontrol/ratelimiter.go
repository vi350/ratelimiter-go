package floodcontrol

import (
	"context"
	"strconv"
	"time"
)

type RateLimiter struct {
	storage Storage
	limit   int64
	period  time.Duration
}

func NewRateLimiter(storage Storage, period time.Duration, limit int64) *RateLimiter {
	return &RateLimiter{storage: storage, period: period, limit: limit}
}

func (r *RateLimiter) Check(ctx context.Context, userID int64) (bool, error) {
	res, err := r.storage.Get(ctx, strconv.FormatInt(userID, 10), r.limit, r.period)
	if err != nil {
		return false, err
	}
	return !res.Reached, nil
}
