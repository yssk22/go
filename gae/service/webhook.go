package service

import (
	"path"

	"github.com/speedland/go/services/facebook/messenger"
)

func (s *Service) MessengerWebhook(webhook messenger.Webhook) {
	s.Post(path.Join(s.APIConfig.WebhookBasePath, "messenger/"), messenger.NewWebhookHandler(webhook))
}
