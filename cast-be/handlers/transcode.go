package handlers

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/streadway/amqp"

	conf "github.com/daystram/cast/cast-be/config"
	"github.com/daystram/cast/cast-be/constants"
	"github.com/daystram/cast/cast-be/datatransfers"
)

func (m *module) TranscodeListenerWorker() {
	var err error
	var receiver <-chan amqp.Delivery
	if receiver, err = m.mq.Consume(
		conf.AppConfig.RabbitMQQueueProgress,
		"",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		log.Printf("[TranscodeListenerWorker] Failed creating RabbitMQ consumer. %+v\n", err)
	}
	log.Println("[TranscodeListenerWorker] TranscodeListenerWorker started")
	for msg := range receiver {
		hash := strings.Split(string(msg.Body), ":")[0]
		resolution, err := strconv.Atoi(strings.Split(string(msg.Body), ":")[1])
		if err != nil {
			log.Println("[TranscodeListenerWorker] Failed parsing message from transcoder")
			continue
		}
		var video datatransfers.Video
		if video, err = m.db.videoOrm.GetOneByHash(hash); err != nil {
			log.Printf("[TranscodeListenerWorker] Unknown video with hash %s\n", hash)
			continue
		}
		if err = m.db.videoOrm.SetResolution(hash, resolution); err != nil {
			log.Printf("[TranscodeListenerWorker] Failed updating video %s resolution\n", hash)
			continue
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
				Message:   fmt.Sprintf("%s just uploaded %s! Watch now!", video.Author.Username, video.Title),
				Username:  video.Author.Username,
				Hash:      video.Hash,
				CreatedAt: time.Now(),
			})
		}
		log.Printf("[TranscodeListenerWorker] Done transcoding %s to %s\n", hash, constants.Resolutions[resolution])
	}
}

func (m *module) StartTranscode(hash string) {
	_ = m.mq.Publish("", conf.AppConfig.RabbitMQQueueTask, true, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(hash),
		})
	log.Printf("[StartTranscode] Transcoding task for %s commencing\n", hash)
}
