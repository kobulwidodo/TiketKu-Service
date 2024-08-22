package usecase

import (
	"go-clean/src/business/domain"
	"go-clean/src/business/usecase/booking"
	"go-clean/src/business/usecase/category"
	"go-clean/src/business/usecase/event"
	"go-clean/src/business/usecase/seat"
	"go-clean/src/business/usecase/user"
	"go-clean/src/lib/auth"
	"go-clean/src/lib/nsq"
)

type Usecase struct {
	User     user.Interface
	Event    event.Interface
	Category category.Interface
	Seat     seat.Interface
	Booking  booking.Interface
}

func Init(auth auth.Interface, d *domain.Domains, nsq nsq.Interface) *Usecase {
	uc := &Usecase{
		User:     user.Init(d.User, auth),
		Event:    event.Init(d.Event),
		Category: category.Init(d.Category),
		Seat:     seat.Init(d.Seat, d.Category, d.Event),
		Booking:  booking.Init(auth, d.Booking, d.Seat, nsq),
	}

	return uc
}
