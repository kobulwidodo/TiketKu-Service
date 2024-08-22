package event

import (
	"go-clean/src/business/entity"

	eventDom "go-clean/src/business/domain/event"
)

type Interface interface {
	GetList(param entity.EventParam) ([]entity.Event, error)
	Get(param entity.EventParam) (entity.Event, error)
}

type event struct {
	event eventDom.Interface
}

func Init(ed eventDom.Interface) Interface {
	e := &event{
		event: ed,
	}

	return e
}

func (e *event) GetList(param entity.EventParam) ([]entity.Event, error) {
	events, err := e.event.GetList(param)
	if err != nil {
		return events, err
	}

	return events, nil
}

func (e *event) Get(param entity.EventParam) (entity.Event, error) {
	event, err := e.event.Get(param)
	if err != nil {
		return event, err
	}

	return event, nil
}
