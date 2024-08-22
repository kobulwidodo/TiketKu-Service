package booking

import (
	"context"
	"fmt"
	"go-clean/src/business/usecase"
	"go-clean/src/lib/log"
	"go-clean/src/utils/config"
	"os"
	"os/signal"
	"syscall"

	"github.com/nsqio/go-nsq"
)

type Interface interface {
	Run()
}

type worker struct {
	consumer *nsq.Consumer
	conf     config.WorkerConfig
	log      log.Interface
}

func Init(conf config.WorkerConfig, uc *usecase.Usecase, log log.Interface) Interface {
	w := &worker{}

	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer(conf.Topic, conf.Channel, config)
	if err != nil {
		log.Fatal(context.Background(), fmt.Sprintf("could bot create consumer for topics %s: %v", conf.Topic, err))
	}

	w = &worker{
		consumer: consumer,
		conf:     conf,
		log:      log,
	}

	// add the handler
	w.consumer.AddConcurrentHandlers(w, 3)

	return w
}

func (w *worker) Run() {
	if err := w.consumer.ConnectToNSQLookupd(w.conf.LookupdAddress); err != nil {
		w.log.Fatal(context.Background(), fmt.Sprintf("could bot connect to nsqlookupd %s: %v", w.conf.LookupdAddress, err))
	}
	// Signal handling for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit // Wait for an interrupt signal
	w.log.Info(context.Background(), "Shutting down worker...")

	// Gracefully stop the consumer
	w.consumer.Stop()

	// Wait for all messages to be processed before exiting
	<-w.consumer.StopChan
	w.log.Info(context.Background(), "Worker stopped")
}
