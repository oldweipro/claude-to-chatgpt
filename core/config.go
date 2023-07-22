package core

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const (
	ConfigDefaultFile = "config.yaml"
)

var ServerConfig Config

func NewViper() {
	var config string
	config = ConfigDefaultFile
	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	// TODO check
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err = v.Unmarshal(&ServerConfig); err != nil {
			fmt.Println(err)
		}
	})
	if err = v.Unmarshal(&ServerConfig); err != nil {
		fmt.Println(err)
	}
}
