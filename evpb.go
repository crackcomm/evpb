package evpb

// Consumer - evpb message consumer interface
type Consumer func([]byte) error

// Interface - evpb interface.
type Interface interface {
	// Send - Sends message to a given topic.
	Send(string, []byte) error

	// Consume - Registers message consumer.
	Consume(string, Consumer) error

	// Stop - Stops consuming and producing.
	Stop() error
}
