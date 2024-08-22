package payment

import (
	"go-clean/src/business/entity"

	"gorm.io/gorm"
)

type Interface interface {
	Create(payment entity.Payment) error
}

type payment struct {
	db *gorm.DB
}

func Init(db *gorm.DB) Interface {
	p := &payment{
		db: db,
	}

	return p
}

func (p *payment) Create(payment entity.Payment) error {
	if err := p.db.Create(&payment).Error; err != nil {
		return err
	}

	return nil
}
