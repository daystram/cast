package handlers

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"gitlab.com/daystram/cast/cast-be/constants"
	"gitlab.com/daystram/cast/cast-be/datatransfers"

	"cloud.google.com/go/pubsub"
)

func (m *module) TranscodeListenerWorker() {
	var mutex sync.Mutex
	for {
		fmt.Println("[TranscodeListenerWorker] TranscodeListenerWorker started")
		_ = m.mq.completeSubscription.Receive(context.Background(), func(ctx context.Context, msg *pubsub.Message) {
			msg.Ack()
			mutex.Lock()
			defer mutex.Unlock()
			hash := strings.Split(string(msg.Data), ":")[0]
			resolution, err := strconv.Atoi(strings.Split(string(msg.Data), ":")[1])
			if err != nil {
				fmt.Println("[TranscodeListenerWorker] Failed parsing message from transcoder")
				return
			}
			var video datatransfers.Video
			if video, err = m.db.videoOrm.GetOneByHash(hash); err != nil {
				fmt.Printf("[TranscodeListenerWorker] Unknown video with hash %s\n", hash)
				return
			}
			if err = m.db.videoOrm.SetResolution(hash, resolution); err != nil {
				fmt.Printf("[TranscodeListenerWorker] Failed updating video %s resolution\n", hash)
				return
			}
			if resolution >= 1 {
				m.PushNotification(video.Author.ID, datatransfers.NotificationOutgoing{
					Message:   fmt.Sprintf("%s is now ready in %s!", video.Title, constants.VideoResolutions[resolution]),
					Username:  video.Author.Username,
					Hash:      video.Hash,
					CreatedAt: time.Now(),
				})
			}
			if resolution == 1 {
				m.BroadcastNotificationSubscriber(video.Author.ID, datatransfers.NotificationOutgoing{
					Message:   fmt.Sprintf("%s just uploaded %s! Watch now!", video.Author.Name, video.Title),
					Username:  video.Author.Username,
					Hash:      video.Hash,
					CreatedAt: time.Now(),
				})
			}
			fmt.Printf("[TranscodeListenerWorker] Done transcoding %s to %s\n", hash, constants.Resolutions[resolution])
		})
	}
}

func (m *module) StartTranscode(hash string) {
	m.mq.transcodeTopic.Publish(context.Background(), &pubsub.Message{Data: []byte(hash)})
	fmt.Printf("[StartTranscode] Transcoding task for %s commencing\n", hash)
}
