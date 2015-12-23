package nsqpb

import (
	"errors"
	"strings"
	"sync"

	"github.com/nsqio/go-nsq"

	"github.com/crackcomm/evpb"
)

var defaultChannel = "default"

// New - Creates evpb on nsq.
func New(opts ...Option) evpb.Interface {
	c := &queue{
		concurrency: 1,
		consumers:   make(map[string]*nsq.Consumer),
		channel:     defaultChannel,
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.config == nil {
		c.config = nsq.NewConfig()
	}
	return c
}

type queue struct {
	config *nsq.Config

	topicPrefix string

	addrsNsqlookup []string
	addrsNsq       []string

	// Consumer
	concurrency int
	consumers   map[string]*nsq.Consumer
	channel     string

	// Producer
	producerOnce sync.Once
	producer     *nsq.Producer
}

func (q *queue) Send(topic string, body []byte) (err error) {
	if len(q.addrsNsq) == 0 {
		return errors.New("At least one nsq address is required to create a producer")
	}
	q.producerOnce.Do(func() {
		q.producer, err = nsq.NewProducer(q.addrsNsq[0], q.config)
	})
	if err != nil {
		return
	}
	return q.producer.Publish(q.topic(topic), body)
}

func (q *queue) Consume(topic string, consumer evpb.Consumer) (err error) {
	cons, err := nsq.NewConsumer(q.topic(topic), q.channel, q.config)
	if err != nil {
		return
	}
	cons.AddConcurrentHandlers(nsq.HandlerFunc(func(msg *nsq.Message) error {
		return consumer(msg.Body)
	}), q.concurrency)
	err = cons.ConnectToNSQLookupds(q.addrsNsqlookup)
	if err != nil {
		return
	}
	err = cons.ConnectToNSQDs(q.addrsNsq)
	if err != nil {
		return
	}
	q.consumers[topic] = cons
	return nil
}

func (q *queue) Stop() (err error) {
	for _, cons := range q.consumers {
		cons.Stop()
	}
	return
}

func (q *queue) topic(t string) string {
	if q.topicPrefix == "" {
		return t
	}
	return strings.Join([]string{q.topicPrefix, t}, "")
}
