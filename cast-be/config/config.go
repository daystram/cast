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

	GoogleProjectID          string
	JSONKey                  string
	TopicNameTranscode       string
	SubscriptionNameComplete string

	UploadsDirectory string

	MailgunDomain string
	MailgunAPIKey string
	Hostname      string
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

	AppConfig.MongoDBURI = viper.GetString("MONGODB_URI")
	AppConfig.MongoDBName = viper.GetString("MONGODB_NAME")

	if AppConfig.UploadsDirectory = viper.GetString("UPLOADS_DIR"); AppConfig.UploadsDirectory == "" {
		log.Fatalln("[INIT] UPLOADS_DIR is not set")
	}
}
