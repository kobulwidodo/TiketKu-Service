package midtrans

import (
	midtransSdk "go-clean/src/lib/midtrans"

	"github.com/midtrans/midtrans-go/coreapi"
)

type Interface interface {
	Create(params midtransSdk.CreateOrderParam) (*coreapi.ChargeResponse, error)
	HandleNotification(id string) (*coreapi.TransactionStatusResponse, error)
}

type midtrans struct {
	m midtransSdk.Interface
}

func Init(m midtransSdk.Interface) Interface {
	ms := &midtrans{
		m: m,
	}

	return ms
}

func (m *midtrans) Create(params midtransSdk.CreateOrderParam) (*coreapi.ChargeResponse, error) {
	result, err := m.m.CreateOrder(params)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (m *midtrans) HandleNotification(id string) (*coreapi.TransactionStatusResponse, error) {
	result, err := m.m.HandleNotification(id)
	if err != nil {
		return result, err
	}

	return result, nil
}
