package service

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	http2 "github.com/bogdanfinn/fhttp"
	"github.com/bogdanfinn/tls-client"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/oldweipro/claude-to-chatgpt/global"
	"github.com/oldweipro/claude-to-chatgpt/model"
	"io"
	"time"
)

var (
	jar     = tls_client.NewCookieJar()
	options = []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(3600),
		tls_client.WithClientProfile(tls_client.Safari_Ipad_15_6),
		tls_client.WithNotFollowRedirects(),
		// create cookieJar instance and pass it as argument
		tls_client.WithCookieJar(jar),
		// Disable SSL verification
		tls_client.WithInsecureSkipVerify(),
	}
	client, _ = tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
)

func RequestClaudeToResponse(c *gin.Context, params *model.ChatMessageRequest, stream bool) {
	sessionKey := GetSessionKey(c)
	organizationUuid, err := GetOrganizations(sessionKey)
	if err != nil {
		HandleErrorResponse(c, err.Error())
		return
	}
	appendMessageApi := global.ServerConfig.BaseUrl + "/api/append_message"
	err = client.SetProxy(global.ServerConfig.HttpProxy)
	if err != nil {
		HandleErrorResponse(c, err.Error())
		return
	}
	// 设置两个参数
	newStringUuid := uuid.NewString()
	_, err = CreateChatConversations(newStringUuid, sessionKey)
	if err != nil {
		HandleErrorResponse(c, err.Error())
		return
	}
	defer func(newStringUuid, sessionKey string) {
		err := DeleteChatConversations(newStringUuid, sessionKey)
		if err != nil {
			fmt.Println("delete err:", err, newStringUuid)
		}
	}(newStringUuid, sessionKey)
	params.ConversationUuid = newStringUuid
	params.OrganizationUuid = organizationUuid
	// 发起请求
	marshal, err := json.Marshal(params)
	if err != nil {
		HandleErrorResponse(c, err.Error())
		return
	}
	request, err := http2.NewRequest(http2.MethodPost, appendMessageApi, bytes.NewBuffer(marshal))
	if err != nil {
		HandleErrorResponse(c, err.Error())
		return
	}
	SetHeaders(request, sessionKey)
	response, err := client.Do(request)
	if response.StatusCode != 200 {
		HandleErrorResponse(c, "claude错误")
		return
	}
	reader := bufio.NewReader(response.Body)
	var originalResponse model.ChatMessageResponse
	var isRole = true
	if stream {
		// Response content type is text/event-stream
		c.Header("Content-Type", "text/event-stream")
	} else {
		// Response content type is application/json
		c.Header("Content-Type", "application/json")
	}
	var fullResponseText string
	completionResponse := model.ChatCompletionStreamResponse{
		ID:      "chatcmpl-7f1DmyzTWtiysnyfSS4i187kus2Ao",
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   "gpt-3.5-turbo-0613",
		Choices: []model.ChatCompletionStreamChoice{
			{
				Index: 0,
				Delta: model.ChatCompletionStreamChoiceDelta{
					Content: originalResponse.Completion,
				},
				FinishReason: nil,
			},
		},
	}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
		if len(line) < 6 {
			continue
		}
		line = line[6:]
		if isRole {
			completionResponse.Choices[0].Delta.Role = "assistant"
		} else {
			completionResponse.Choices[0].Delta.Content = originalResponse.Completion
			fullResponseText += originalResponse.Completion
		}
		completionResponse.Choices[0].Delta.Role = ""
		isRole = false
		if stream {
			resp, _ := json.Marshal(completionResponse)
			responseString := "data: " + string(resp) + "\n\n"
			c.Writer.WriteString(responseString)
			c.Writer.Flush()
		}
		err = json.Unmarshal([]byte(line), &originalResponse)
		if err != nil {
			continue
		}
		if originalResponse.Stop != "" && stream {
			completionResponse.Choices[0].FinishReason = "stop"
			completionResponse.Choices[0].Delta = model.ChatCompletionStreamChoiceDelta{}
			resp, _ := json.Marshal(completionResponse)
			responseString := "data: " + string(resp) + "\n\n"
			c.Writer.WriteString(responseString)
			c.Writer.Flush()
		}
	}
	if stream {
		c.Writer.WriteString("data: [DONE]\n\n")
		c.Writer.Flush()
	} else {
		c.JSON(200, NewChatCompletion(fullResponseText))
	}
}

