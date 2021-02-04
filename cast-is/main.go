package main

import (
	"fmt"
	"io"
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
	UseCUDA bool

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
	FlagsAudio   = "-i video.mp4 -vn -c:a aac -b 128k -dash 1 -y audio.m4a"
	BaseVideoCPU = "-i video.mp4 -c:v libx264 -vsync passthrough -x264-params keyint=25:min-keyint=25:no-scenecut -movflags +faststart -y"
	BaseVideoGPU = "-hwaccel cuvid -i video.mp4 -c:v h264_nvenc -vsync passthrough -x264-params keyint=25:min-keyint=25:no-scenecut -movflags +faststart -y"
	Flags240     = "-an -vf scale=-2:240 -b:v 400k temp_240.mp4"
	Flags360     = "-an -vf scale=-2:360 -b:v 800k temp_360.mp4"
	Flags480     = "-an -vf scale=-2:480 -b:v 1200k temp_480.mp4"
	Flags720     = "-an -vf scale=-2:720 -b:v 2400k temp_720.mp4"
	Flags1080    = "-an -vf scale=-2:1080 -b:v 4800k temp_1080.mp4"
	FlagsDASH    = "-dash 10000 -rap -frag-rap -bs-switching no -url-template -dash-profile onDemand -segment-name segment_$RepresentationID$ -out manifest.mpd"
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
		UseCUDA:               viper.GetBool("USE_CUDA"),
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
		fmt.Printf("[Initialization] Failed connecting to RabbitMQ at %s. %+v\n", config.RabbitMQURI, err)
		panic(err)
	}
	var mqInitCh *amqp.Channel
	if mqInitCh, err = mqConn.Channel(); err != nil {
		fmt.Printf("[Initialization] Failed opening RabbitMQ init channel. %+v\n", err)
		panic(err)
	}
	if _, err = mqInitCh.QueueDeclare(config.RabbitMQQueueTask, true, false, false, false, nil); err != nil {
		fmt.Printf("[Initialization] Failed declaring RabbitMQ queue %s. %+v\n", config.RabbitMQQueueTask, err)
		panic(err)
	}
	if _, err = mqInitCh.QueueDeclare(config.RabbitMQQueueProgress, true, false, false, false, nil); err != nil {
		fmt.Printf("[Initialization] Failed declaring RabbitMQ queue %s. %+v\n", config.RabbitMQQueueProgress, err)
		panic(err)
	}
	var mqPubCh *amqp.Channel
	var mqSubCh *amqp.Channel
	if mqPubCh, err = mqConn.Channel(); err != nil {
		fmt.Printf("[Initialization] Failed opening RabbitMQ publisher channel. %+v\n", err)
		panic(err)
	}
	if mqSubCh, err = mqConn.Channel(); err != nil {
		fmt.Printf("[Initialization] Failed opening RabbitMQ subscription channel. %+v\n", err)
		panic(err)
	}
	fmt.Printf("[Initialization] Successfully connected to RabbitMQ!\n")

	// Init S3
	var s3Session *session.Session
	if s3Session, err = session.NewSession(&aws.Config{
		Endpoint:         aws.String(config.S3URI),
		Region:           aws.String(config.S3Region),
		Credentials:      credentials.NewStaticCredentials(config.S3AccessKey, config.S3SecretKey, ""),
		DisableSSL:       aws.Bool(false),
		S3ForcePathStyle: aws.Bool(true),
	}); err != nil {
		fmt.Printf("[Initialization] Failed creating S3 session to %s. %+v\n", config.S3URI, err)
		panic(err)
	}
	s3Client := s3.New(s3Session)
	if _, err := s3Client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(config.S3Bucket),
	}); err != nil {
		fmt.Printf("[Initialization] Failed connecting to S3. %+v\n", err)
		panic(err)
	}
	fmt.Printf("[Initialization] Successfully connected to S3!\n")

	// Check CUDA
	if config.UseCUDA {
		reader, writer := io.Pipe()
		cmd1 := exec.Command("ldconfig", "-p")
		cmd2 := exec.Command("grep", "nvcuvid")
		cmd1.Stdout = writer
		cmd2.Stdin = reader
		fmt.Printf("[Initialization] Checking nvcuvid CUDA driver availability... ")
		if err = cmd1.Start(); err != nil {
			fmt.Printf("ERROR: %+v\n", err)
			panic(err)
		}
		if err = cmd2.Start(); err != nil {
			fmt.Printf("ERROR: %+v\n", err)
			panic(err)
		}
		if err = cmd1.Wait(); err != nil {
			fmt.Printf("ERROR: %+v\n", err)
			panic(err)
		}
		writer.Close()
		if err = cmd2.Wait(); err != nil {
			fmt.Printf("ERROR: nvcuvid driver not found! %+v\n", err)
			panic(err)
		}
		reader.Close()
		fmt.Printf("OK\n")
	}

	module = Module{mqPub: mqPubCh, mqSub: mqSubCh, s3: s3Client}
}

