package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oldweipro/claude-to-chatgpt/api"
	"github.com/oldweipro/claude-to-chatgpt/core"
	"github.com/oldweipro/claude-to-chatgpt/middleware"
	"github.com/oldweipro/claude-to-chatgpt/service"
	"strings"
)

func main() {
	core.NewViper()
	if core.ServerConfig.Claude.SessionKey == "" {
		fmt.Println("SessionKey cannot be empty.")
		return
	}
	if core.ServerConfig.Claude.BaseUrl == "" {
		core.ServerConfig.Claude.BaseUrl = "https://claude.ai"
	}
	if !strings.HasPrefix(core.ServerConfig.Claude.SessionKey, "sessionKey=") {
		core.ServerConfig.Claude.SessionKey = "sessionKey=" + core.ServerConfig.Claude.SessionKey
	}
	// 获取OrganizationUuid验证网络及参数是否正确
	organizations, err := service.GetOrganizations()
	if err != nil {
		return
	}
	if core.ServerConfig.Claude.OrganizationUuid == "" {
		core.ServerConfig.Claude.OrganizationUuid = organizations[0].Uuid
	}
	r := gin.Default()
	r.Use(middleware.CORS)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.OPTIONS("/v1/chat/completions", api.OptionsHandler)
	r.POST("/v1/chat/completions", api.ChatCompletionsHandler)
	err = r.Run(":8787")
	if err != nil {
		return
	}
}
