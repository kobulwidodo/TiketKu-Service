package bookingdetail

import (
	"go-clean/src/business/entity"

	"gorm.io/gorm"
)

type Interface interface {
	Create(bookingDetail entity.BookingDetail) error
}

type bookingDetail struct {
	db *gorm.DB
}

func Init(db *gorm.DB) Interface {
	bd := &bookingDetail{
		db: db,
	}

	return bd
}

func (bd *bookingDetail) Create(bookingDetail entity.BookingDetail) error {
	if err := bd.db.Create(&bookingDetail).Error; err != nil {
		return err
	}

	return nil
}
