package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"

	"cloud.google.com/go/pubsub"
	"github.com/jhw0604/StreamLog/handle"
)

//Errors
var (
	ErrNotExistsTopic = errors.New("not exists topic")
)

//configs
var (
	address   = flag.String("addr", "127.0.0.1:8080", "server binding ip:port")
	projectID = flag.String("project", "", "pubsub project id")
	topicID   = flag.String("topic", "", "pubsub topic id")
)

func main() {
	flag.Parse()

	ctx, cancle := context.WithCancel(context.Background())
	client, err := pubsub.NewClient(ctx, *projectID)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	topic := client.Topic(*topicID)
	defer topic.Stop()

	exists, err := topic.Exists(ctx)
	if err != nil {
		panic(err)
	}
	if !exists {
		panic(ErrNotExistsTopic)
	}

	svr := http.Server{
		Addr: *address,
		Handler: handle.LogAPI{
			Topic:   topic,
			Context: ctx,
			Cancle:  cancle,
			CancleCondition: func(r *http.Request) bool {
				if r.URL.Path == "/ServerCommand/Shutdown" {
					return true
				}
				return false
			},
		},
	}
	go func() {
		err = svr.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	log.Println("Server Started.")
	<-ctx.Done()
	log.Println("Server stoped.")
}
