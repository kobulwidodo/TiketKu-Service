package paymentoption

import (
	"go-clean/src/business/entity"

	"gorm.io/gorm"
)

type Interface interface {
	GetList(param entity.PaymentOptionParam) ([]entity.PaymentOption, error)
	Get(param entity.PaymentOptionParam) (entity.PaymentOption, error)
}

type paymentOption struct {
	db *gorm.DB
}

func Init(db *gorm.DB) Interface {
	c := &paymentOption{
		db: db,
	}

	return c
}

func (c *paymentOption) GetList(param entity.PaymentOptionParam) ([]entity.PaymentOption, error) {
	res := []entity.PaymentOption{}

	if err := c.db.Where(param).Find(&res).Error; err != nil {
		return res, err
	}

	return res, nil
}

func (c *paymentOption) Get(param entity.PaymentOptionParam) (entity.PaymentOption, error) {
	res := entity.PaymentOption{}

	if err := c.db.Where(param).First(&res).Error; err != nil {
		return res, err
	}

	return res, nil
}
