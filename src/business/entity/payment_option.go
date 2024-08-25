package entity

import "gorm.io/gorm"

type PaymentOption struct {
	gorm.Model
	Name string
	Code string
}

type PaymentOptionParam struct {
	Code string
}
