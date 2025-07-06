package domain

import (
	"errors"
	"net/http"
)

var ErrGiveawayNotFound = errors.New("giveaway not found")

type HTTPError struct {
	Code    int
	Message string
	Err     error
}

func (e *HTTPError) Error() string {
	return e.Message
}

func NewNotFoundError(message string) *HTTPError {
	return &HTTPError{
		Code:    http.StatusNotFound,
		Message: message,
		Err:     ErrGiveawayNotFound,
	}
}

func NewInternalServerError(message string, err error) *HTTPError {
	return &HTTPError{
		Code:    http.StatusInternalServerError,
		Message: message,
		Err:     err,
	}
}
