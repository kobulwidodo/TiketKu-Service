package entity

import (
	"gorm.io/gorm"
)

const (
	BookingTopic                  string = "booking_topic"
	WaitingForSelectPaymentStatus string = "waiting_for_select_payment"
	WaitingToPay                  string = "waiting_to_pay"
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
	ID        uint
	BookingID string `uri:"booking_id"`
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
	BookingID string
	Status    string
}

type BookingTopicPayload struct {
	BookingID  string
	UserID     uint
	EventID    uint
	CategoryID uint
	SeatIDs    []uint
	RequestID  string
}

type BookingDetailResponse struct {
	ID           uint
	BookingID    string
	Status       string
	UserId       uint
	EventId      uint
	EventName    string
	EventVenue   string
	CategoryName string
	TotalAmount  int
	TotalPrice   int
	Seat         []BookingSeatDetail
}

type BookingSeatDetail struct {
	SeatID uint
	Row    string
	Number int
	Price  uint
}

type BookingStatusResponse struct {
	BookingID string
	Status    string
}
