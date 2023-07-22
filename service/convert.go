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
	//len := len(completionMessages)
	// 如果message只有一个的话，可以直接message
	//message := completionMessages[len-1]
	return model.NewChatMessageRequest(string(marshal))
}
