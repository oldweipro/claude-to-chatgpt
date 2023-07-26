package initialize

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/oldweipro/claude-to-chatgpt/global"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"os"
)

const (
	ConfigDefaultFile = "config.yaml"
)

func NewViper() {
	var config string
	flag.StringVar(&config, "c", "", "choose config file.")
	flag.Parse()
	if config == "" {
		config = ConfigDefaultFile
	}
	// check config file
	_, err := os.Stat(config)
	if os.IsNotExist(err) {
		file, err := os.Create(config)
		// 其他处理
		if err != nil {
			return
		}
		defer file.Close()
		encoder := yaml.NewEncoder(file)
		encoder.SetIndent(2)
		if err := encoder.Encode(&global.ServerConfig); err != nil {
			panic(err)
		}
		fmt.Println("File created and data written successfully.")
	}
	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	err = v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err = v.Unmarshal(&global.ServerConfig); err != nil {
			fmt.Println(err)
		}
		SyncServerConfig()
		PrintServerConfig()
	})
	if err = v.Unmarshal(&global.ServerConfig); err != nil {
		fmt.Println(err)
	}
	SyncServerConfig()
}

func SyncServerConfig() {
	if global.ServerConfig.BaseUrl == "" {
		global.ServerConfig.BaseUrl = "https://claude.ai"
	}
	if global.ServerConfig.Proxy.Protocol != "" || global.ServerConfig.Proxy.Host != "" || global.ServerConfig.Proxy.Port != "" {
		if global.ServerConfig.Proxy.Username != "" || global.ServerConfig.Proxy.Password != "" {
			global.HttpProxy = global.ServerConfig.Proxy.Protocol + "://" + global.ServerConfig.Proxy.Username + ":" + global.ServerConfig.Proxy.Password + "@" + global.ServerConfig.Proxy.Host + ":" + global.ServerConfig.Proxy.Port
			fmt.Println(global.HttpProxy)
		} else {
			global.HttpProxy = global.ServerConfig.Proxy.Protocol + "://" + global.ServerConfig.Proxy.Host + ":" + global.ServerConfig.Proxy.Port
		}
	}
}

func PrintServerConfig() {
	indent, _ := json.MarshalIndent(global.ServerConfig, "", "    ")
	fmt.Println(string(indent))
}
