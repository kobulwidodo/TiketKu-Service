package booking

import (
	"go-clean/src/business/entity"

	"gorm.io/gorm"
)

type Interface interface {
	Create(booking entity.Booking) error
}

type booking struct {
	db *gorm.DB
}

func Init(db *gorm.DB) Interface {
	b := &booking{
		db: db,
	}

	return b
}

func (b *booking) Create(booking entity.Booking) error {
	if err := b.db.Create(&booking).Error; err != nil {
		return err
	}

	return nil
}
