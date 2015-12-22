package evpb

import "github.com/golang/protobuf/proto"

// Consumer - evpb message consumer interface
type Consumer func(proto.Message) error

// Interface - evpb interface.
type Interface interface {
	// Send - Sends message.
	Send(proto.Message) error

	// Consume - Registers message consumer.
	Consume(proto.Message, Consumer) error

	// Stop - Stops consuming and producing.
	Stop() error
}
