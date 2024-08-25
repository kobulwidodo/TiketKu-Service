package entity

import "gorm.io/gorm"

const (
	StatusChallange = "challange"
	StatusSuccess   = "success"
	StatusDeny      = "deny"
	StatusFailure   = "failure"
	StatusPending   = "pending"
)

type MidtransTransaction struct {
	gorm.Model
	TransactionID uint
	MidtransID    string
	OrderID       string
	PaymentType   string
	Status        string
	PaymentData   string
}

type PaymentData struct {
	Key string `json:"key"`
	Qr  string `json:"qr"`
}

type MidtransTransactionParam struct {
	ID            uint   `json:"id"`
	OrderID       string `uri:"order_id" json:"order_id"`
	TransactionID uint   `uri:"transaction_id"`
}

type UpdateMidtransTransactionParam struct {
	Status string
}

type MidtransTransactionPaymentDetail struct {
	Status      string      `json:"status"`
	MidtransID  string      `json:"midtrans_id"`
	PaymentData PaymentData `json:"payment_data"`
}
