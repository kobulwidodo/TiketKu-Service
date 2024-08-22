package nsq

import (
	"fmt"
	"log"

	"github.com/nsqio/go-nsq"
)

type Interface interface {
	Publish(topic string, data []byte) error
}

type Config struct {
	Host string
	Port string
}

type nsqMq struct {
	conf     Config
	producer *nsq.Producer
}

func Init(cfg Config) Interface {
	n := &nsqMq{
		conf: cfg,
	}
	n.newProducer()
	return n
}

func (n *nsqMq) newProducer() {
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer(fmt.Sprintf("%s:%s", n.conf.Host, n.conf.Port), config)
	if err != nil {
		log.Fatalf("[FATAL] cannot create nsq producer on address @%s:%v, with error: %s", n.conf.Host, n.conf.Port, err)
	}
	n.producer = producer
	log.Printf("NSQ: Procedure Address @%s:%v", n.conf.Host, n.conf.Port)
}

func (n *nsqMq) Publish(topic string, data []byte) error {
	if err := n.producer.Publish(topic, data); err != nil {
		return err
	}

	return nil
}
