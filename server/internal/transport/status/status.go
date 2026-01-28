package httpstatus

import (
	"context"
	"errors"
	"net/http"

	coreerror "github.com/got-many-wheels/dwarf/server/internal/core/error"
)

// Status resolves the HTTP status code for any error.
func Status(err error) int {
	if err == nil {
		return http.StatusOK
	}

	var se statusCoder
	if errors.As(err, &se) && se.StatusCode() != 0 {
		return se.StatusCode()
	}

	if derr, ok := coreerror.As(err); ok {
		return statusForCode(derr.Code)
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return http.StatusGatewayTimeout
	}
	if errors.Is(err, context.Canceled) {
		return http.StatusRequestTimeout
	}

	return http.StatusInternalServerError
}

// WithStatus wraps err with an explicit HTTP status override.
func WithStatus(err error, status int) error {
	if err == nil {
		return nil
	}
	return &statusCodeError{err: err, status: status}
}

type statusCoder interface {
	StatusCode() int
}

type statusCodeError struct {
	err    error
	status int
}

func (e *statusCodeError) Error() string   { return e.err.Error() }
func (e *statusCodeError) Unwrap() error   { return e.err }
func (e *statusCodeError) StatusCode() int { return e.status }

func statusForCode(code coreerror.Code) int {
	switch code {
	case coreerror.CodeInvalid:
		return http.StatusBadRequest
	case coreerror.CodeUnauthorized:
		return http.StatusUnauthorized
	case coreerror.CodeForbidden:
		return http.StatusForbidden
	case coreerror.CodeNotFound:
		return http.StatusNotFound
	case coreerror.CodeConflict:
		return http.StatusConflict
	case coreerror.CodeRateLimited:
		return http.StatusTooManyRequests
	case coreerror.CodeTimeout:
		return http.StatusGatewayTimeout
	case coreerror.CodeUnavailable:
		return http.StatusServiceUnavailable
	case coreerror.CodeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
