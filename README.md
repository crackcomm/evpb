# evpb

[![GoDoc](https://godoc.org/github.com/crackcomm/evpb?status.svg)](https://godoc.org/github.com/crackcomm/evpb) [![Circle CI](https://img.shields.io/circleci/project/crackcomm/evpb.svg)](https://circleci.com/gh/crackcomm/evpb)

[Protobuf](https://github.com/golang/protobuf/) events production and consumption.

This library was designed only with protobuf messages in mind but it has `[]byte`
interfaces so any format is pluggable.

We are shipping with `protoc-gen-go` with builtin `evpb` plugin.

## Usage

Protobuf messages do not have to be changed for this to work,
but You can use our `protoc-gen-go` generator like in example [Makefile](https://github.com/crackcomm/evpb/blob/master/Makefile)
to generate a static typed consumer and sender functions.

```Go
package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/crackcomm/evpb/nsqpb"

	pb "github.com/crackcomm/evpb/example/pb"
)

func consumer(movie *pb.Movie) (err error) {
	log.Printf("title=%q year=%d", movie.Title, movie.Year)
	return
}

func main() {
	// Create new NSQ consumer with default nsq addr
	queue := nsqpb.New(
		nsqpb.WithAddrs("127.0.0.1:4150"),
	)

	// Register movies consumer
	// without generating: would need a type casting
	if err := pb.ConsumeMovie(queue, consumer); err != nil {
		log.Fatal(err)
	}

	// Start sending movies to NSQ
	go func() {
		for range time.Tick(500 * time.Millisecond) {
			// without generating: evpb.Send(queue, &pb.Movie{...})
			// would need to lookup topic name by reflected type
			if err := pb.SendMovie(queue, &pb.Movie{
				Title: "Pulp Fiction",
				Year:  1994,
			}); err != nil {
				log.Fatal(err)
			}
		}
	}()

	// Wait for interrupt signal
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
}
```

## queues

  * [nsq](http://nsq.io/)
