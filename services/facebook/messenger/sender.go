package messenger

import (
	"net/http"

	"context"

	"github.com/yssk22/go/services/facebook"
)

// Sender is a struct to send messages
type Sender struct {
	client *facebook.Client
}

// NewSender returns a new messenger API client to send messages
func NewSender(client *http.Client, pageAccessToken string) *Sender {
	return &Sender{
		client: facebook.NewClient(client, pageAccessToken),
	}
}

// Send sends a message
func (s *Sender) Send(ctx context.Context, msg *SendMessage) (*SendResponse, error) {
	var resp SendResponse
	if err := s.client.Post(ctx, "/me/messages", nil, msg, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SendResponse is a response structure for Send API
type SendResponse struct {
	RecipientID string `json:"recipient_id"`
	MessageID   string `json:"message_dd"`
}

// SendMessage is a message structure for Send API
type SendMessage struct {
	Recipient *Recipient          `json:"recipient"`
	Message   *SendMessageContent `json:"message"`
}

// Recipient is a structure to specify a recipent in Send API
type Recipient struct {
	ID          string `json:"id,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Name        string `json:"name,omitempty"`
}

// SendMessageContent is a content structure in SendMessage
type SendMessageContent struct {
	Text       string           `json:"text,omitempty"`
	QuickReply []SendQuickReply `json:"quick_reply,omitempty"`
	Attachment *SendAttachment  `json:"attachment,omitempty"`
	Metadata   string           `json:"metadata,omitempty"`
}

// SendQuickReply is a quick reply structure for Send API
type SendQuickReply struct {
	ContentType string `json:"content_type"`
	Title       string `json:"title"`
	Payload     string `json:"payload"`
	ImageURL    string `json:"image_url"`
}

// SendAttachment is an attachment structure for Send API
type SendAttachment struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// SendPayload is a payload structure for Send AP{
type SendPayload struct {
	URL           string `json:"url"`
	IsReusable    string `json:"is_reusable"`
	AttachementID string `json:"attachment_id"`
}

// GenericTemplatePayload
type GenericTemplatePayload struct {
	TemplateType string                           `json:"template_type"`
	Elements     []*GenericTemplatePayloadElement `json:"elements,omitempty"`
}

type GenericTemplatePayloadElement struct {
	Title         string        `json:"title"`
	Subtitle      string        `json:"subtitle,omitempty"`
	ImageURL      string        `json:"image_url,omitempty"`
	DefaultAction *URLButton    `json:"default_action,omitempty"`
	Buttons       []interface{} `json:"buttons,omitempty"`
}

// URLButton is for url button structure described at
// https://developers.facebook.com/docs/messenger-platform/send-api-reference/url-button
type URLButton struct {
	Type                string `json:"type"`
	Title               string `json:"title,omitempty"`
	URL                 string `json:"url"`
	WebviewHeightRatio  string `json:"webview_height_ratio,omitempty"`
	MessengerExtentions bool   `json:"messenger_extentions,omitempty"`
	FallbackURL         string `json:"fallback_url,omitempty"`
	WebviewShareButton  string `json:"webview_share_button,omitempty"`
}

// Constants for URLButton.WebviewShareButton fields
const (
	WebviewShareButtonShow = "show"
	WebviewShareButtonHide = "hide"
)
