package entity

import "gorm.io/gorm"

type Payment struct {
	gorm.Model
	BookingId uint
	Amount    uint
	Method    string
	Status    string
}
