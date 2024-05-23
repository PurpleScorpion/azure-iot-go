package models

const (
	USER_AGENT     = "User-Agent"
	ACCEPT         = "Accept"
	ACCEPT_VALUE   = "application/json"
	ACCEPT_CHARSET = "charset=utf-8"
	CONTENT_TYPE   = "Content-Type"
	CONTENT_LENGTH = "Content-Length"
	AUTHORIZATION  = "Authorization"
	REQUEST_ID     = "Request-Id"
	IF_MATCH       = "If-Match"
)

type HttpRequest struct {
	Method  string
	Url     string
	Headers map[string]string
	Body    []byte
}
