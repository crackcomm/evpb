// Package chanpb implements evpb.Interface that sends protobuf marshaled
// messages using channels. It was created for debugging purposes.
// chanpb ignores errors returned by consumers.
package chanpb

import (
	"github.com/crackcomm/evpb"
	"github.com/crackcomm/evpb/syncmap"
)

// Option - Queue option.
type Option func(*options)

// New - Creates new evpb.Interface that transports messages using channels.
func New(opts ...Option) evpb.Interface {
	q := &queue{options: &options{capacity: 10000, concurrency: 1}}
	for _, opt := range opts {
		opt(q.options)
	}
	q.Map = syncmap.New(func(interface{}) (interface{}, error) {
		return make(chan []byte, q.options.capacity), nil
	})
	return q
}

// WithCapacity - Sets queue capacity.
func WithCapacity(capacity int) Option {
	return func(o *options) {
		o.capacity = capacity
	}
}

// WithConcurrency - Sets nsq handler concurrency.
func WithConcurrency(concurrency int) Option {
	return func(o *options) {
		o.concurrency = concurrency
	}
}

type options struct {
	capacity    int
	concurrency int
}

type queue struct {
	*options
	syncmap.Map
}

func (q *queue) Send(topic string, body []byte) (err error) {
	q.channel(topic) <- body
	return
}

func (q *queue) Consume(topic string, consume evpb.Consumer) (err error) {
	channel := q.channel(topic)
	for n := 0; n < q.options.concurrency; n++ {
		go func() {
			for body := range channel {
				consume(body)
			}
		}()
	}
	return
}

func (q *queue) Stop() (err error) {
	q.Map.Lock()
	for _, ch := range q.Map.Values() {
		close(ch.(chan []byte))
	}
	q.Map.Unlock()
	q.Map.Clear()
	return
}

func (q *queue) channel(key string) chan []byte {
	ch, _ := q.Map.Get(key)
	return ch.(chan []byte)
}
