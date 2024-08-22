package entity

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Name       string
	Date       time.Time
	Venue      string
	TotalSeats int
	IsActive   bool
}

type EventParam struct {
	ID       uint `uri:"event_id"`
	Name     string
	Venue    string
	IsActive bool
}
