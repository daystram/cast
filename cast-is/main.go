package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

var config Config

type Config struct {
	TempDir string

	RabbitMQURI           string
	RabbitMQQueueTask     string
	RabbitMQQueueProgress string

	S3URI       string
	S3Bucket    string
	S3Region    string
	S3AccessKey string
	S3SecretKey string

	FFMpegExecutable string
	MP4BoxExecutable string
}

var module Module

type Module struct {
	mqPub *amqp.Channel
	mqSub *amqp.Channel
	s3    *s3.S3
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

var resolutions = []Resolution{
	{"Audio", FlagsAudio},
	{"240p", Flags240},
	{"360p", Flags360},
	{"480p", Flags480},
	{"720p", Flags720},
	{"1080p", Flags1080},
}

func init() {
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

	config = Config{
		TempDir:               viper.GetString("TEMP_DIR"),
		RabbitMQURI:           viper.GetString("RABBITMQ_URI"),
		RabbitMQQueueTask:     viper.GetString("RABBITMQ_QUEUE_TASK"),
		RabbitMQQueueProgress: viper.GetString("RABBITMQ_QUEUE_PROGRESS"),
		S3URI:                 viper.GetString("S3_URI"),
		S3Bucket:              viper.GetString("S3_BUCKET"),
		S3Region:              viper.GetString("S3_REGION"),
		S3AccessKey:           viper.GetString("S3_ACCESS_KEY"),
		S3SecretKey:           viper.GetString("S3_SECRET_KEY"),
		FFMpegExecutable:      viper.GetString("FFMPEG_EXECUTABLE"),
		MP4BoxExecutable:      viper.GetString("MP4BOX_EXECUTABLE"),
	}

	// Init RabbitMQ
	var mqConn *amqp.Connection
	if mqConn, err = amqp.Dial(config.RabbitMQURI); err != nil {
		log.Fatalf("[Initialization] Failed connecting to RabbitMQ at %s. %+v\n", config.RabbitMQURI, err)
	}
	var mqInitCh *amqp.Channel
	if mqInitCh, err = mqConn.Channel(); err != nil {
		log.Fatalf("[Initialization] Failed opening RabbitMQ init channel. %+v\n", err)
	}
	if _, err = mqInitCh.QueueDeclare(config.RabbitMQQueueTask, true, false, false, false, nil); err != nil {
		log.Fatalf("[Initialization] Failed declaring RabbitMQ queue %s. %+v\n", config.RabbitMQQueueTask, err)
	}
	if _, err = mqInitCh.QueueDeclare(config.RabbitMQQueueProgress, true, false, false, false, nil); err != nil {
		log.Fatalf("[Initialization] Failed declaring RabbitMQ queue %s. %+v\n", config.RabbitMQQueueProgress, err)
	}
	log.Printf("[Initialization] Successfully connected to RabbitMQ!\n")

	// Init S3
	var s3Session *session.Session
	if s3Session, err = session.NewSession(&aws.Config{
		Endpoint:         aws.String(config.S3URI),
		Region:           aws.String(config.S3Region),
		Credentials:      credentials.NewStaticCredentials(config.S3AccessKey, config.S3SecretKey, ""),
		DisableSSL:       aws.Bool(false),
		S3ForcePathStyle: aws.Bool(true),
	}); err != nil {
		log.Fatalf("[Initialization] Failed creating S3 session to %s. %+v\n", config.S3URI, err)
	}
	s3Client := s3.New(s3Session)
	if _, err := s3Client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(config.S3Bucket),
	}); err != nil {
		log.Fatalf("[Initialization] Failed connecting to S3. %+v\n", err)
	}
	log.Printf("[Initialization] Successfully connected to S3!\n")

	module = Module{mq: mq}
}

	module = Module{mqPub: mqPubCh, mqSub: mqSubCh, s3: s3Client}
}

func main() {
	// Listen to messages
	var err error
	var receiver <-chan amqp.Delivery
	if receiver, err = module.mq.Consume(
		config.RabbitMQQueueTask,
		"",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		log.Printf("[cast-is] Failed creating RabbitMQ consumer. %+v\n", err)
	}
	fmt.Println("[cast-is] Ready to transcode!")
	for msg := range receiver {
		hash := string(msg.Body)
		workDir := fmt.Sprintf("%s/%s", config.TempDir, hash)
		_ = os.MkdirAll(workDir, 0755)
		outfile, _ := os.Create(fmt.Sprintf("%s/transcode.log", workDir))
		// TODO: get raw from S3

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
				fmt.Fprintf(outfile, "[cast-is] ----------------------- Cancelled\n")
				continue
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
			_ = module.mq.Publish("", config.RabbitMQQueueProgress, true, false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(fmt.Sprintf("%s:%d", hash, i)),
				})
			// TODO: write into S3
			fmt.Fprintf(outfile, "[cast-is] ----------------------- Completed\n")
		}
		cleanUp(workDir)
		fmt.Printf("[cast-is] Done: %s\n", hash)
		fmt.Fprintf(outfile, "[cast-is] ----------------------- Done: %s\n", hash)
		outfile.Close()
	}
}

func cleanUp(path string) {
	files, _ := filepath.Glob(fmt.Sprintf("%s/temp*", path))
	for _, file := range files {
		_ = os.Remove(file)
	}
}
