package entity

import "gorm.io/gorm"

type BookingDetail struct {
	gorm.Model
	BookingId uint
	SeatId    uint
	Price     uint
}
