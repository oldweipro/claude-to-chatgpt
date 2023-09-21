package service

import (
	"encoding/json"
	"fmt"
	"github.com/oldweipro/claude-to-chatgpt/model"
)

// OpenaiToClaudeParams 转换成claude的参数
func OpenaiToClaudeParams(chatCompletionRequest model.ChatCompletionRequest) *model.ChatMessageRequest {
	completionMessages := chatCompletionRequest.Messages
	text := completionMessages[len(completionMessages)-1]
	history := completionMessages[:len(completionMessages)-1]
	textMarshal, err := json.Marshal(text)
	if err != nil {
		fmt.Println("Text marshal err:", err)
	}
	textMessage := string(textMarshal)
	historyMessage := ""
	if len(history) > 0 {
		historyMarshal, err := json.Marshal(history)
		if err != nil {
			fmt.Println("History marshal err:", err)
		}
		historyMessage = string(historyMarshal)
	}
	//textMessage = "附件中存放着对话的上下文，作答时请忽略上下文的json格式，请以最新的prompt作答：" + textMessage
	textMessage = "The attachment contains the context of the conversation. When answering, please ignore the json format of the context. Please answer with the latest prompt: " + textMessage
	return model.NewChatMessageRequest(textMessage, historyMessage)
}
