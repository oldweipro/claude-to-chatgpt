package model

import "time"

type ChatMessageRequest struct {
	Completion       *Completion   `json:"completion"`
	OrganizationUuid string        `json:"organization_uuid"`
	ConversationUuid string        `json:"conversation_uuid"`
	Text             string        `json:"text"`
	Attachments      []Attachments `json:"attachments"`
}
type ChatMessageResponse struct {
	Completion   string       `json:"completion"`
	StopReason   string       `json:"stop_reason"`
	Model        string       `json:"model"`
	Stop         string       `json:"stop"`
	LogId        string       `json:"log_id"`
	MessageLimit MessageLimit `json:"messageLimit"`
}

type MessageLimit struct {
	Type string `json:"type"`
}

func NewChatMessageRequest(text, history string) *ChatMessageRequest {
	var attachments []Attachments
	if history != "" {
		attachments = append(attachments, Attachments{
			FileName:         "paste.txt",
			FileType:         "txt",
			FileSize:         uint(len(history)),
			ExtractedContent: history,
		})

	}
	return &ChatMessageRequest{
		Completion:       NewCompletion(text),
		OrganizationUuid: "",
		ConversationUuid: "",
		Text:             text,
		Attachments:      attachments,
	}
}

type Completion struct {
	Prompt      string `json:"prompt"`
	Timezone    string `json:"timezone"`
	Model       string `json:"model"`
	Incremental bool   `json:"incremental"`
}

func NewCompletion(prompt string) *Completion {
	return &Completion{
		Prompt:      prompt,
		Timezone:    "Asia/Shanghai",
		Model:       "claude-2.1",
		Incremental: true,
	}
}

type Attachments struct {
	FileName         string `json:"file_name"`
	FileType         string `json:"file_type"`
	FileSize         uint   `json:"file_size"`
	ExtractedContent string `json:"extracted_content"`
}

type ChatConversationRequest struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
}

func NewChatConversationRequest(uuid, name string) *ChatConversationRequest {
	return &ChatConversationRequest{
		uuid,
		name,
	}
}

type ChatConversationResponse struct {
	Uuid      string `json:"uuid"`
	Name      string `json:"name"`
	Summary   string `json:"summary"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type OrganizationsResponse struct {
	Uuid         string       `json:"uuid"`
	Name         string       `json:"name"`
	Settings     Settings     `json:"settings"`
	Capabilities []string     `json:"capabilities"`
	JoinToken    string       `json:"join_token"`
	CreatedAt    string       `json:"created_at"`
	UpdatedAt    string       `json:"updated_at"`
	ActiveFlags  []ActiveFlag `json:"active_flags"`
}

type ActiveFlag struct {
	Id          string    `json:"id"`
	Type        string    `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
	DismissedAt time.Time `json:"dismissed_at"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type Settings struct {
	ClaudeConsolePrivacy string `json:"claude_console_privacy"`
}

type ErrorResponse struct {
	Error  Error  `json:"error"`
	Detail string `json:"detail"`
}
type Error struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}
