package builtin

import (
	"context"

	"github.com/yssk22/go/services/facebook/messenger"
)

// TODO: issue #11 to generate API code
func postFacebookMessengerWebHook(ctx context.Context, message *messenger.ReceivedMessage) (bool, error) {
	return true, nil
}
