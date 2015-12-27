package evpb

import (
	"reflect"

	"github.com/gogo/protobuf/proto"
)

// ProtoConsumer - evpb protobuf message consumer interface
type ProtoConsumer func(proto.Message) error

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
