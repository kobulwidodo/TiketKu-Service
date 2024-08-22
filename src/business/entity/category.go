package entity

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	EventId uint
	Name    string
	Price   uint
}

type CategoryParam struct {
	ID      uint
	EventID uint `uri:"event_id"`
}