func HandleErrorResponse(c *gin.Context, err string) {
	fmt.Println(err)
	c.JSON(500, gin.H{"error": gin.H{
		"message": "Unknown error",
		"type":    "internal_server_error",
		"param":   nil,
		"code":    "500",
		"details": err,
	}})
}

func CreateChatConversations(newStringUuid, sessionKey string) (model.ChatConversationResponse, error) {
	var chatConversationResponse model.ChatConversationResponse
	organizationUuid, err := GetOrganizations(sessionKey)
	if err != nil {
		return chatConversationResponse, err
	}
	chatConversationsApi := global.ServerConfig.BaseUrl + "/api/organizations/" + organizationUuid + "/chat_conversations"
	err = client.SetProxy(global.ServerConfig.HttpProxy)
	if err != nil {
		return chatConversationResponse, err
	}
	conversation := model.NewChatConversationRequest(newStringUuid, "")
	marshal, err := json.Marshal(conversation)
	if err != nil {
		return chatConversationResponse, err
	}
	request, err := http2.NewRequest(http2.MethodPost, chatConversationsApi, bytes.NewBuffer(marshal))

	if err != nil {
		return chatConversationResponse, err
	}
	SetHeaders(request, sessionKey)

	res, err := client.Do(request)
	if err != nil {
		return chatConversationResponse, err
	}
	defer res.Body.Close()

	if res.StatusCode != 201 {
		return chatConversationResponse, errors.New("Claude创建会话出错: " + res.Status)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return chatConversationResponse, err
	}
	err = json.Unmarshal(body, &chatConversationResponse)
	if err != nil {
		return chatConversationResponse, err
	}
	return chatConversationResponse, err
}

func DeleteChatConversations(newStringUuid, sessionKey string) error {
	organizationUuid, err := GetOrganizations(sessionKey)
	if err != nil {
		return err
	}
	err = client.SetProxy(global.ServerConfig.HttpProxy)
	if err != nil {
		return err
	}
	chatConversationsApi := global.ServerConfig.BaseUrl + "/api/organizations/" + organizationUuid + "/chat_conversations/"
	request, err := http2.NewRequest(http2.MethodDelete, chatConversationsApi+newStringUuid, nil)
	if err != nil {
		return err
	}
	SetHeaders(request, sessionKey)

	res, err := client.Do(request)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	//body, err := io.ReadAll(res.Body)
	//if err != nil {
	//	fmt.Println(err)
	//	return err
	//}
	//fmt.Println("delete:", newStringUuid, string(body))
	if res.StatusCode != 204 {
		all, _ := io.ReadAll(res.Body)
		return errors.New("delete chat conversations err" + string(all))
	}
	return nil
}

func GetOrganizations(sessionKey string) (string, error) {
	err := client.SetProxy(global.ServerConfig.HttpProxy)
	if err != nil {
		return "", err
	}
	organizationsApi := global.ServerConfig.BaseUrl + "/api/organizations"
	request, err := http2.NewRequest(http2.MethodGet, organizationsApi, nil)

	if err != nil {
		return "", err
	}
	SetHeaders(request, sessionKey)
	res, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("Claude获取组织出错: %s, %s", res.Status, string(body)))
	}
	var response []model.OrganizationsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}
	return response[0].Uuid, err
}

func SetHeaders(r *http2.Request, sessionKey string) {
	r.Header.Add("Cookie", sessionKey)
	r.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Accept", "*/*")
	r.Header.Add("Host", "claude.ai")
	r.Header.Add("Connection", "keep-alive")
}

func NewChatCompletion(fullResponseText string) model.ChatCompletionResponse {
	return model.ChatCompletionResponse{
		ID:      "chatcmpl-QXlha2FBbmROaXhpZUFyZUF3ZXNvbWUK",
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   "gpt-3.5-turbo-0613",
		Usage: model.Usage{
			PromptTokens:     0,
			CompletionTokens: 0,
			TotalTokens:      0,
		},
		Choices: []model.ChatCompletionChoice{
			{
				Message: model.ChatCompletionMessage{
					Content: fullResponseText,
					Role:    "assistant",
				},
				Index:        0,
				FinishReason: "stop",
			},
		},
	}
}
