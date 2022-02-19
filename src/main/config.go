package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

var Conf = viper.New()

func confInit() {
	Conf.SetConfigType("toml")
	Conf.SetConfigName("avalon-aria2")
	Conf.AddConfigPath(`./soft/avalon/config/`)
	Conf.SetDefault("service.addr", ":6211")
	Conf.SetDefault("aria2.server.host", "127.0.0.1:6800")
	Conf.SetDefault("aria2.server.path", "jsonrpc")
	Conf.SetDefault("aria2.server.autoPause", 300)
	Conf.SetDefault("aria2.server.tempPath", "./")
	replacer := strings.NewReplacer(".", "_")
	Conf.SetEnvKeyReplacer(replacer)
	err := Conf.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		_, err := os.Create("./avalon-aria2.toml")
		if err != nil {
			log.Println(err)
			return
		}
	}

	err = Conf.WriteConfig()
	if err != nil {
		log.Println(err)
		return
	}

	Conf.WatchConfig()
	Conf.OnConfigChange(func(in fsnotify.Event) {
		err := Conf.ReadInConfig()
		if err != nil {
			log.Println(err)
		}
	})
}

func init() {
	confInit()
}
