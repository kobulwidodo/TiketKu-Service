package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"go-clean/src/business/entity"
	"go-clean/src/lib/errors"

	"github.com/nsqio/go-nsq"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func (w *worker) ProcessBooking(msg *nsq.Message) error {
	var payload entity.BookingTopicPayload
	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		return errors.NewError(err.Error(), err.Error())
	}

	propagator := otel.GetTextMapPropagator()
	carrier := propagation.MapCarrier(payload.TraceContext)
	ctx := propagator.Extract(context.Background(), carrier)

	w.log.Info(ctx, fmt.Sprintf("processing new message : %#v", payload))

	if err := w.uc.Booking.ProcessBooking(ctx, payload); err != nil {
		return err
	}

	return nil
}
