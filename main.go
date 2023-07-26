package main

import (
	"github.com/gin-gonic/gin"
	"github.com/oldweipro/claude-to-chatgpt/handles"
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
	r.OPTIONS("/v1/chat/completions", handles.OptionsHandler)
	r.POST("/v1/chat/completions", handles.ChatCompletionsHandler)
	err := r.Run(":8787")
	if err != nil {
		return
	}
}
