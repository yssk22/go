package response

// HTTPStatus is an enum for http status code
//go:generate enum -type=HTTPStatus
type HTTPStatus int

// Available status code for 2xx
const (
	HTTPStatusOK HTTPStatus = 200 + iota
	HTTPStatusCreated
	HTTPStatusAccepted
	HTTPStatusNonAuthoritativeInformation
	HTTPStatusNoContent
	HTTPStatusResetContent
	HTTPStatusPartialContent
)

// Available status code for 3xx
const (
	HTTPStatusMovedParmanently HTTPStatus = 301 + iota
	HTTPStatusFound
	HTTPStatusSeeOther
	HTTPStatusNotModified
	HTTPStatusUseProxy
)

// Available status code for 4xx
const (
	HTTPStatusBadRequest HTTPStatus = 400 + iota
	HTTPStatusUnauthorized
	HTTPStatusPaymentRequired
	HTTPStatusForbidden
	HTTPStatusNotFound
	HTTPStatusMethodNotAllowed
	HTTPStatusNotAcceptable
	HTTPStatusProxyAuthenticationRequired
	HTTPStatusRequestTimeout
	HTTPStatusCoflict
	HTTPStatusGone
	HTTPStatusLengthRequired
	HTTPStatusPreconditionFailed
	HTTPStatusPayloadTooLarge
	HTTPStatusURITooLong
	HTTPStatusRangeNotSatisfiable
	HTTPStatusExpectationFailed
)

// Available status code for 5xx
const (
	HTTPStatusInternalServerError HTTPStatus = 500 + iota
	HTTPStatusNotImplemented
	HTTPStatusBadGateway
	HTTPStatusServiceUnavailable
	HTTPStatusGatewayTimeout
)
