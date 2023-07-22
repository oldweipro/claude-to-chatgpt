package initialize

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/oldweipro/claude-to-chatgpt/global"
	"github.com/oldweipro/claude-to-chatgpt/service"
	"github.com/spf13/viper"
	"strings"
)

const (
	ConfigDefaultFile = "config.yaml"
)

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
		if err = v.Unmarshal(&global.ServerConfig); err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%s\n", global.ServerConfig)
	})
	if err = v.Unmarshal(&global.ServerConfig); err != nil {
		fmt.Println(err)
	}
	SyncServerConfig()
}

func SyncServerConfig() {
	if global.ServerConfig.Claude.SessionKey == "" {
		fmt.Println("SessionKey cannot be empty.")
		return
	}
	if global.ServerConfig.Claude.BaseUrl == "" {
		global.ServerConfig.Claude.BaseUrl = "https://claude.ai"
	}
	if !strings.HasPrefix(global.ServerConfig.Claude.SessionKey, "sessionKey=") {
		global.ServerConfig.Claude.SessionKey = "sessionKey=" + global.ServerConfig.Claude.SessionKey
	}
	if global.ServerConfig.Claude.OrganizationUuid == "" {
		// 获取OrganizationUuid验证网络及参数是否正确
		organizations, err := service.GetOrganizations()
		if err != nil {
			fmt.Println("get organizations err")
		}
		global.ServerConfig.Claude.OrganizationUuid = organizations[0].Uuid
	}
}
