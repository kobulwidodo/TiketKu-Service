package seat

import (
	"context"
	"fmt"
	"time"
)

const (
	lockBookSeatKey = `tiketku:bookseat:lock:%d`
)

func (s *seat) batchLockCache(ctx context.Context, keys []string, expTime time.Duration) (bool, error) {
	ok, err := s.redis.BatchLock(ctx, keys, expTime)
	if !ok || err != nil {
		s.log.Error(ctx, fmt.Sprintf("failed to batchlock %s : %v", keys, err))
	}
	return ok, err
}

func (s *seat) batchReleaseLockCache(ctx context.Context, keys []string) error {
	return s.redis.BatchReleaseLock(ctx, keys)
}
