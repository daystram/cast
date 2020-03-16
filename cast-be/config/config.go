package config

import (
	"github.com/astaxie/beego"
)

var AppConfig Config

type Config struct {
	JWTSecret string
	Debug     bool

	MongoDBURI  string
	MongoDBName string
}

func InitializeAppConfig() {
	_ = beego.LoadAppConfig("ini", "config/app.conf")
	var err error

	// JWT secret
	if AppConfig.JWTSecret = beego.AppConfig.String("secret"); AppConfig.JWTSecret == "" {
		panic("secret in app.conf is missing or is in wrong format")
	}

	// Debug
	if AppConfig.Debug, err = beego.AppConfig.Bool("debug"); err != nil {
		panic("debug is missing in app.conf")
	}

	// MongoDB URI
	if AppConfig.MongoDBURI = beego.AppConfig.String("mongodb_uri"); AppConfig.MongoDBURI == "" {
		panic("mongodb_uri is missing in app.conf")
	}

	// MongoDB DB Name
	if AppConfig.MongoDBName = beego.AppConfig.String("mongodb_name"); AppConfig.MongoDBName == "" {
		panic("mongodb_name is missing in app.conf")
	}
}
