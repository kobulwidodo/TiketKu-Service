package booking

import (
	"context"
	"encoding/json"
	"fmt"
	bookingDom "go-clean/src/business/domain/booking"
	seatDom "go-clean/src/business/domain/seat"
	"go-clean/src/business/entity"
	"go-clean/src/lib/auth"
	"go-clean/src/lib/errors"
	"go-clean/src/lib/nsq"
	"time"
)

type Interface interface {
	Create(ctx context.Context, param entity.CreateBookingParam) (entity.BookingResponse, error)
}

type booking struct {
	booking bookingDom.Interface
	seat    seatDom.Interface
	nsq     nsq.Interface
	auth    auth.Interface
}

func Init(auth auth.Interface, bd bookingDom.Interface, sd seatDom.Interface, nsq nsq.Interface) Interface {
	b := &booking{
		auth:    auth,
		booking: bd,
		seat:    sd,
		nsq:     nsq,
	}

	return b
}

func (b *booking) Create(ctx context.Context, param entity.CreateBookingParam) (entity.BookingResponse, error) {
	user, err := b.auth.GetUserAuthInfo(ctx)
	if err != nil {
		return entity.BookingResponse{}, errors.NewError("authenticate failed", err.Error())
	}

	lock, err := b.seat.LockBatchSeat(ctx, param.SeatIDs)
	if err != nil {
		return entity.BookingResponse{}, errors.NewError("failed to process the seat", err.Error())
	}

	defer func() {
		if err != nil {
			_ = b.seat.ReleaseLockBatchSeat(ctx, param.SeatIDs)
		}
	}()

	if !lock {
		return entity.BookingResponse{}, errors.NewError("one of your seat has on process by another user, pick another or wait user cancel their order", "")
	}

	isAvailable, err := b.seat.CheckBatchSeatReserved(param.SeatIDs)
	if err != nil {
		return entity.BookingResponse{}, errors.NewError("failed to get seat status", err.Error())
	}

	if !isAvailable {
		err = errors.NewError("one of your seat has already taken, pick another", "seat has already taken")
		return entity.BookingResponse{}, err
	}

	bookingID := fmt.Sprintf("ID-TIKETKU-%d", time.Now().UnixNano())
	bookingTopicPayload := entity.BookingTopicPayload{
		BookingID:  bookingID,
		UserID:     user.User.ID,
		EventID:    param.EventId,
		CategoryID: param.CategoryId,
		SeatIDs:    param.SeatIDs,
	}
	payloadMarshalled, err := json.Marshal(bookingTopicPayload)
	if err != nil {
		return entity.BookingResponse{}, errors.NewError("there was something wrong, contact the administrator", err.Error())
	}

	err = b.nsq.Publish(entity.BookingTopic, payloadMarshalled)
	if err != nil {
		return entity.BookingResponse{}, errors.NewError("failed to booking, try again later", err.Error())
	}

	return entity.BookingResponse{
		TempBookingID: bookingID,
		Status:        "PROCESSING",
	}, nil
}
