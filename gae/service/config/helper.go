package config

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/appengine/urlfetch"

	"github.com/speedland/go/retry"
	"github.com/speedland/go/services/facebook/messenger"
	"github.com/speedland/go/web/response/view/react"
	"github.com/speedland/go/x/xlog"
	"github.com/speedland/go/x/xtime"
)

// HTTPClientLoggerKey is a xlog key for this package
const HTTPClientLoggerKey = "gae.service.config.http"

const (
	ckURLFetchDeadline                = "urlfetch.deadline"
	ckURLFetchAllowInvalidCertificate = "urlfetch.allow_invalid_certificate"
	ckURLFetchMaxRetries              = "urlfetch.max_retries"
	ckURLFetchRetryBackoff            = "urlfetch.retry_backoff"
	ckFacebookPageID                  = "facebook.page_id"
	ckFacebookPageToken               = "facebook.page_token"
	ckFacebookAppID                   = "facebook.app_id"
	ckFacebookAppSecret               = "facebook.app_secret"
	ckFacebookPixelID                 = "facebook.pixel_id"
	ckMessengerVerificationToken      = "messenger.verification_token"
	ckGoogleAnalyticsID               = "google.analytics_id"
)

func init() {
	Global(ckURLFetchDeadline, "30", "timeout seconds for urlfetch")
	Global(ckURLFetchAllowInvalidCertificate, "0", "allow urlfetch to access the invalid certificate host")
	Global(ckURLFetchMaxRetries, "5", "allow urlfetch to retry request if an error occurred")
	Global(ckURLFetchRetryBackoff, "1s", "retry backoff configuration (time.Duration format)")
	Global(ckFacebookAppID, "", "facebook app id")
	Global(ckFacebookAppSecret, "", "facebook app secret")
	Global(ckFacebookPageID, "", "facebook page id")
	Global(ckFacebookPageToken, "", "facebook page access token")
	Global(ckFacebookPixelID, "", "facebook pixel id")
	Global(ckMessengerVerificationToken, "", "messenger verification token")
	Global(ckGoogleAnalyticsID, "", "google analytics ID")
}

// OAuth2Config is a configuration object for oauth2 clients
type OAuth2Config struct {
	ClientID     string
	ClientSecret string
}

// GenPageConfig returns *react.PageConfig from configuration datastore.
func (c *Config) GenPageConfig(ctx context.Context) *react.PageConfig {
	return &react.PageConfig{
		FacebookAppID:     c.GetValue(ctx, ckFacebookAppID),
		FacebookPageID:    c.GetValue(ctx, ckFacebookPageID),
		FacebookPixelID:   c.GetValue(ctx, ckFacebookPixelID),
		GoogleAnalyticsID: c.GetValue(ctx, ckGoogleAnalyticsID),
	}
}

// GetFacebookConfig returns OAuth2 config for facebook.
func (c *Config) GetFacebookConfig(ctx context.Context) *OAuth2Config {
	appID := c.GetValue(ctx, ckFacebookAppID)
	if appID == "" {
		return nil
	}
	appSecret := c.GetValue(ctx, ckFacebookAppSecret)
	return &OAuth2Config{
		ClientID:     appID,
		ClientSecret: appSecret,
	}
}

// GetMessengerVerificationToken returns messenger verification token string
func (c *Config) GetMessengerVerificationToken(ctx context.Context) string {
	return c.GetValue(ctx, ckMessengerVerificationToken)
}

type httpTransport struct {
	context          context.Context
	logger           *xlog.Logger
	base             http.RoundTripper
	deadline         int // second, zero means 5 second.
	allowInvalidCert bool
}

func (transport *httpTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	ctx, cancel := context.WithTimeout(transport.context, time.Duration(transport.deadline)*time.Second)
	defer cancel()
	base := &urlfetch.Transport{
		Context: ctx,
		AllowInvalidServerCertificate: transport.allowInvalidCert,
	}
	t := xtime.Benchmark(func() {
		resp, err = base.RoundTrip(req)
	})
	var status string
	if resp == nil {
		status = "FAIL"
	} else {
		status = resp.Status
	}
	transport.logger.Infof("[%s] %s (%s)", status, req.URL.String(), t)
	return resp, err
}

// NewHTTPClient is an http client available in this context
func (c *Config) NewHTTPClient(ctx context.Context) *http.Client {
	const maxRetryHardLimit = 30
	const deadlineHardLimit = 60
	logger := xlog.WithContext(ctx).WithKey(HTTPClientLoggerKey)
	deadline := c.GetIntValue(ctx, ckURLFetchDeadline)
	allowInvalidCert := c.GetIntValue(ctx, ckURLFetchAllowInvalidCertificate)
	maxRetries := c.GetIntValue(ctx, ckURLFetchMaxRetries)
	backoff := c.GetValue(ctx, ckURLFetchRetryBackoff)

	if deadline <= 0 || deadline > deadlineHardLimit {
		deadline = c.GetIntDefaultValue(ckURLFetchDeadline)
	}
	if maxRetries < 0 || maxRetries > maxRetryHardLimit {
		maxRetries = c.GetIntDefaultValue(ckURLFetchMaxRetries)
	}
	backoffDuration, err := time.ParseDuration(backoff)
	if err != nil {
		logger.Warnf("Could not set a backoff duration by %q: %v", backoff, err)
		backoffDuration, _ = time.ParseDuration(c.GetDefaultValue(ckURLFetchRetryBackoff))
	}

	return &http.Client{
		Transport: retry.NewHTTPTransport(
			&httpTransport{
				context:          ctx,
				logger:           logger,
				deadline:         deadline,
				allowInvalidCert: allowInvalidCert == 1,
			},
			&retryLogger{
				base: retry.HTTPAnd(
					retry.HTTPMaxRetry(maxRetries),
					retry.HTTPServerErrorChecker(),
				),
				maxRetries: maxRetries,
				logger:     logger,
			},
			retry.HTTPConstBackoff(backoffDuration),
		),
	}
}

// NewMessengerSender returns a new *messenger.Sender from server configuration
func (c *Config) NewMessengerSender(ctx context.Context) (*messenger.Sender, error) {
	httpClient := c.NewHTTPClient(ctx)
	pageToken := c.GetValue(ctx, ckFacebookPageToken)
	if pageToken == "" {
		return nil, fmt.Errorf("%q is configured in server configuration", ckFacebookPageToken)
	}
	return messenger.NewSender(httpClient, pageToken), nil
}

type retryLogger struct {
	base       retry.HTTPChecker
	maxRetries int
	logger     *xlog.Logger
}

func (rl *retryLogger) NeedRetry(attempt int, req *http.Request, resp *http.Response, err error) bool {
	needRetry := rl.base.NeedRetry(attempt, req, resp, err)
	if needRetry {
		if err != nil {
			rl.logger.Infof("[%d/%d] Retry attempting %s (last error: %v)", attempt, rl.maxRetries, req.URL.String(), err)
		} else {
			rl.logger.Infof("[%d/%d] Retry attempting %s (last status code: %d)", attempt, rl.maxRetries, req.URL.String(), resp.StatusCode)
		}
		return true
	}
	return false
}
