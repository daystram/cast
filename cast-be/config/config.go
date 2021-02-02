package config

import (
	"log"

	"github.com/astaxie/beego"
	"github.com/spf13/viper"
)

var AppConfig Config

type Config struct {
	Domain    string
	JWTSecret string
	Debug     bool

	RTMPPort int

	MongoDBURI  string
	MongoDBName string

	RabbitMQURI           string
	RabbitMQQueueTask     string
	RabbitMQQueueProgress string

	S3URI       string
	S3Bucket    string
	S3Region    string
	S3AccessKey string
	S3SecretKey string

	RatifyIssuer       string
	RatifyClientID     string
	RatifyClientSecret string

	UploadsDirectory string
}

// Load configuration
func InitializeAppConfig() {
	_ = beego.LoadAppConfig("ini", "./config/app.conf")

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/")
	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv()
	_ = viper.ReadInConfig()

	if AppConfig.Domain = viper.GetString("DOMAIN"); AppConfig.Domain == "" {
		log.Fatalln("[INIT] DOMAIN is not set")
	}
	AppConfig.JWTSecret = viper.GetString("SECRET")
	AppConfig.Debug = viper.GetBool("DEBUG")

	AppConfig.RTMPPort = viper.GetInt("RTMP_PORT")

	if AppConfig.RatifyIssuer = viper.GetString("RATIFY_ISSUER"); AppConfig.RatifyIssuer == "" {
		log.Fatalln("[INIT] RATIFY_ISSUER is not set")
	}
	if AppConfig.RatifyClientID = viper.GetString("RATIFY_CLIENT_ID"); AppConfig.RatifyClientID == "" {
		log.Fatalln("[INIT] RATIFY_CLIENT_ID is not set")
	}
	if AppConfig.RatifyClientSecret = viper.GetString("RATIFY_CLIENT_SECRET"); AppConfig.RatifyClientSecret == "" {
		log.Fatalln("[INIT] RATIFY_CLIENT_SECRET is not set")
	}

	if AppConfig.MongoDBURI = viper.GetString("MONGODB_URI"); AppConfig.MongoDBURI == "" {
		log.Fatalln("[INIT] MONGODB_URI is not set")
	}
	if AppConfig.MongoDBName = viper.GetString("MONGODB_NAME"); AppConfig.MongoDBName == "" {
		log.Fatalln("[INIT] MONGODB_NAME is not set")
	}

	if AppConfig.RabbitMQURI = viper.GetString("RABBITMQ_URI"); AppConfig.RabbitMQURI == "" {
		log.Fatalln("[INIT] RABBITMQ_URI is not set")
	}
	if AppConfig.RabbitMQQueueTask = viper.GetString("RABBITMQ_QUEUE_TASK"); AppConfig.RabbitMQQueueTask == "" {
		log.Fatalln("[INIT] RABBITMQ_QUEUE_TASK is not set")
	}
	if AppConfig.RabbitMQQueueProgress = viper.GetString("RABBITMQ_QUEUE_PROGRESS"); AppConfig.RabbitMQQueueProgress == "" {
		log.Fatalln("[INIT] RABBITMQ_QUEUE_PROGRESS is not set")
	}

	if AppConfig.S3URI = viper.GetString("S3_URI"); AppConfig.S3URI == "" {
		log.Fatalln("[INIT] S3_URI is not set")
	}
	if AppConfig.S3Bucket = viper.GetString("S3_BUCKET"); AppConfig.S3Bucket == "" {
		log.Fatalln("[INIT] S3_BUCKET is not set")
	}
	if AppConfig.S3Region = viper.GetString("S3_REGION"); AppConfig.S3Region == "" {
		log.Fatalln("[INIT] S3_REGION is not set")
	}
	if AppConfig.S3AccessKey = viper.GetString("S3_ACCESS_KEY"); AppConfig.S3AccessKey == "" {
		log.Fatalln("[INIT] S3_ACCESS_KEY is not set")
	}
	if AppConfig.S3SecretKey = viper.GetString("S3_SECRET_KEY"); AppConfig.S3SecretKey == "" {
		log.Fatalln("[INIT] S3_SECRET_KEY is not set")
	}

	if AppConfig.UploadsDirectory = viper.GetString("UPLOADS_DIR"); AppConfig.UploadsDirectory == "" {
		log.Fatalln("[INIT] UPLOADS_DIR is not set")
	}
}
