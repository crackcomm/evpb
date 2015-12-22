package nsqpb

import (
	"errors"
	"sync"

	"github.com/golang/protobuf/proto"
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
		addrsNsq:    []string{"127.0.0.1:4150"},
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

func (q *queue) Send(msg proto.Message) (err error) {
	if len(q.addrsNsq) == 0 {
		return errors.New("At least one nsq address is required to create a producer")
	}

	q.producerOnce.Do(func() {
		q.producer, err = nsq.NewProducer(q.addrsNsq[0], q.config)
	})
	if err != nil {
		return
	}

	body, err := proto.Marshal(msg)
	if err != nil {
		return
	}

	topic := proto.MessageName(msg)
	return q.producer.Publish(topic, body)
}

func (q *queue) Consume(msg proto.Message, consumer evpb.Consumer) (err error) {
	topic := proto.MessageName(msg)
	cons, err := nsq.NewConsumer(topic, q.channel, q.config)
	if err != nil {
		return
	}
	cons.AddConcurrentHandlers(&handler{Consumer: consumer, Message: msg}, q.concurrency)
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

type handler struct {
	proto.Message
	evpb.Consumer
}

func (h *handler) HandleMessage(message *nsq.Message) (err error) {
	msg := proto.Clone(h.Message)
	err = proto.Unmarshal(message.Body, msg)
	if err != nil {
		return
	}
	return h.Consumer(msg)
}
