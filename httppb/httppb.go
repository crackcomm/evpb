package httppb

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/crackcomm/evpb"
	"github.com/hailocab/go-hostpool"
)

// New - Creates new http evpb.Interface.
func New(opts ...Option) evpb.Interface {
	q := &queue{
		basePath:    "/api/v1/httppb",
		serverMux:   http.NewServeMux(),
		concurrency: 1,
		msgQueue:    make(chan *message, 10000),
		handlers:    make(map[string]evpb.Consumer),
	}
	for _, opt := range opts {
		opt(q)
	}
	return q
}

type queue struct {
	hostPool    hostpool.HostPool
	basePath    string
	topicPrefix string
	concurrency int

	msgQueue chan *message
	handlers map[string]evpb.Consumer

	listenOnce sync.Once
	listenAddr string
	httpListen net.Listener
	httpServer *http.Server
	serverMux  *http.ServeMux
}

type message struct {
	topic string
	body  []byte
}

func (q *queue) Send(topic string, body []byte) (err error) {
	host := q.hostPool.Get()
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s%s", host.Host(), q.path(topic)), bytes.NewReader(body))
	if err != nil {
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		host.Mark(err)
		return
	}
	// Close response body
	resp.Body.Close()
	return
}

func (q *queue) Consume(topic string, consume evpb.Consumer) (err error) {
	q.handlers[topic] = consume
	q.serverMux.HandleFunc(q.path(topic), func(w http.ResponseWriter, r *http.Request) {
		// Read request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Request body reading error: %v", err)
			return
		}
		q.msgQueue <- &message{
			topic: topic,
			body:  body,
		}
	})
	q.listenOnce.Do(func() {
		q.httpListen, err = net.Listen("tcp", q.listenAddr)
		if err != nil {
			return
		}
		q.httpServer = &http.Server{
			Addr:         q.listenAddr,
			Handler:      q.serverMux,
			ReadTimeout:  time.Minute,
			WriteTimeout: time.Minute,
		}
		go func() {
			if err := q.httpServer.Serve(q.httpListen); err != nil {
				panic(err)
			}
		}()
		for index := 0; index < q.concurrency; index++ {
			go func() {
				for msg := range q.msgQueue {
					handler := q.handlers[msg.topic]
					err := handler(msg.body)
					if err != nil {
						log.Printf("Consumer error: (topic: %q) %v", msg.topic, err)
					}
				}
			}()
		}
	})
	return
}

func (q *queue) Stop() (err error) {
	return q.httpListen.Close()
}

func (q *queue) path(t string) string {
	return path.Join(q.basePath, q.topic(t))
}

func (q *queue) topic(t string) string {
	if q.topicPrefix == "" {
		return t
	}
	return strings.Join([]string{q.topicPrefix, t}, "")
}
