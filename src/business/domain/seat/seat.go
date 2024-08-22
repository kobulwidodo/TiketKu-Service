package seat

import (
	"context"
	"fmt"
	"go-clean/src/business/entity"
	"go-clean/src/lib/log"
	"go-clean/src/lib/redis"
	"time"

	"gorm.io/gorm"
)

type Interface interface {
	GetList(param entity.SeatParam) ([]entity.Seat, error)
	CheckBatchSeatReserved(seatIDs []uint) (bool, error)
	LockBatchSeat(ctx context.Context, seatIDs []uint) (bool, error)
	ReleaseLockBatchSeat(ctx context.Context, seatIDs []uint) error
}

type seat struct {
	db    *gorm.DB
	redis redis.Interface
	log   log.Interface
}

func Init(db *gorm.DB, rd redis.Interface, log log.Interface) Interface {
	s := &seat{
		db:    db,
		redis: rd,
		log:   log,
	}

	return s
}

func (s *seat) GetList(param entity.SeatParam) ([]entity.Seat, error) {
	res := []entity.Seat{}

	if err := s.db.Where(param).Find(&res).Error; err != nil {
		return res, err
	}

	return res, nil
}

func (s *seat) CheckBatchSeatReserved(seatIDs []uint) (bool, error) {
	var count int64
	if err := s.db.Model(&entity.Seat{}).Where("id IN ? AND is_reserved = true", seatIDs).Count(&count).Error; err != nil {
		return false, err
	}

	if count != 0 {
		return false, nil
	}

	return true, nil
}

func (s *seat) LockBatchSeat(ctx context.Context, seatIDs []uint) (bool, error) {
	keys := []string{}

	for _, id := range seatIDs {
		keys = append(keys, fmt.Sprintf(lockBookSeatKey, id))
	}

	return s.batchLockCache(ctx, keys, time.Minute*15)
}

func (s *seat) ReleaseLockBatchSeat(ctx context.Context, seatIDs []uint) error {
	keys := []string{}

	for _, id := range seatIDs {
		keys = append(keys, fmt.Sprintf(lockBookSeatKey, id))
	}

	return s.batchReleaseLockCache(ctx, keys)
}
