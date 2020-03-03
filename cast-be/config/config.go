package config

import (
	"github.com/astaxie/beego"
)

var AppConfig Config

type Config struct {
	JWTSecret string

	Debug bool
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
}
