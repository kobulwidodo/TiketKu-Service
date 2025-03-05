package event

import (
	"context"
	"go-clean/src/business/entity"
	tracer "go-clean/src/lib/otel"

	eventDom "go-clean/src/business/domain/event"
)

type Interface interface {
	GetList(ctx context.Context, param entity.EventParam) ([]entity.Event, error)
	Get(param entity.EventParam) (entity.Event, error)
}

type event struct {
	event      eventDom.Interface
	oteltracer tracer.Interface
}

func Init(ed eventDom.Interface, ot tracer.Interface) Interface {
	e := &event{
		event:      ed,
		oteltracer: ot,
	}

	return e
}

func (e *event) GetList(ctx context.Context, param entity.EventParam) ([]entity.Event, error) {
	ctx, span := e.oteltracer.Start(ctx, "GetListEvent")
	defer span.End()

	events, err := e.event.GetList(ctx, param)
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
