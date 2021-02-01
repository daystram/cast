package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

type Config struct {
	UploadsDir string

	RabbitMQURI           string
	RabbitMQQueueTask     string
	RabbitMQQueueProgress string

	FFMpegExecutable string
	MP4BoxExecutable string
}

type Resolution struct {
	Name  string
	Flags string
}

var (
	TempFileNames = []string{"temp_1080.mp4", "temp_720.mp4", "temp_480.mp4", "temp_360.mp4", "temp_240.mp4", "audio.m4a"}
)

const (
	FlagsAudio = "-i video.mp4 -vn -acodec aac -ab 128k -dash 1 -y audio.m4a"
	BaseVideo  = "-i video.mp4 -vsync passthrough -c:v libx264 -x264-params keyint=25:min-keyint=25:no-scenecut -movflags +faststart -y "
	Flags240   = BaseVideo + "-an -vf scale=-2:240 -b:v 400k -preset faster temp_240.mp4"
	Flags360   = BaseVideo + "-an -vf scale=-2:360 -b:v 800k -preset faster temp_360.mp4"
	Flags480   = BaseVideo + "-an -vf scale=-2:480 -b:v 1200k -preset faster temp_480.mp4"
	Flags720   = BaseVideo + "-an -vf scale=-2:720 -b:v 2400k -preset faster temp_720.mp4"
	Flags1080  = BaseVideo + "-an -vf scale=-2:1080 -b:v 4800k -preset faster temp_1080.mp4"
	FlagsDASH  = "-dash 10000 -rap -frag-rap -bs-switching no -url-template -dash-profile onDemand -segment-name segment_$RepresentationID$ -out manifest.mpd"
)

func main() {
	var err error

	// Init Configuration
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/")
	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv()
	_ = viper.ReadInConfig()

	config := Config{
		UploadsDir:            viper.GetString("UPLOADS_DIR"),
		RabbitMQURI:           viper.GetString("RABBITMQ_URI"),
		RabbitMQQueueTask:     viper.GetString("RABBITMQ_QUEUE_TASK"),
		RabbitMQQueueProgress: viper.GetString("RABBITMQ_QUEUE_PROGRESS"),
		FFMpegExecutable:      viper.GetString("FFMPEG_EXECUTABLE"),
		MP4BoxExecutable:      viper.GetString("MP4BOX_EXECUTABLE"),
	}
	resolutions := []Resolution{
		{"Audio", FlagsAudio},
		{"240p", Flags240},
		{"360p", Flags360},
		{"480p", Flags480},
		{"720p", Flags720},
		{"1080p", Flags1080},
	}
	fmt.Println("[Initialization] config loaded")

	// Init RabbitMQ
	var mqConn *amqp.Connection
	if mqConn, err = amqp.Dial(config.RabbitMQURI); err != nil {
		log.Fatalf("[Initialization] Failed connecting to RabbitMQ at %s. %+v\n", config.RabbitMQURI, err)
	}
	var mq *amqp.Channel
	if mq, err = mqConn.Channel(); err != nil {
		log.Fatalf("[Initialization] Failed opening RabbitMQ channel. %+v\n", err)
	}
	if _, err = mq.QueueDeclare(config.RabbitMQQueueTask, true, false, false, false, nil); err != nil {
		log.Fatalf("[Initialization] Failed declaring RabbitMQ queue %s. %+v\n", config.RabbitMQQueueTask, err)
	}
	if _, err = mq.QueueDeclare(config.RabbitMQQueueProgress, true, false, false, false, nil); err != nil {
		log.Fatalf("[Initialization] Failed declaring RabbitMQ queue %s. %+v\n", config.RabbitMQQueueProgress, err)
	}
	log.Printf("[Initialization] Successfully connected to RabbitMQ!\n")

	// Listen to messages
	var receiver <-chan amqp.Delivery
	if receiver, err = mq.Consume(
		config.RabbitMQQueueTask,
		"",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		log.Printf("[TranscodeListenerWorker] Failed creating RabbitMQ consumer. %+v\n", err)
	}
	fmt.Println("[cast-is] Ready to transcode!")
	for msg := range receiver {
		hash := string(msg.Body)
		workDir := fmt.Sprintf("%s/%s", config.UploadsDir, hash)
		outfile, _ := os.Create(fmt.Sprintf("%s/transcode.log", workDir))
		defer outfile.Close()

		fmt.Printf("[cast-is] Start: %s\n", hash)
		fmt.Fprintf(outfile, "\n[cast-is] ----------------------- Start: %s", hash)
		for i, resolution := range resolutions {
			fmt.Printf("[cast-is] %s -> %s\n", hash, resolution.Name)
			fmt.Fprintf(outfile, "\n[cast-is] ----------------------- %s -> %s\n", hash, resolution.Name)
			cmd := exec.Command(config.FFMpegExecutable, strings.Split(resolution.Flags, " ")...)
			cmd.Stderr = outfile
			cmd.Dir = workDir
			if err := cmd.Run(); err != nil {
				fmt.Printf("[cast-is] Stopped: %s\n", hash)
				fmt.Fprintf(outfile, "[cast-is] ----------------------- Stopped: %s\n", hash)
				return
			}
			if resolution.Name == "Audio" {
				continue
			}
			time.Sleep(1 * time.Second)
			_ = os.Remove(fmt.Sprintf("%s/manifest.mpd", workDir))
			time.Sleep(2 * time.Second)
			fmt.Fprintf(outfile, "[cast-is] ----------------------- %s -> %s DASH\n", hash, resolution.Name)
			cmd = exec.Command(config.MP4BoxExecutable, append(strings.Split(FlagsDASH, " "), TempFileNames[5-i:]...)...)
			cmd.Stderr = outfile
			cmd.Dir = workDir
			if err := cmd.Run(); err != nil {
				fmt.Println(err)
			}
			_ = mq.Publish("", config.RabbitMQQueueProgress, true, false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(fmt.Sprintf("%s:%d", hash, i)),
				})
		}
		cleanUp(workDir)
		fmt.Printf("[cast-is] Done: %s\n", hash)
		fmt.Fprintf(outfile, "[cast-is] ----------------------- Done: %s\n", hash)
	}
}

func cleanUp(path string) {
	files, _ := filepath.Glob(fmt.Sprintf("%s/temp*", path))
	for _, file := range files {
		_ = os.Remove(file)
	}
}
