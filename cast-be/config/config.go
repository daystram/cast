package config

import (
	"github.com/astaxie/beego"
)

var AppConfig Config

type Config struct {
	Domain    string
	JWTSecret string
	Debug     bool

	RTMPPort int

	MongoDBURI  string
	MongoDBName string

	GoogleProjectID          string
	JSONKey                  string
	TopicNameTranscode       string
	SubscriptionNameComplete string

	UploadsDirectory string

	MailgunDomain string
	MailgunAPIKey string
	Hostname      string
}

func InitializeAppConfig() {
	_ = beego.LoadAppConfig("ini", "config/app.conf")
	var err error

	// Domain
	if AppConfig.Domain = beego.AppConfig.String("domain"); AppConfig.Domain == "" {
		panic("domain in app.conf is missing or is in wrong format")
	}

	// Debug
	if AppConfig.Debug, err = beego.AppConfig.Bool("debug"); err != nil {
		panic("debug is missing in app.conf")
	}

	// RTMP Port
	if AppConfig.RTMPPort, err = beego.AppConfig.Int("rtmp_port"); err != nil {
		AppConfig.RTMPPort = 1935
	}

	// MongoDB URI
	if AppConfig.MongoDBURI = beego.AppConfig.String("mongodb_uri"); AppConfig.MongoDBURI == "" {
		panic("mongodb_uri is missing in app.conf")
	}

	// MongoDB DB Name
	if AppConfig.MongoDBName = beego.AppConfig.String("mongodb_name"); AppConfig.MongoDBName == "" {
		panic("mongodb_name is missing in app.conf")
	}

	// Google Project ID
	if AppConfig.GoogleProjectID = beego.AppConfig.String("google_project_id"); AppConfig.GoogleProjectID == "" {
		panic("google_project_id is missing in app.conf")
	}

	// Google Service Worker JSON Key
	if AppConfig.JSONKey = beego.AppConfig.String("google_api_key"); AppConfig.JSONKey == "" {
		panic("google_api_key is missing in app.conf")
	}

	// PubSub Transcode Topic Name
	if AppConfig.TopicNameTranscode = beego.AppConfig.String("pubsub_topic_transcode"); AppConfig.TopicNameTranscode == "" {
		panic("pubsub_topic_transcode is missing in app.conf")
	}

	// PubSub Complete Subscription Name
	if AppConfig.SubscriptionNameComplete = beego.AppConfig.String("pubsub_subscription_complete"); AppConfig.SubscriptionNameComplete == "" {
		panic("pubsub_subscription_complete is missing in app.conf")
	}

	// Uploads Directory
	if AppConfig.UploadsDirectory = beego.AppConfig.String("uploads_dir"); AppConfig.UploadsDirectory == "" {
		panic("uploads_dir is missing in app.conf")
	}

	// Mailgun Domain
	if AppConfig.MailgunDomain = beego.AppConfig.String("mailgun_domain"); AppConfig.MailgunDomain == "" {
		panic("mailgun_domain is missing in app.conf")
	}

	// Mailgun API Key
	if AppConfig.MailgunAPIKey = beego.AppConfig.String("mailgun_api_key"); AppConfig.MailgunAPIKey == "" {
		panic("mailgun_api_key is missing in app.conf")
	}

	// Hostname
	if AppConfig.Hostname = beego.AppConfig.String("hostname"); AppConfig.Hostname == "" {
		panic("hostname is missing in app.conf")
	}
}
