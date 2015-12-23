// Code generated by protoc-gen-go.
// source: github.com/crackcomm/evpb/example/pb/message.proto
// DO NOT EDIT!

/*
Package pb is a generated protocol buffer package.

It is generated from these files:
	github.com/crackcomm/evpb/example/pb/message.proto

It has these top-level messages:
	Movie
	Video
*/
package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	evpb "github.com/crackcomm/evpb"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type Movie struct {
	// Id - Id of movie to find.
	Id string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	// Title - Movie title.
	Title string `protobuf:"bytes,2,opt,name=title" json:"title,omitempty"`
	// Year - Movie year.
	Year uint32 `protobuf:"varint,3,opt,name=year" json:"year,omitempty"`
}

func (m *Movie) Reset()                    { *m = Movie{} }
func (m *Movie) String() string            { return proto.CompactTextString(m) }
func (*Movie) ProtoMessage()               {}
func (*Movie) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Video struct {
	Url string `protobuf:"bytes,1,opt,name=url" json:"url,omitempty"`
}

func (m *Video) Reset()                    { *m = Video{} }
func (m *Video) String() string            { return proto.CompactTextString(m) }
func (*Video) ProtoMessage()               {}
func (*Video) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func init() {
	proto.RegisterType((*Movie)(nil), "evpb.example.Movie")
	proto.RegisterType((*Video)(nil), "evpb.example.Video")
}

var topicNameMovie = "evpb.example.Movie"

func ConsumeMovie(q evpb.Interface, h func(*Movie) error) error {
	return q.Consume(topicNameMovie, func(body []byte) error {
		msg := new(Movie)
		if err := proto.Unmarshal(body, msg); err != nil {
			return err
		}
		return h(msg)
	})
}

func SendMovie(q evpb.Interface, msg *Movie) error {
	body, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	return q.Send(topicNameMovie, body)
}

var topicNameVideo = "evpb.example.Video"

func ConsumeVideo(q evpb.Interface, h func(*Video) error) error {
	return q.Consume(topicNameVideo, func(body []byte) error {
		msg := new(Video)
		if err := proto.Unmarshal(body, msg); err != nil {
			return err
		}
		return h(msg)
	})
}

func SendVideo(q evpb.Interface, msg *Video) error {
	body, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	return q.Send(topicNameVideo, body)
}

var fileDescriptor0 = []byte{
	// 155 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x32, 0x4a, 0xcf, 0x2c, 0xc9,
	0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0x4f, 0x2e, 0x4a, 0x4c, 0xce, 0x06, 0x32, 0x72, 0xf5,
	0x53, 0xcb, 0x0a, 0x92, 0xf4, 0x53, 0x2b, 0x12, 0x73, 0x0b, 0x72, 0x52, 0xf5, 0x81, 0xcc, 0xdc,
	0xd4, 0xe2, 0xe2, 0xc4, 0xf4, 0x54, 0xbd, 0x82, 0xa2, 0xfc, 0x92, 0x7c, 0x21, 0x1e, 0x90, 0xb4,
	0x1e, 0x54, 0x5a, 0xc9, 0x80, 0x8b, 0xd5, 0x37, 0xbf, 0x2c, 0x33, 0x55, 0x88, 0x8b, 0x8b, 0x29,
	0x33, 0x45, 0x82, 0x51, 0x81, 0x51, 0x83, 0x53, 0x88, 0x97, 0x8b, 0xb5, 0x24, 0xb3, 0x24, 0x27,
	0x55, 0x82, 0x09, 0xcc, 0xe5, 0xe1, 0x62, 0xa9, 0x4c, 0x4d, 0x2c, 0x92, 0x60, 0x06, 0xf2, 0x78,
	0x95, 0x44, 0xb8, 0x58, 0xc3, 0x32, 0x53, 0x52, 0xf3, 0x85, 0xb8, 0xb9, 0x98, 0x4b, 0x8b, 0x72,
	0x20, 0x5a, 0x9c, 0x58, 0xa2, 0x98, 0x0a, 0x92, 0x92, 0xd8, 0xc0, 0x56, 0x18, 0x03, 0x02, 0x00,
	0x00, 0xff, 0xff, 0x78, 0xc0, 0x23, 0x14, 0x98, 0x00, 0x00, 0x00,
}
