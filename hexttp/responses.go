package hexttp

import (
	"net/http"
)

type LogLevel string

const (
	LogNone  LogLevel = ""
	LogWarn  LogLevel = "warn"
	LogError LogLevel = "error"
)

type HTTPResponse struct {
	StatusCode int
	Body       any
	LogLevel   LogLevel
}

func NewHTTPResponse(status int, body any, logLevel LogLevel) *HTTPResponse {
	return &HTTPResponse{
		StatusCode: status,
		Body:       body,
		LogLevel:   logLevel,
	}
}

func OK(data any) *HTTPResponse {
	return NewHTTPResponse(http.StatusOK, data, LogNone)
}

func Created(data any) *HTTPResponse {
	return NewHTTPResponse(http.StatusCreated, data, LogNone)
}

func Updated(data any) *HTTPResponse {
	return NewHTTPResponse(http.StatusAccepted, data, LogNone)
}

func NoContent() *HTTPResponse {
	return NewHTTPResponse(http.StatusNoContent, nil, LogNone)
}

func InvalidJSON() *HTTPResponse {
	return NewHTTPResponse(http.StatusBadRequest, "invalid body request", LogWarn)
}

func InvalidRequestData(errors map[string]string) *HTTPResponse {
	return NewHTTPResponse(http.StatusBadRequest, errors, LogWarn)
}

func InvalidID() *HTTPResponse {
	return NewHTTPResponse(http.StatusBadRequest, "invalid id", LogWarn)
}

func NotFound(err error) *HTTPResponse {
	return NewHTTPResponse(http.StatusNotFound, err.Error(), LogWarn)
}

func AlreadyExist(err error) *HTTPResponse {
	return NewHTTPResponse(http.StatusConflict, err.Error(), LogWarn)
}
func InternalError(msg string) *HTTPResponse {
	return NewHTTPResponse(http.StatusInternalServerError, "internal error", LogError)
}

func Unauthorized() *HTTPResponse {
	return NewHTTPResponse(http.StatusUnauthorized, "unauthorized", LogWarn)
}

func Forbidden() *HTTPResponse {
	return NewHTTPResponse(http.StatusForbidden, "forbidden", LogWarn)
}
