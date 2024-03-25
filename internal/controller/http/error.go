package http

import "fmt"

type httpError struct {
	Type       ErrorType `json:"-"`
	Code       ErrorCode `json:"code"`
	StatusCode int       `json:"-"`
	Message    string    `json:"message"`
}

func (e *httpError) Error() string {
	return fmt.Sprintf("%s: code=%s, message=%s", e.Type, e.Code, e.Message)
}
