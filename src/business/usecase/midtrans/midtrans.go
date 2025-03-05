package midtrans

import (
	bookingDomain "go-clean/src/business/domain/booking"
	midtransDomain "go-clean/src/business/domain/midtrans"
	midtransTransactionDomain "go-clean/src/business/domain/midtrans_transaction"
	"go-clean/src/business/entity"
	"go-clean/src/lib/errors"
)

type Interface interface {
	HandleNotification(payload map[string]interface{}) error
}

type midtrans struct {
	midtrans            midtransDomain.Interface
	booking             bookingDomain.Interface
	midtransTransaction midtransTransactionDomain.Interface
}

func Init(md midtransDomain.Interface, bd bookingDomain.Interface, mtt midtransTransactionDomain.Interface) Interface {
	m := &midtrans{
		midtrans:            md,
		booking:             bd,
		midtransTransaction: mtt,
	}

	return m
}

func (md *midtrans) HandleNotification(payload map[string]interface{}) error {
	orderId, exist := payload["order_id"].(string)
	if !exist {
		return errors.NewError("order id not exist", "order id not exist")
	}

	transactionResponse, err := md.midtrans.HandleNotification(orderId)
	if err != nil {
		return err
	}

	midtransTransaction, err := md.midtransTransaction.Get(entity.MidtransTransactionParam{
		OrderID: orderId,
	})
	if err != nil {
		return err
	}

	status := ""

	if transactionResponse != nil {
		// 5. Do set transaction status based on response from check transaction status
		if transactionResponse.TransactionStatus == "capture" {
			if transactionResponse.FraudStatus == "challenge" {
				// TODO set transaction status on your database to 'challenge'
				status = entity.StatusChallange
				// e.g: 'Payment status challenged. Please take action on your Merchant Administration Portal
			} else if transactionResponse.FraudStatus == "accept" {
				// TODO set transaction status on your database to 'success'
				status = entity.StatusSuccess
			}
		} else if transactionResponse.TransactionStatus == "settlement" {
			// TODO set transaction status on your databaase to 'success'
			status = entity.StatusSuccess
		} else if transactionResponse.TransactionStatus == "deny" {
			// TODO you can ignore 'deny', because most of the time it allows payment retries
			// and later can become success
			status = entity.StatusDeny
		} else if transactionResponse.TransactionStatus == "cancel" || transactionResponse.TransactionStatus == "expire" {
			// TODO set transaction status on your databaase to 'failure'
			status = entity.StatusFailure
		} else if transactionResponse.TransactionStatus == "pending" {
			// TODO set transaction status on your databaase to 'pending' / waiting payment
			status = entity.StatusPending
		}
	}

	if err := md.midtransTransaction.Update(entity.MidtransTransactionParam{
		ID: midtransTransaction.ID,
	}, entity.UpdateMidtransTransactionParam{
		Status: status,
	}); err != nil {
		return err
	}

	if status == entity.StatusSuccess {
		if err := md.booking.Update(entity.BookingParam{
			Status: entity.WaitingToPay,
			ID:     midtransTransaction.BookingId,
		}, entity.UpdateBookingParam{
			Status: entity.PaymentSuccessStatus,
		}); err != nil {
			return err
		}
	}

	return nil
}
