package service

import (
	"encoding/json"
	"fmt"
	"github.com/oldweipro/claude-to-chatgpt/model"
)

// OpenaiToClaudeParams 转换成claude的参数
func OpenaiToClaudeParams(chatCompletionRequest model.ChatCompletionRequest) *model.ChatMessageRequest {
	completionMessages := chatCompletionRequest.Messages
	marshal, err := json.Marshal(completionMessages)
	if err != nil {
		fmt.Println("Marshal err:", err)
	}
	content := string(marshal)
	return model.NewChatMessageRequest(content)
}
