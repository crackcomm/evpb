package nsqpb

import "github.com/nsqio/go-nsq"

// Option - Queue option.
type Option func(*queue)

// WithChannel - Sets queue channel.
func WithChannel(channel string) Option {
	return func(c *queue) {
		c.channel = channel
	}
}

// WithConfig - Sets queue nsq config.
func WithConfig(config *nsq.Config) Option {
	return func(c *queue) {
		c.config = config
	}
}

// WithConcurrency - Sets nsq handler concurrency.
func WithConcurrency(concurrency int) Option {
	return func(c *queue) {
		c.concurrency = concurrency
	}
}

// WithAddrs - Sets queue nsq addresses.
func WithAddrs(addrs ...string) Option {
	return func(c *queue) {
		c.addrsNsq = addrs
	}
}

// WithLookupAddrs - Sets queue nsq lookup addresses.
func WithLookupAddrs(addrs ...string) Option {
	return func(c *queue) {
		c.addrsNsqlookup = addrs
	}
}
