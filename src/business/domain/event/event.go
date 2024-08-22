package event

import (
	"go-clean/src/business/entity"

	"gorm.io/gorm"
)

type Interface interface {
	GetList(param entity.EventParam) ([]entity.Event, error)
	Get(param entity.EventParam) (entity.Event, error)
}

type event struct {
	db *gorm.DB
}

func Init(db *gorm.DB) Interface {
	e := &event{
		db: db,
	}

	return e
}

func (e *event) GetList(param entity.EventParam) ([]entity.Event, error) {
	res := []entity.Event{}

	if err := e.db.Where(param).Find(&res).Error; err != nil {
		return res, err
	}

	return res, nil
}

func (e *event) Get(param entity.EventParam) (entity.Event, error) {
	res := entity.Event{}

	if err := e.db.Where(param).First(&res).Error; err != nil {
		return res, err
	}

	return res, nil
}
