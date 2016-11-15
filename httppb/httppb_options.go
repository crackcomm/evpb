package httppb

import "github.com/hailocab/go-hostpool"

// Option - Queue option.
type Option func(*queue)

// WithTopicPrefix - Sets topic prefix.
func WithTopicPrefix(topicPrefix string) Option {
	return func(c *queue) {
		c.topicPrefix = topicPrefix
	}
}

// WithBasePath - Sets base http path (/api/v1/httppb is default).
func WithBasePath(basePath string) Option {
	return func(c *queue) {
		c.basePath = basePath
	}
}

// WithAddrs - Sets queue nsq addresses.
func WithAddrs(addrs ...string) Option {
	return func(c *queue) {
		c.hostPool = hostpool.NewEpsilonGreedy(addrs, 0, new(hostpool.LinearEpsilonValueCalculator))
	}
}

// WithServer - Sets listening server address.
func WithServer(listenAddr string) Option {
	return func(c *queue) {
		c.listenAddr = listenAddr
	}
}

// WithConcurrency - Sets consumer handler concurrency.
func WithConcurrency(concurrency int) Option {
	return func(c *queue) {
		c.concurrency = concurrency
	}
}
