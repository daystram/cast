package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"google.golang.org/api/option"
	"log"
	"sync"
)

func main() {
	var mu sync.Mutex
	// Init Google PubSub
	pubsubClient, err := pubsub.NewClient(context.Background(), "daystram-cast", option.WithCredentialsFile("c:/key.json"))
	if err != nil {
		log.Fatalf("Failed connecting to Google PubSub. %+v\n", err)
	}
	sub := pubsubClient.Subscription("cast-is")
	for {
		fmt.Println("READY")
		err = sub.Receive(context.Background(), func(ctx context.Context, msg *pubsub.Message) {
			fmt.Printf("Got message: %q\n", string(msg.Data))
			msg.Ack()
			mu.Lock()
			defer mu.Unlock()
		})
		//fmt.Println(err)
	}
}
