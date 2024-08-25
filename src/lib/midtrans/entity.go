package midtrans

import (
	midtransSdk "github.com/midtrans/midtrans-go"
)

type paymentCode string

const (
	GopayPayment    paymentCode = "gopay"
	TransferBRI     paymentCode = "transfer-bri"
	TransferMandiri paymentCode = "transfer-mandiri"
	TransferBNI     paymentCode = "transfer-bni"
)

type CreateOrderParam struct {
	PaymentID       paymentCode
	OrderID         uint
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
