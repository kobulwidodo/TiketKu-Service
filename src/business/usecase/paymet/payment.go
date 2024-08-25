package paymet

import (
	"context"
	"encoding/json"
	bookingDom "go-clean/src/business/domain/booking"
	bookingDetailDom "go-clean/src/business/domain/booking_detail"
	midtransDom "go-clean/src/business/domain/midtrans"
	midtransTransactionDom "go-clean/src/business/domain/midtrans_transaction"
	paymentOptionDom "go-clean/src/business/domain/payment_option"
	"go-clean/src/business/entity"
	"go-clean/src/lib/auth"
	"go-clean/src/lib/errors"
	"go-clean/src/lib/midtrans"
	"strconv"
)

type Interface interface {
	Create(ctx context.Context, param entity.CreatePaymentParam) (entity.PaymentRes, error)
}

type payment struct {
	auth                auth.Interface
	midtransTransaction midtransTransactionDom.Interface
	booking             bookingDom.Interface
	bookingDetail       bookingDetailDom.Interface
	paymentOptionDom    paymentOptionDom.Interface
	midtrans            midtransDom.Interface
}

func Init(auth auth.Interface, mdt midtransTransactionDom.Interface, bd bookingDom.Interface, bdd bookingDetailDom.Interface, pod paymentOptionDom.Interface, md midtransDom.Interface) Interface {
	p := &payment{
		auth:                auth,
		midtransTransaction: mdt,
		booking:             bd,
		bookingDetail:       bdd,
		paymentOptionDom:    pod,
		midtrans:            md,
	}

	return p
}

func (p *payment) Create(ctx context.Context, param entity.CreatePaymentParam) (entity.PaymentRes, error) {
	res := entity.PaymentRes{}

	user, err := p.auth.GetUserAuthInfo(ctx)
	if err != nil {
		return res, errors.NewError("failed to get user info", err.Error())
	}

	booking, err := p.booking.Get(entity.BookingParam{
		BookingID: param.BookingID,
	})
	if err != nil {
		return res, errors.NewError("failed to get booking data", err.Error())
	}

	bookingDetails, err := p.bookingDetail.GetList(entity.BookingDetailParam{
		BookingId: booking.ID,
	})
	if err != nil {
		return res, errors.NewError("failed to get booking details", err.Error())
	}

	var totalPrice int64 = 0
	for _, bd := range bookingDetails {
		totalPrice += int64(bd.Price)
	}

	paymentOption, err := p.paymentOptionDom.Get(entity.PaymentOptionParam{
		Code: param.PaymentCode,
	})
	if err != nil {
		return res, errors.NewError("failed to get payment option", err.Error())
	}

	coreApiRes, err := p.midtrans.Create(midtrans.CreateOrderParam{
		BookingID:    booking.BookingID,
		PaymentID:    midtrans.PaymentCode(paymentOption.Code),
		GrossAmount:  totalPrice,
		ItemsDetails: p.convertToItemDetails(bookingDetails),
		CustomerDetails: midtrans.CustomerDetails{
			Email: user.User.Email,
		},
	})
	if err != nil {
		return res, errors.NewError("failed to make payment", err.Error())
	}

	paymentDataMarshalled, err := json.Marshal(coreApiRes)
	if err != nil {
		return res, errors.NewError("failed to store payment data", err.Error())
	}

	_, err = p.midtransTransaction.Create(entity.MidtransTransaction{
		TransactionID: booking.ID,
		MidtransID:    coreApiRes.TransactionID,
		OrderID:       coreApiRes.OrderID,
		PaymentType:   coreApiRes.PaymentType,
		Status:        entity.StatusPending,
		PaymentData:   string(paymentDataMarshalled),
	})
	if err != nil {
		return res, errors.NewError("failed to store payment detail", err.Error())
	}

	err = p.booking.Update(entity.BookingParam{
		ID: booking.ID,
	}, entity.UpdateBookingParam{
		Status: entity.WaitingToPay,
	})
	if err != nil {
		return res, errors.NewError("failed to update booking status", err.Error())
	}

	res.Status = entity.StatusPending

	return res, nil
}

func (p *payment) convertToItemDetails(bookingDetails []entity.BookingDetail) []midtrans.ItemsDetails {
	res := []midtrans.ItemsDetails{}

	for _, bd := range bookingDetails {
		res = append(res, midtrans.ItemsDetails{
			ID:    strconv.Itoa(int(bd.ID)),
			Price: int64(bd.Price),
			Qty:   1,
			Name:  "Tiket",
		})
	}

	return res
}
