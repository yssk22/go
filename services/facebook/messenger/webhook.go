package messenger

import (
	"io"

	"io/ioutil"

	"bytes"

	"github.com/speedland/go/web"
	"github.com/speedland/go/web/response"
	"github.com/speedland/go/x/xlog"
	"context"
)

func NewVericationHandler(token string) web.Handler {
	return web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		if req.Query.GetStringOr("hub.mode", "") == "subscribe" &&
			req.Query.GetStringOr("hub.verify_token", "") == token {
			return response.NewText(req.Query.GetStringOr("hub.challenge", ""))
		}
		return response.NewTextWithStatus("invalid request", response.HTTPStatusForbidden)
	})
}

// Webhook is an interface to process messenger message in webhook
type Webhook interface {
	Process(context.Context, *ReceivedMessage) error
}

// WebhookFunc is a simple interface generator for a webhook function
type WebhookFunc func(context.Context, *ReceivedMessage) error

// Process implements Webhook#Process
func (f WebhookFunc) Process(ctx context.Context, messages *ReceivedMessage) error {
	return f(ctx, messages)
}

// MaxWebhookPayloadSize is a maximum size of payload (512KB) in a single webhook call
const MaxWebhookPayloadSize = 1024 * 1024 * 512

// NewWebhookHandler return s new web.Handler to handle webhook request
func NewWebhookHandler(hook Webhook) web.Handler {
	return web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		ctx := req.Context()
		logger := xlog.WithContext(ctx).WithPrefix("[Messenger WebHook] ")
		buff, err := ioutil.ReadAll(&io.LimitedReader{R: req.Body, N: MaxWebhookPayloadSize})
		if err != nil {
			logger.Infof("io error: %v (Content-Length = %s)", err, req.Header.Get("content-length"))
			return response.NewTextWithStatus("io error", response.HTTPStatusForbidden)
		}
		messages, err := Parse(bytes.NewBuffer(buff))
		if err != nil {
			logger.Errorf("could not parse the messenger payload: %v -- %s", err, string(buff))
		} else {
			var noContentCount int
			for _, message := range messages {
				if message.Content == nil {
					noContentCount++
					continue
				}
				if err = hook.Process(ctx, &message); err != nil {
					logger.Errorf("webhook error: %v", err)
				}
			}
			if noContentCount > 0 {
				logger.Warnf("%d messages are ignored because of no content, raw json ---> %s", noContentCount, string(buff))
			}
		}
		return response.NewText("OK")
	})
}
