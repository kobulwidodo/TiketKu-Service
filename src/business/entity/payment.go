package entity

type CreatePaymentParam struct {
	PaymentCode string `json:"payment_code" binding:"required"`
	BookingID   string `json:"booking_id" binding:"required"`
}

type PaymentRes struct {
	Status string
}
