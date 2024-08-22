package entity

import (
	"gorm.io/gorm"
)

const (
	BookingTopic                  string = "booking_topic"
	WaitingForSelectPaymentStatus string = "waiting_for_select_payment"
	FailedStatus                  string = "failed"
)

type Booking struct {
	gorm.Model
	BookingID   string
	UserId      uint
	EventId     uint
	TotalAmount int
	Status      string
}

type BookingParam struct {
	ID uint
}

type CreateBookingParam struct {
	EventId    uint   `uri:"event_id"`
	CategoryId uint   `uri:"category_id"`
	SeatIDs    []uint `json:"seat_ids"`
}

type UpdateBookingParam struct {
	Status string
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
	RequestID  string
}
