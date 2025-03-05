package booking

import (
	"context"
	"encoding/json"
	"fmt"
	bookingDom "go-clean/src/business/domain/booking"
	bookingDetailDom "go-clean/src/business/domain/booking_detail"
	categoryDom "go-clean/src/business/domain/category"
	eventDom "go-clean/src/business/domain/event"
	seatDom "go-clean/src/business/domain/seat"
	"go-clean/src/business/entity"
	"go-clean/src/lib/appcontext"
	"go-clean/src/lib/auth"
	"go-clean/src/lib/errors"
	"go-clean/src/lib/log"
	"go-clean/src/lib/nsq"
	tracer "go-clean/src/lib/otel"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type Interface interface {
	Create(ctx context.Context, param entity.CreateBookingParam) (entity.BookingResponse, error)
	ProcessBooking(ctx context.Context, param entity.BookingTopicPayload) error
	Get(ctx context.Context, param entity.BookingParam) (entity.BookingDetailResponse, error)
	CheckStatus(ctx context.Context, param entity.BookingParam) (entity.BookingStatusResponse, error)
}

type booking struct {
	booking       bookingDom.Interface
	category      categoryDom.Interface
	bookingDetail bookingDetailDom.Interface
	event         eventDom.Interface
	seat          seatDom.Interface
	nsq           nsq.Interface
	auth          auth.Interface
	log           log.Interface
	oteltracer    tracer.Interface
}

func Init(auth auth.Interface, bd bookingDom.Interface, cd categoryDom.Interface, bdd bookingDetailDom.Interface, sd seatDom.Interface, ed eventDom.Interface, nsq nsq.Interface, log log.Interface, ot tracer.Interface) Interface {
	b := &booking{
		auth:          auth,
		booking:       bd,
		category:      cd,
		bookingDetail: bdd,
		event:         ed,
		seat:          sd,
		nsq:           nsq,
		log:           log,
		oteltracer:    ot,
	}

	return b
}

func (b *booking) Create(ctx context.Context, param entity.CreateBookingParam) (entity.BookingResponse, error) {
	ctx, span := b.oteltracer.Start(ctx, "Create")
	defer span.End()

	// get user info
	user, err := b.auth.GetUserAuthInfo(ctx)
	if err != nil {
		err := errors.NewError("authenticate failed", err.Error())
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return entity.BookingResponse{}, err
	}

	// validate the seat
	seats, err := b.seat.GetListByIDs(ctx, param.SeatIDs, entity.SeatParam{
		EventId:    param.EventId,
		CategoryId: param.CategoryId,
	})
	if err != nil {
		err := errors.NewError("there is something wrong", err.Error())
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return entity.BookingResponse{}, err
	}

	if len(seats) != len(param.SeatIDs) {
		err := errors.NewError("different categories in the same order", "different categories in the same order")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return entity.BookingResponse{}, err
	}

	// lock the seat on redis
	lock, err := b.seat.LockBatchSeat(ctx, param.SeatIDs)
	if err != nil {
		err := errors.NewError("failed to process the seat", err.Error())
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return entity.BookingResponse{}, err
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
	isAvailable, err := b.seat.CheckBatchSeatReserved(ctx, param.SeatIDs)
	if err != nil {
		err := errors.NewError("failed to get seat status", err.Error())
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return entity.BookingResponse{}, err
	}

	span.AddEvent("Seat Validation Complete", trace.WithAttributes(
		attribute.Bool("isAvailable", isAvailable),
	))

	// if its not available, than cancel the book
	if !isAvailable {
		err = errors.NewError("one of your seat has already taken, pick another", "seat has already taken")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
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

	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	bookingTopicPayload.TraceContext = carrier

	payloadMarshalled, err := json.Marshal(bookingTopicPayload)
	if err != nil {
		err = errors.NewError("there was something wrong, contact the administrator", err.Error())
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return entity.BookingResponse{}, err
	}

	// publish booking message to nsq
	err = b.nsq.Publish(entity.BookingTopic, payloadMarshalled)
	if err != nil {
		err = errors.NewError("failed to booking, try again later", err.Error())
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return entity.BookingResponse{}, err
	}

	span.AddEvent("Booking published to NSQ", trace.WithAttributes(
		attribute.String("payload", string(payloadMarshalled)),
	))

	b.log.Info(ctx, fmt.Sprintf("sucessfully push new booking message : %s", string(payloadMarshalled)))

	return entity.BookingResponse{
		BookingID: bookingID,
		Status:    "PROCESSING",
	}, nil
}

func (b *booking) ProcessBooking(ctx context.Context, param entity.BookingTopicPayload) error {
	ctx, span := b.oteltracer.Start(ctx, "ProcessBooking")
	defer span.End()

	// set the seat isreserved status to true
	err := b.seat.UpdateStatusBatch(ctx, param.SeatIDs, true)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return errors.NewError(err.Error(), err.Error())
	}

	// create booking data in database - waiting for payment status
	bookingData, err := b.booking.Create(ctx, entity.Booking{
		BookingID:   param.BookingID,
		UserId:      param.UserID,
		EventId:     param.EventID,
		TotalAmount: len(param.SeatIDs),
		Status:      entity.WaitingForSelectPaymentStatus,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return errors.NewError(err.Error(), err.Error())
	}

	category, err := b.category.Get(ctx, entity.CategoryParam{
		ID:      param.CategoryID,
		EventID: param.EventID,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return errors.NewError(err.Error(), err.Error())
	}

	for _, s := range param.SeatIDs {
		if err := b.bookingDetail.Create(ctx, entity.BookingDetail{
			BookingId: bookingData.ID,
			SeatId:    s,
			Price:     category.Price,
		}); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return errors.NewError(err.Error(), err.Error())
		}
	}

	// defer unlock the redis lock
	defer func() {
		b.seat.ReleaseLockBatchSeat(ctx, param.SeatIDs)
		// if there is an error in process, then revert every changes
		if err != nil {
			b.seat.UpdateStatusBatch(ctx, param.SeatIDs, false)
			b.booking.Update(entity.BookingParam{
				ID: bookingData.ID,
			}, entity.UpdateBookingParam{
				Status: entity.FailedStatus,
			})
			b.log.Error(ctx, err)
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	b.log.Info(ctx, fmt.Sprintf("sucessfully process the message : ", param.BookingID))

	return nil
}

func (b *booking) Get(ctx context.Context, param entity.BookingParam) (entity.BookingDetailResponse, error) {
	res := entity.BookingDetailResponse{}

	booking, err := b.booking.Get(param)
	if err != nil {
		return res, errors.NewError("booking data is not available", err.Error())
	}

	bookingDetails, err := b.bookingDetail.GetList(entity.BookingDetailParam{
		BookingId: booking.ID,
	})
	if err != nil {
		return res, errors.NewError("failed to get booking details", err.Error())
	}

	if len(bookingDetails) == 0 {
		errmsg := "there is no booking details"
		return res, errors.NewError(errmsg, errmsg)
	}

	seatIDs := []uint{}
	for _, b := range bookingDetails {
		seatIDs = append(seatIDs, b.SeatId)
	}

	seats, err := b.seat.GetListByIDs(ctx, seatIDs, entity.SeatParam{})
	if err != nil {
		return res, errors.NewError("failed to get seat data", err.Error())
	}

	event, err := b.event.Get(entity.EventParam{
		ID: booking.EventId,
	})
	if err != nil {
		return res, errors.NewError("failed to get event detail", err.Error())
	}

	category, err := b.category.Get(ctx, entity.CategoryParam{
		ID:      booking.ID,
		EventID: booking.EventId,
	})
	if err != nil {
		return res, errors.NewError("failed to get category detail", err.Error())
	}

	res.ID = booking.ID
	res.BookingID = booking.BookingID
	res.Status = booking.Status
	res.UserId = booking.UserId
	res.EventId = booking.EventId
	res.EventName = event.Name
	res.EventVenue = event.Venue
	res.CategoryName = category.Name
	res.TotalAmount = booking.TotalAmount
	res.TotalPrice = booking.TotalAmount * int(category.Price)

	for _, s := range seats {
		res.Seat = append(res.Seat, entity.BookingSeatDetail{
			SeatID: s.ID,
			Row:    s.Row,
			Number: s.Number,
			Price:  category.Price,
		})
	}

	return res, nil
}

func (b *booking) CheckStatus(ctx context.Context, param entity.BookingParam) (entity.BookingStatusResponse, error) {
	res := entity.BookingStatusResponse{}

	booking, err := b.booking.Get(entity.BookingParam{
		BookingID: param.BookingID,
	})
	if err != nil {
		return res, errors.NewError("failed to get booking data", err.Error())
	}

	res.BookingID = booking.BookingID
	res.Status = booking.Status

	return res, nil
}
