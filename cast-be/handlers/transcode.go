package handlers

import (
	"context"
	"fmt"
	"gitlab.com/daystram/cast/cast-be/constants"
	"strconv"
	"strings"
	"sync"

	"cloud.google.com/go/pubsub"
)

func (m *module) TranscodeListenerWorker() {
	var mutex sync.Mutex
	for {
		fmt.Println("[TranscodeListenerWorker] TranscodeListenerWorker started")
		_ = m.mq().completeSubscription.Receive(context.Background(), func(ctx context.Context, msg *pubsub.Message) {
			msg.Ack()
			mutex.Lock()
			defer mutex.Unlock()
			hash := strings.Split(string(msg.Data), ":")[0]
			resolution, err := strconv.Atoi(strings.Split(string(msg.Data), ":")[1])
			if err != nil {
				fmt.Println("[TranscodeListenerWorker] Failed parsing message from transcoder")
				return
			}
			if err = m.db().videoOrm.SetResolution(hash, resolution); err != nil {
				fmt.Printf("[TranscodeListenerWorker] Failed updating video %s resolution\n", hash)
				return
			}
			fmt.Printf("[TranscodeListenerWorker] Done transcoding %s to %s\n", hash, constants.Resolutions[resolution])
		})
	}
}

func (m *module) StartTranscode(hash string) {
	m.mq().transcodeTopic.Publish(context.Background(), &pubsub.Message{Data: []byte(hash)})
	fmt.Printf("[StartTranscode] Transcoding task for %s commencing\n", hash)
}
