package api

import (
	"github.com/gin-gonic/gin"
	"github.com/oldweipro/claude-to-chatgpt/core"
	"github.com/oldweipro/claude-to-chatgpt/model"
	"github.com/oldweipro/claude-to-chatgpt/service"
)

func OptionsHandler(c *gin.Context) {
	// Set headers for CORS
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST")
	c.Header("Access-Control-Allow-Headers", "*")
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func ChatCompletionsHandler(c *gin.Context) {
	var chatCompletionRequest model.ChatCompletionRequest
	err := c.BindJSON(&chatCompletionRequest)
	if err != nil {
		c.JSON(400, gin.H{"error": gin.H{
			"message": "Request must be proper JSON",
			"type":    "invalid_request_error",
			"param":   nil,
			"code":    err.Error(),
		}})
		return
	}
	// 流程是：接收openAI的参数，转换成claude的参数，请求claude返回结果，claude结果转openAI，返回给用户。
	params := service.OpenaiToClaudeParams(chatCompletionRequest)
	service.RequestClaudeToResponse(c, params, core.ServerConfig.Claude.SessionKey, "http://127.0.0.1:7890", chatCompletionRequest.Stream)
}
