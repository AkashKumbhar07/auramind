package errors

import (
	"fmt"
	"net/http"
)

type Kind int

const (
	KindInvalid       Kind = iota
	KindNotFound      Kind = http.StatusNotFound
	KindUnauthorized  Kind = http.StatusUnauthorized
	KindForbidden     Kind = http.StatusForbidden
	KindConflict      Kind = http.StatusConflict
	KindBadRequest    Kind = http.StatusBadRequest
	KindInternal      Kind = http.StatusInternalServerError
	KindUnavailable   Kind = http.StatusServiceUnavailable
)

type Error struct {
	Kind    Kind
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func New(kind Kind, msg string) *Error {
	return &Error{Kind: kind, Message: msg}
}

func Wrap(kind Kind, msg string, err error) *Error {
	return &Error{Kind: kind, Message: msg, Err: err}
}

func NotFound(msg string) *Error {
	return New(KindNotFound, msg)
}

func Unauthorized(msg string) *Error {
	return New(KindUnauthorized, msg)
}

func Forbidden(msg string) *Error {
	return New(KindForbidden, msg)
}

func BadRequest(msg string) *Error {
	return New(KindBadRequest, msg)
}

func Conflict(msg string) *Error {
	return New(KindConflict, msg)
}

func Internal(msg string) *Error {
	return New(KindInternal, msg)
}

func IsNotFound(err error) bool {
	var e *Error
	if as, ok := err.(*Error); ok {
		return as.Kind == KindNotFound
	}
	_ = e
	return false
}

func KindOf(err error) Kind {
	var e *Error
	if as, ok := err.(*Error); ok {
		return as.Kind
	}
	_ = e
	return KindInternal
}
