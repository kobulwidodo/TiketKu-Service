package booking

import (
	"context"
	"fmt"

	"github.com/nsqio/go-nsq"
)

func (w *worker) HandleMessage(msg *nsq.Message) error {
	w.log.Info(context.Background(), fmt.Sprintf("accept new message : %s", string(msg.Body)))
	return nil
}
