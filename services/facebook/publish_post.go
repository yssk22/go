package facebook

import (
	"context"
	"fmt"
)

type PagePostParams struct {
	Message     string `json:"message,omitempty"`
	Link        string `json:"link,omitempty"`
	ActionLinks []struct {
		Name string `json:"name,omitempty"`
		Link string `json:"link,omitempty"`
	} `json:"action_links,omitempty"`
	Place                    string        `json:"place,omitempty"`
	Tags                     string        `json:"string,omitempty"`
	ObjectAttachment         string        `json:"object_attachment,omitempty"`
	Published                *bool         `json:"published,omitempty"`
	ScheduledPublishTime     int           `json:"scheduled_publish_time,omitempty"`
	BackdatedTime            int           `json:"backdated_time,omitempty"`
	BackdatedTimeGranularity string        `json:"backdated_time_granularity,omitempty"`
	ChildAttachments         []*Attachment `json:"child_attachments,omitempty"`
	MultiShareOptimized      *bool         `json:"multi_share_optimized,omitempty"`
	MultiShareEndCard        *bool         `json:"multi_share_end_card,omitempty"`
}

type Attachment struct {
	Link        string `json:"link,omitempty"`
	Picture     string `json:"picture,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	VideoID     string `json:"video_id,omitempty"`
}

func (c *Client) PublishPost(ctx context.Context, id string, params *PagePostParams) (string, error) {
	var r map[string]string
	err := c.Post(ctx, fmt.Sprintf("/%s/feed", id), nil, params, &r)
	if err != nil {
		return "", err
	}
	return r["id"], nil
}
