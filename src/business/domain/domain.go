package domain

import (
	"go-clean/src/business/domain/booking"
	bookingdetail "go-clean/src/business/domain/booking_detail"
	"go-clean/src/business/domain/category"
	"go-clean/src/business/domain/event"
	"go-clean/src/business/domain/midtrans"
	"go-clean/src/business/domain/payment"
	paymentoption "go-clean/src/business/domain/payment_option"
	"go-clean/src/business/domain/seat"
	"go-clean/src/business/domain/user"
	"go-clean/src/lib/log"
	midtransLib "go-clean/src/lib/midtrans"
	"go-clean/src/lib/redis"

	"gorm.io/gorm"
)

type Domains struct {
	User          user.Interface
	Event         event.Interface
	Seat          seat.Interface
	Booking       booking.Interface
	BookingDetail bookingdetail.Interface
	Payment       payment.Interface
	Category      category.Interface
	Midtrans      midtrans.Interface
	PaymentOption paymentoption.Interface
}

func Init(db *gorm.DB, redis redis.Interface, m midtransLib.Interface, log log.Interface) *Domains {
	d := &Domains{
		User:          user.Init(db),
		Event:         event.Init(db),
		Seat:          seat.Init(db, redis, log),
		Booking:       booking.Init(db),
		BookingDetail: bookingdetail.Init(db),
		Payment:       payment.Init(db),
		Category:      category.Init(db),
		Midtrans:      midtrans.Init(m),
		PaymentOption: paymentoption.Init(db),
	}

	return d
}