func main() {
	// Listen to transcode tasks
	var err error
	var receiver <-chan amqp.Delivery
	if receiver, err = module.mqSub.Consume(
		config.RabbitMQQueueTask,
		"",
		false,
		false,
		false,
		false,
		nil,
	); err != nil {
		fmt.Printf("[cast-is] Failed creating RabbitMQ consumer. %+v\n", err)
	}
	forever := make(chan bool)
	go func() {
		for msg := range receiver {
			// Init working directory
			startAll := time.Now().UnixNano()
			hash := string(msg.Body)
			workDir := fmt.Sprintf("%s/%s", config.TempDir, hash)
			_ = os.MkdirAll(workDir, 0755)
			logFile, _ := os.Create(fmt.Sprintf("%s/transcode.log", workDir))
			fmt.Printf("[cast-is] Start: %s\n", hash)
			fmt.Fprintf(logFile, "\n[cast-is] ----------------------- Start: %s", hash)

			// Fetch assets
			fmt.Printf("[cast-is] Retrieving assets... ")
			fmt.Fprintf(logFile, "\n[cast-is] ----------------------- Retrieving assets... ")
			start := time.Now().UnixNano()
			if object, err := module.s3.GetObject(&s3.GetObjectInput{
				Bucket: aws.String(config.S3Bucket),
				Key:    aws.String(fmt.Sprintf("video/%s/video.mp4", hash)),
			}); err != nil {
				fmt.Printf("ERROR: %v\n", err)
				fmt.Fprintf(logFile, "ERROR\n")
				logFile.Close()
				msg.Nack(false, false)
				continue
			} else {
				video, _ := os.Create(fmt.Sprintf("%s/video.mp4", workDir))
				_, _ = io.Copy(video, object.Body)
				video.Close()
				object.Body.Close()
				fmt.Printf("OK (%.2fs)\n", float64(time.Now().UnixNano()-start)/1e9)
				fmt.Fprintf(logFile, "OK (%.2fs)\n", float64(time.Now().UnixNano()-start)/1e9)
			}

			// Begin transcoding
			for i, resolution := range resolutions {
				// Transcode to resolution
				fmt.Printf("[cast-is] Transcode -> %s ", resolution.Name)
				fmt.Fprintf(logFile, "\n[cast-is] ----------------------- Transcode -> %s\n", resolution.Name)
				flags := strings.Split(resolution.Flags, " ")
				if resolution.Name != "Audio" {
					if config.UseCUDA {
						flags = append(strings.Split(BaseVideoGPU, " "), flags...) // prepend
					} else {
						flags = append(strings.Split(BaseVideoCPU, " "), flags...) // prepend
					}
				}
				cmd := exec.Command(config.FFMpegExecutable, flags...)
				cmd.Stderr = logFile
				cmd.Dir = workDir
				start = time.Now().UnixNano()
				if err := cmd.Run(); err != nil {
					fmt.Printf("ERROR: %v\n", err)
					fmt.Fprintf(logFile, "[cast-is] ----------------------- ERROR\n")
					continue
				}
				fmt.Printf("DONE (%.2fs) ", float64(time.Now().UnixNano()-start)/1e9)
				fmt.Fprintf(logFile, "[cast-is] ----------------------- DONE (%.2fs)\n", float64(time.Now().UnixNano()-start)/1e9)
				if resolution.Name == "Audio" {
					// Upload audio
					start = time.Now().UnixNano()
					file, _ := os.Open(fmt.Sprintf("%s/audio.m4a", workDir))
					_, _ = module.s3.PutObject(&s3.PutObjectInput{
						Bucket:      aws.String(config.S3Bucket),
						Key:         aws.String(fmt.Sprintf("video/%s/audio.m4a", hash)),
						Body:        file,
						ContentType: aws.String("audio/m4a"),
					})
					fmt.Printf("; UPLOADED (%.2fs)\n", float64(time.Now().UnixNano()-start)/1e9)
					fmt.Fprintf(logFile, "[cast-is] ----------------------- UPLOADED (%.2fs)\n", float64(time.Now().UnixNano()-start)/1e9)
					continue
				}

				// Generate DASH manifest
				fmt.Fprintf(logFile, "[cast-is] ----------------------- DASH -> %s\n", resolution.Name)
				cmd = exec.Command(config.MP4BoxExecutable, append(strings.Split(FlagsDASH, " "), TempFileNames[5-i:]...)...)
				cmd.Stderr = logFile
				cmd.Dir = workDir
				start = time.Now().UnixNano()
				if err := cmd.Run(); err != nil {
					fmt.Printf("; ERROR: %v\n", err)
					fmt.Fprintf(logFile, "[cast-is] ----------------------- ERROR\n")
					continue
				}
				fmt.Printf("; DASH (%.2fs) ", float64(time.Now().UnixNano()-start)/1e9)
				fmt.Fprintf(logFile, "\n[cast-is] ----------------------- DONE (%.2fs)\n", float64(time.Now().UnixNano()-start)/1e9)

				// Notify cast-be
				_ = module.mqPub.Publish("", config.RabbitMQQueueProgress, true, false,
					amqp.Publishing{
						ContentType: "text/plain",
						Body:        []byte(fmt.Sprintf("%s:%d", hash, i)),
					})

				// Upload to S3
				files, _ := filepath.Glob(fmt.Sprintf("%s/segment_*", workDir))
				files = append(files, fmt.Sprintf("%s/manifest.mpd", workDir))
				start = time.Now().UnixNano()
				for _, path := range files {
					file, _ := os.Open(path)
					_, fileName := filepath.Split(path)
					var mime string
					if fileName == "manifest.mpd" {
						mime = "application/dash+xml"
					} else {
						mime = "video/mp4"
					}
					_, _ = module.s3.PutObject(&s3.PutObjectInput{
						Bucket:      aws.String(config.S3Bucket),
						Key:         aws.String(fmt.Sprintf("video/%s/%s", hash, fileName)),
						Body:        file,
						ContentType: aws.String(mime),
					})
				}
				fmt.Printf("; UPLOADED (%.2fs)\n", float64(time.Now().UnixNano()-start)/1e9)
				fmt.Fprintf(logFile, "[cast-is] ----------------------- UPLOADED (%.2fs)\n", float64(time.Now().UnixNano()-start)/1e9)
				fmt.Fprintf(logFile, "[cast-is] ----------------------- Completed -> %s\n", resolution.Name)
			}
			// Upload log
			logFile.Close()
			logFile, _ = os.Open(logFile.Name())
			fmt.Fprintf(logFile, "[cast-is] ----------------------- Done: %s (%.2fs)\n", hash, float64(time.Now().UnixNano()-startAll)/1e9)
			_, _ = module.s3.PutObject(&s3.PutObjectInput{
				Bucket:      aws.String(config.S3Bucket),
				Key:         aws.String(fmt.Sprintf("video/%s/transcode.log", hash)),
				Body:        logFile,
				ContentType: aws.String("text/plain"),
			})
			logFile.Close()
			fmt.Printf("[cast-is] Done: %s (%.2fs)\n", hash, float64(time.Now().UnixNano()-startAll)/1e9)
			fmt.Println("[cast-is] Ready!")

			// Cleanup
			cleanUp(workDir)
			msg.Ack(false)
		}
	}()
	fmt.Println("[cast-is] Ready!")
	<-forever
}

func cleanUp(path string) {
	files, _ := filepath.Glob(fmt.Sprintf("%s/*", path))
	for _, file := range files {
		_ = os.Remove(file)
	}
	_ = os.Remove(path)
}
