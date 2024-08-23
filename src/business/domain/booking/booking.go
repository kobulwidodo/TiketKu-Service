package booking

import (
	"go-clean/src/business/entity"

	"gorm.io/gorm"
)

type Interface interface {
	Get(param entity.BookingParam) (entity.Booking, error)
	Create(booking entity.Booking) (entity.Booking, error)
	Update(selectParam entity.BookingParam, updateParam entity.UpdateBookingParam) error
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

func (b *booking) Get(param entity.BookingParam) (entity.Booking, error) {
	res := entity.Booking{}
	if err := b.db.Where(param).First(&res).Error; err != nil {
		return res, err
	}

	return res, nil
}

func (b *booking) Create(booking entity.Booking) (entity.Booking, error) {
	if err := b.db.Create(&booking).Error; err != nil {
		return entity.Booking{}, err
	}

	return booking, nil
}

func (b *booking) Update(selectParam entity.BookingParam, updateParam entity.UpdateBookingParam) error {
	if err := b.db.Model(&entity.Booking{}).Where(selectParam).Updates(updateParam).Error; err != nil {
		return err
	}

	return nil
}
