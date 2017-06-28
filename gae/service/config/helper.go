package config

import (
	"net/http"
	"time"

	"github.com/speedland/go/retry"
	"github.com/speedland/go/x/xlog"
	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
)

// HTTPClientLoggerKey is a xlog key for this package
const HTTPClientLoggerKey = "gae.service.config.http"

const (
	ckURLFetchDeadline                = "urlfetch.deadline"
	ckURLFetchAllowInvalidCertificate = "urlfetch.allow_invalid_certificate"
	ckURLFetchMaxRetries              = "urlfetch.max_retries"
	ckURLFetchRetryBackoff            = "urlfetch.retry_backoff"
)

func init() {
	Global(ckURLFetchDeadline, "30", "timeout seconds for urlfetch")
	Global(ckURLFetchAllowInvalidCertificate, "0", "allow urlfetch to access the invalid certificate host")
	Global(ckURLFetchMaxRetries, "5", "allow urlfetch to retry request if an error occurred")
	Global(ckURLFetchRetryBackoff, "1s", "retry backoff configuration (time.Duration format)")
}

// NewHTTPClient returns a new *http.Client based on configurations
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

	ctx, _ = context.WithTimeout(ctx, time.Duration(deadline)*time.Second)

	return &http.Client{
		Transport: retry.NewHTTPTransport(
			&urlfetch.Transport{
				Context: ctx,
				AllowInvalidServerCertificate: allowInvalidCert == 1,
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
