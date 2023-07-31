package service

import (
	"github.com/gin-gonic/gin"
	"github.com/oldweipro/claude-to-chatgpt/global"
	"strings"
)

func GetSessionKey(c *gin.Context) (sk string) {
	auth := c.Request.Header.Get("Authorization")
	hasPrefix := strings.HasPrefix(auth, "Bearer ")
	if hasPrefix && len(auth) > 7 {
		sk = "sessionKey=" + auth[7:]
	}
	if sk == "" {
		sk = global.ServerConfig.Claude.GetSessionKey()
	}
	return sk
}
