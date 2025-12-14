package pkg

import (
	"fmt"
	"net/http"
)

var ErrInternalServer = &Error{
	Message: "Internal error",
	Code:    "INTERNAL_SERVER_ERROR",
	Status:  http.StatusInternalServerError,
}

var ErrNotFound = &Error{
	Message: "Not found",
	Code:    "NOT_FOUND",
	Status:  http.StatusNotFound,
}

var ErrBadRequest = &Error{
	Message: "Bad request",
	Code:    "BAD_REQUEST",
	Status:  http.StatusBadRequest,
}

var ErrConflict = &Error{
	Message: "Conflict",
	Code:    "CONFLICT",
	Status:  http.StatusConflict,
}

var ErrUnauthorized = &Error{
	Message: "Unauthorized",
	Code:    "UNAUTHORIZED",
	Status:  http.StatusUnauthorized,
}

var ErrUnprocessableEntity = &Error{
	Message: "Unprocessable Entity",
	Code:    "UNPROCESSABLE_ENTITY",
	Status:  http.StatusUnprocessableEntity,
}

var allowStatusError = map[int]struct{}{
	http.StatusInternalServerError: {},
	http.StatusNotFound:            {},
	http.StatusBadRequest:          {},
	http.StatusUnauthorized:        {},
	http.StatusConflict:            {},
	http.StatusUnprocessableEntity: {},
}

type ErrorCode string

type Error struct {
	Message string
	Code    ErrorCode
	Status  int
	Info    string
}

func (e *Error) Error() string {
	if e == nil || e.Code == "" {
		return ""
	}

	if e.Info == "" {
		return fmt.Sprintf("[%s] %s", e.Code, e.Message)
	}

	return fmt.Sprintf("[%s] %s (%s)", e.Code, e.Message, e.Info)

}

func (e *Error) HTTPStatus() int {
	if e == nil {
		return http.StatusOK
	}
	if _, ok := allowStatusError[e.Status]; !ok {
		return http.StatusInternalServerError
	}
	return e.Status
}

func (e *Error) SetInfo(s string) *Error {
	return &Error{
		Message: e.Message,
		Code:    e.Code,
		Status:  e.Status,
		Info:    s,
	}
}

func (e *Error) SetInfof(s string, v ...any) *Error {
	return e.SetInfo(fmt.Sprintf(s, v...))
}
