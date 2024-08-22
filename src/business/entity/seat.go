package entity

import "gorm.io/gorm"

type Seat struct {
	gorm.Model
	EventId    uint
	CategoryId uint
	Row        string
	Number     int
	IsReserved bool
}

type SeatParam struct {
	EventId    uint `uri:"event_id"`
	CategoryId uint `uri:"category_id"`
}

type UpdateSeatParam struct {
	ID         uint
	EventId    uint
	CategoryId uint
	IsReserved bool
}

type SeatResponse struct {
	EventName string
	Category  string
	Rows      []SeatRows
}

type SeatRows struct {
	Row   string
	Seats []SeatListRows
}

type SeatListRows struct {
	SeatID     uint
	Number     int
	IsReserved bool
}
