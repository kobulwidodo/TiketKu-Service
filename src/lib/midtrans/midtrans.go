package midtrans

import (
	"fmt"
	"time"

	midtransSdk "github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

type Interface interface {
	CreateOrder(param CreateOrderParam) (*coreapi.ChargeResponse, error)
	HandleNotification(id string) (*coreapi.TransactionStatusResponse, error)
}

type Config struct {
	ServerKey string
}

type midtrans struct {
	conf    Config
	coreapi *coreapi.Client
}

func Init(cfg Config) Interface {
	m := &midtrans{
		conf: cfg,
	}
	m.connect()
	return m
}

func (m *midtrans) connect() {
	c := coreapi.Client{}
	c.New(m.conf.ServerKey, midtransSdk.Sandbox)
	m.coreapi = &c
}

func (m *midtrans) CreateOrder(param CreateOrderParam) (*coreapi.ChargeResponse, error) {
	fmt.Printf("isi param : %#v\n", param)
	chargeReq := &coreapi.ChargeReq{
		TransactionDetails: midtransSdk.TransactionDetails{
			OrderID:  fmt.Sprintf("%s-%d-%d", "SYN", param.OrderID, time.Now().Unix()),
			GrossAmt: param.GrossAmount,
		},
		Items: param.convertToItemDetails(),
		CustomerDetails: &midtransSdk.CustomerDetails{
			FName: param.CustomerDetails.Name,
			Email: param.CustomerDetails.Email,
		},
	}

	switch param.PaymentID {
	case GopayPayment:
		chargeReq.PaymentType = coreapi.PaymentTypeGopay
	case TransferBRI:
		chargeReq.PaymentType = coreapi.PaymentTypeBankTransfer
		chargeReq.BankTransfer = &coreapi.BankTransferDetails{
			Bank: midtransSdk.BankBri,
		}
	case TransferBNI:
		chargeReq.PaymentType = coreapi.PaymentTypeBankTransfer
		chargeReq.BankTransfer = &coreapi.BankTransferDetails{
			Bank: midtransSdk.BankBni,
		}
	case TransferMandiri:
		chargeReq.PaymentType = coreapi.PaymentTypeEChannel
		chargeReq.EChannel = &coreapi.EChannelDetail{
			BillInfo1: "Payment :",
			BillInfo2: "TiketKu Payment",
		}
	}

	coreApiRes, err := m.coreapi.ChargeTransaction(chargeReq)
	if err != nil {
		return coreApiRes, err
	}

	return coreApiRes, nil
}

func (m *midtrans) HandleNotification(id string) (*coreapi.TransactionStatusResponse, error) {
	midtransReport, err := m.coreapi.CheckTransaction(id)
	if err != nil {
		return midtransReport, err
	}

	return midtransReport, nil
}
