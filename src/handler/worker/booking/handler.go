package booking

import (
	"context"
	"encoding/json"
	"fmt"
	"go-clean/src/business/entity"
	"go-clean/src/lib/errors"

	"github.com/nsqio/go-nsq"
)

func (w *worker) HandleMessage(msg *nsq.Message) error {
	var payload entity.BookingTopicPayload
	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		return errors.NewError(err.Error(), err.Error())
	}

	ctx := w.initContext(context.Background(), payload.RequestID)

	w.log.Info(ctx, fmt.Sprintf("processing new message : %#v", payload))

	if err := w.uc.Booking.ProcessBooking(ctx, payload); err != nil {
		return err
	}

	return nil
}
