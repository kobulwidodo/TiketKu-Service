package entity

import "gorm.io/gorm"

const (
	BookingTopic string = "booking_topic"
)

type Booking struct {
	gorm.Model
	UserId        uint
	EventId       uint
	TotalAmount   int
	PaymentStatus string
}

type CreateBookingParam struct {
	EventId    uint   `uri:"event_id"`
	CategoryId uint   `uri:"category_id"`
	SeatIDs    []uint `json:"seat_ids"`
}

type BookingResponse struct {
	TempBookingID string
	Status        string
}

type BookingTopicPayload struct {
	BookingID  string
	UserID     uint
	EventID    uint
	CategoryID uint
	SeatIDs    []uint
}
