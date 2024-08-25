package paymentoption

import (
	"go-clean/src/business/entity"

	paymentOptDom "go-clean/src/business/domain/payment_option"
)

type Interface interface {
	GetList(param entity.PaymentOptionParam) ([]entity.PaymentOption, error)
}

type paymentOption struct {
	paymentOption paymentOptDom.Interface
}

func Init(ed paymentOptDom.Interface) Interface {
	e := &paymentOption{
		paymentOption: ed,
	}

	return e
}

func (e *paymentOption) GetList(param entity.PaymentOptionParam) ([]entity.PaymentOption, error) {
	paymentOptions, err := e.paymentOption.GetList(param)
	if err != nil {
		return paymentOptions, err
	}

	return paymentOptions, nil
}
