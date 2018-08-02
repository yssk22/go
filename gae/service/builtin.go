package service

import "github.com/speedland/go/services/facebook/messenger"

// BuiltInAPIConfig is a configuration object for setupBuiltInAPIs
type BuiltInAPIConfig struct {
	AuthNamespace    string
	MessengerWebHook messenger.Webhook
}
