package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/jhw0604/StreamLog/handle"
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
		log.Println("pub/sub client can't start", err)
		os.Exit(1)
	}
	defer client.Close()

	topic := client.Topic(*topicID)
	defer topic.Stop()

	exists, err := topic.Exists(ctx)
	if err != nil {
		log.Println("topic exists chacke fail:", err)
		os.Exit(1)
	}
	if !exists {
		log.Println("not exists topic")
		os.Exit(1)
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
			log.Println("server brocked:", err)
			os.Exit(1)
		}
	}()

	log.Println("server started")
	<-ctx.Done()
	log.Println("server stoped")
}
