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
	if err := pb.ConsumeMovie(queue, consumer); err != nil {
		log.Fatal(err)
	}

	// Start sending movies to NSQ
	go func() {
		for range time.Tick(500 * time.Millisecond) {
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
