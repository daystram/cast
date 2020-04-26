package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
	"gopkg.in/ini.v1"
)

type Config struct {
	ProjectID             string
	APIKey                string
	TopicTranscode        string
	TopicComplete         string
	SubscriptionTranscode string
	UploadsDir            string
	FFMpegExecutable      string
	MP4BoxExecutable      string
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
	// Init Configuration
	configFile, err := ini.Load("app.conf")
	if err != nil {
		log.Fatalf("[Initialize] unable to find app.conf. %+v\n", err)
	}
	config := Config{
		ProjectID:             configFile.Section("").Key("google_project_id").String(),
		APIKey:                configFile.Section("").Key("google_api_key").String(),
		TopicComplete:         configFile.Section("").Key("pubsub_topic_complete").String(),
		SubscriptionTranscode: configFile.Section("").Key("pubsub_subscription_transcode").String(),
		UploadsDir:            configFile.Section("").Key("uploads_dir").String(),
		FFMpegExecutable:      configFile.Section("").Key("ffmpeg_executable").String(),
		MP4BoxExecutable:      configFile.Section("").Key("mp4box_executable").String(),
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

	// Init Google PubSub
	pubsubClient, err := pubsub.NewClient(context.Background(), config.ProjectID, option.WithCredentialsFile(config.APIKey))
	if err != nil {
		log.Fatalf("Failed connecting to Google PubSub. %+v\n", err)
	}
	fmt.Println("[Initialization] Google PubSub connected")

	// Listen to messages
	var mu sync.Mutex
	sub := pubsubClient.Subscription(config.SubscriptionTranscode)
	topic := pubsubClient.Topic(config.TopicComplete)
	fmt.Println("[cast-is] Ready to transcode")
	for {
		_ = sub.Receive(context.Background(), func(ctx context.Context, msg *pubsub.Message) {
			msg.Ack()
			mu.Lock()
			defer mu.Unlock()
			hash := string(msg.Data)
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
				topic.Publish(context.Background(), &pubsub.Message{Data: []byte(fmt.Sprintf("%s:%d", hash, i))})
			}
			cleanUp(workDir)
			fmt.Printf("[cast-is] Done: %s\n", hash)
			fmt.Fprintf(outfile, "[cast-is] ----------------------- Done: %s\n", hash)
		})
	}
}

func cleanUp(path string) {
	files, _ := filepath.Glob(fmt.Sprintf("%s/temp*", path))
	for _, file := range files {
		_ = os.Remove(file)
	}
}
