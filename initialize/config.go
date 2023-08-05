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
	"strings"
)

const (
	ConfigDefaultFile = "config.yaml"
)

var (
	configFlag    = ""
	httpProxyFlag = ""
	baseUrlFlag   = ""
)

func init() {
	flag.StringVar(&configFlag, "c", "", "choose config file.")
	flag.StringVar(&httpProxyFlag, "http_proxy", "", "set http_proxy, for example http://127.0.0.1:7890")
	flag.StringVar(&baseUrlFlag, "base_url", "", "set base_url, for example https://claude.ai")
}

func NewViper() {
	flag.Parse()
	if configFlag == "" {
		configFlag = ConfigDefaultFile
	}
	// check config file
	_, err := os.Stat(configFlag)
	if os.IsNotExist(err) {
		file, err := os.Create(configFlag)
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
	v.SetConfigFile(configFlag)
	v.SetConfigType("yaml")
	// 设置默认值
	v.SetDefault("base-url", "https://claude.ai")
	err = v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		// 因为viper值如果为空（删除）不会复写原来的值，数组内的值删除会出现不生效问题，SessionKeys先置为空再赋值
		global.ServerConfig.Claude.SessionKeys = nil
		if err = v.Unmarshal(&global.ServerConfig); err != nil {
			fmt.Println(err)
		}
		PrintServerConfig()
	})
	if err = v.Unmarshal(&global.ServerConfig); err != nil {
		fmt.Println(err)
	}
	// 设置命令参数
	if baseUrlFlag != "" {
		global.ServerConfig.BaseUrl = baseUrlFlag
	}
	if httpProxyFlag != "" {
		global.ServerConfig.HttpProxy = httpProxyFlag
	}
	// 设置环境变量
	keysEnv := os.Getenv("CLAUDE_SESSION_KEYS")
	keys := strings.Split(keysEnv, ",")
	global.ServerConfig.Claude.SessionKeys = append(global.ServerConfig.Claude.SessionKeys, keys...)
	baseUrlEnv := os.Getenv("CLAUDE_BASE_URL")
	if baseUrlEnv != "" {
		global.ServerConfig.BaseUrl = baseUrlEnv
	}
	httpProxyEnv := os.Getenv("CLAUDE_HTTP_PROXY")
	if httpProxyEnv != "" {
		global.ServerConfig.HttpProxy = httpProxyEnv
	}
}

func PrintServerConfig() {
	indent, _ := json.MarshalIndent(global.ServerConfig, "", "    ")
	fmt.Println(string(indent))
}
