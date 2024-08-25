package midtrans

import (
	midtransSdk "github.com/midtrans/midtrans-go"
)

type PaymentCode string

const (
	GopayPayment    PaymentCode = "gopay"
	TransferBRI     PaymentCode = "transfer-bri"
	TransferMandiri PaymentCode = "transfer-mandiri"
	TransferBNI     PaymentCode = "transfer-bni"
)

type CreateOrderParam struct {
	PaymentID       PaymentCode
	BookingID       string
	GrossAmount     int64
	ItemsDetails    []ItemsDetails
	CustomerDetails CustomerDetails
}

type ItemsDetails struct {
	ID    string
	Price int64
	Qty   int
	Name  string
}

type CustomerDetails struct {
	Name  string
	Email string
}

func (cop *CreateOrderParam) convertToItemDetails() *[]midtransSdk.ItemDetails {
	itemsDetails := []midtransSdk.ItemDetails{}
	for _, i := range cop.ItemsDetails {
		itemDetail := midtransSdk.ItemDetails{
			ID:    i.ID,
			Price: i.Price,
			Qty:   int32(i.Qty),
			Name:  i.Name,
		}
		itemsDetails = append(itemsDetails, itemDetail)
	}

	return &itemsDetails
}
