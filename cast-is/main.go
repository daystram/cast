package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"google.golang.org/api/option"
	"gopkg.in/ini.v1"
	"log"
	"sync"
	"time"
)

func main() {
	// Init Configuration
	config, err := ini.Load("app.conf")
	if err != nil {
		log.Fatalf("[Initialize] unable to find app.conf. %+v\n", err)
	}
	projectID := config.Section("").Key("google_project_id").String()
	apiKey := config.Section("").Key("google_api_key").String()
	transcodeTopic := config.Section("").Key("pubsub_subscription_transcode").String()

	// Init Google PubSub
	pubsubClient, err := pubsub.NewClient(context.Background(), projectID,
		option.WithCredentialsFile(apiKey))
	if err != nil {
		log.Fatalf("Failed connecting to Google PubSub. %+v\n", err)
	}

	// Listen to messages
	var mu sync.Mutex
	sub := pubsubClient.Subscription(transcodeTopic)
	for {
		fmt.Println("READY")
		err = sub.Receive(context.Background(), func(ctx context.Context, msg *pubsub.Message) {
			msg.Ack()
			mu.Lock()
			defer mu.Unlock()
			fmt.Printf("Got message: %q\n", string(msg.Data))
			time.Sleep(5*time.Second)
		})
	}
}
