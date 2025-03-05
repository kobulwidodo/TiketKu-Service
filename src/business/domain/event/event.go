package event

import (
	"context"
	"go-clean/src/business/entity"
	tracer "go-clean/src/lib/otel"

	"gorm.io/gorm"
)

type Interface interface {
	GetList(context.Context, entity.EventParam) ([]entity.Event, error)
	Get(param entity.EventParam) (entity.Event, error)
}

type event struct {
	db         *gorm.DB
	oteltracer tracer.Interface
}

func Init(db *gorm.DB, ot tracer.Interface) Interface {
	e := &event{
		db:         db,
		oteltracer: ot,
	}

	return e
}

func (e *event) GetList(ctx context.Context, param entity.EventParam) ([]entity.Event, error) {
	res := []entity.Event{}

	if err := e.db.WithContext(ctx).Where(param).Find(&res).Error; err != nil {
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
