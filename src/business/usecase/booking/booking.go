package booking

import (
	"context"
	"encoding/json"
	"fmt"
	bookingDom "go-clean/src/business/domain/booking"
	bookingDetailDom "go-clean/src/business/domain/booking_detail"
	categoryDom "go-clean/src/business/domain/category"
	seatDom "go-clean/src/business/domain/seat"
	"go-clean/src/business/entity"
	"go-clean/src/lib/appcontext"
	"go-clean/src/lib/auth"
	"go-clean/src/lib/errors"
	"go-clean/src/lib/log"
	"go-clean/src/lib/nsq"
	"time"
)

type Interface interface {
	Create(ctx context.Context, param entity.CreateBookingParam) (entity.BookingResponse, error)
	ProcessBooking(ctx context.Context, param entity.BookingTopicPayload) error
}

type booking struct {
	booking       bookingDom.Interface
	category      categoryDom.Interface
	bookingDetail bookingDetailDom.Interface
	seat          seatDom.Interface
	nsq           nsq.Interface
	auth          auth.Interface
	log           log.Interface
}

func Init(auth auth.Interface, bd bookingDom.Interface, cd categoryDom.Interface, bdd bookingDetailDom.Interface, sd seatDom.Interface, nsq nsq.Interface, log log.Interface) Interface {
	b := &booking{
		auth:          auth,
		booking:       bd,
		category:      cd,
		bookingDetail: bdd,
		seat:          sd,
		nsq:           nsq,
		log:           log,
	}

	return b
}

func (b *booking) Create(ctx context.Context, param entity.CreateBookingParam) (entity.BookingResponse, error) {
	// get user info
	user, err := b.auth.GetUserAuthInfo(ctx)
	if err != nil {
		return entity.BookingResponse{}, errors.NewError("authenticate failed", err.Error())
	}

	// validate the seat
	seats, err := b.seat.GetListByIDs(param.SeatIDs, entity.SeatParam{
		EventId:    param.EventId,
		CategoryId: param.CategoryId,
	})
	if err != nil {
		return entity.BookingResponse{}, errors.NewError("there is something wrong", err.Error())
	}

	if len(seats) != len(param.SeatIDs) {
		return entity.BookingResponse{}, errors.NewError("different categories in the same order", "different categories in the same order")
	}

	// lock the seat on redis
	lock, err := b.seat.LockBatchSeat(ctx, param.SeatIDs)
	if err != nil {
		return entity.BookingResponse{}, errors.NewError("failed to process the seat", err.Error())
	}

	// incase there is an error, release the lock
	defer func() {
		if err != nil {
			_ = b.seat.ReleaseLockBatchSeat(ctx, param.SeatIDs)
		}
	}()

	// if redis failed to lock the seat
	// eg : the seat already got locked by another process
	if !lock {
		return entity.BookingResponse{}, errors.NewError("one of your seat has on process by another user, pick another or wait user cancel their order", "")
	}

	// check seat is_reserved status
	isAvailable, err := b.seat.CheckBatchSeatReserved(param.SeatIDs)
	if err != nil {
		return entity.BookingResponse{}, errors.NewError("failed to get seat status", err.Error())
	}

	// if its not available, than cancel the book
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
		RequestID:  appcontext.GetRequestID(ctx),
	}
	payloadMarshalled, err := json.Marshal(bookingTopicPayload)
	if err != nil {
		return entity.BookingResponse{}, errors.NewError("there was something wrong, contact the administrator", err.Error())
	}

	// publish booking message to nsq
	err = b.nsq.Publish(entity.BookingTopic, payloadMarshalled)
	if err != nil {
		return entity.BookingResponse{}, errors.NewError("failed to booking, try again later", err.Error())
	}

	b.log.Info(ctx, fmt.Sprintf("sucessfully push new booking message : %s", string(payloadMarshalled)))

	return entity.BookingResponse{
		TempBookingID: bookingID,
		Status:        "PROCESSING",
	}, nil
}

func (b *booking) ProcessBooking(ctx context.Context, param entity.BookingTopicPayload) error {
	// set the seat isreserved status to true
	err := b.seat.UpdateStatusBatch(param.SeatIDs, true)
	if err != nil {
		return errors.NewError(err.Error(), err.Error())
	}

	// create booking data in database - waiting for payment status
	bookingData, err := b.booking.Create(entity.Booking{
		BookingID:   param.BookingID,
		UserId:      param.UserID,
		EventId:     param.EventID,
		TotalAmount: len(param.SeatIDs),
		Status:      entity.WaitingForSelectPaymentStatus,
	})
	if err != nil {
		return errors.NewError(err.Error(), err.Error())
	}

	category, err := b.category.Get(entity.CategoryParam{
		ID:      param.CategoryID,
		EventID: param.EventID,
	})
	if err != nil {
		return errors.NewError(err.Error(), err.Error())
	}

	for _, s := range param.SeatIDs {
		if err := b.bookingDetail.Create(entity.BookingDetail{
			BookingId: bookingData.ID,
			SeatId:    s,
			Price:     category.Price,
		}); err != nil {
			return errors.NewError(err.Error(), err.Error())
		}
	}

	// defer unlock the redis lock
	defer func() {
		b.seat.ReleaseLockBatchSeat(ctx, param.SeatIDs)
		// if there is an error in process, then revert every changes
		if err != nil {
			b.seat.UpdateStatusBatch(param.SeatIDs, false)
			b.booking.Update(entity.BookingParam{
				ID: bookingData.ID,
			}, entity.UpdateBookingParam{
				Status: entity.FailedStatus,
			})
			b.log.Error(ctx, err)
		}
	}()

	b.log.Info(ctx, fmt.Sprintf("sucessfully process the message : ", param.BookingID))

	return nil
}
