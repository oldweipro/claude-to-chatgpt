package main

import (
	"github.com/gin-gonic/gin"
	"github.com/oldweipro/claude-to-chatgpt/api"
	"github.com/oldweipro/claude-to-chatgpt/initialize"
	"github.com/oldweipro/claude-to-chatgpt/middleware"
)

func main() {
	initialize.NewViper()
	r := gin.Default()
	r.Use(middleware.CORS)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.OPTIONS("/v1/chat/completions", api.OptionsHandler)
	r.POST("/v1/chat/completions", api.ChatCompletionsHandler)
	err := r.Run(":8787")
	if err != nil {
		return
	}
}
