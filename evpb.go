package evpb

import (
	"reflect"

	"github.com/golang/protobuf/proto"
)

// Consumer - evpb message consumer interface
type Consumer func([]byte) error

// ProtoConsumer - evpb protobuf message consumer interface
type ProtoConsumer func(proto.Message) error

// Interface - evpb interface.
type Interface interface {
	// Send - Sends message to a given topic.
	Send(string, []byte) error

	// Consume - Registers message consumer.
	Consume(string, Consumer) error

	// Stop - Stops consuming and producing.
	Stop() error
}

// Send - Sends protobuf message to evpb.
func Send(evpb Interface, msg proto.Message) error {
	body, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	return evpb.Send(proto.MessageName(msg), body)
}

// Consume - Registers protobuf message consumer.
func Consume(evpb Interface, typ proto.Message, consumer ProtoConsumer) error {
	elem := reflect.TypeOf(typ).Elem()
	return evpb.Consume(proto.MessageName(typ), func(body []byte) error {
		msg := reflect.New(elem).Interface().(proto.Message)
		err := proto.Unmarshal(body, msg)
		if err != nil {
			return err
		}
		return consumer(msg)
	})
}
