package httperror

import (
	// internal golang package
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Kind ...
type Kind = int

// ErrorWithStatusCode error with http status code
type ErrorWithStatusCode struct {
	Err        string
	StatusCode int
}

func (e *ErrorWithStatusCode) Error() string {
	return e.Err
}

// Error encapsulate error with type of error
type Error struct {
	err     string
	kind    Kind
	message string
}

func New(err error, kind Kind, message string) *Error {
	return &Error{
		err:     fmt.Sprintf("%v", err),
		kind:    kind,
		message: message,
	}
}

// EncodeError ...
func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	code := http.StatusInternalServerError
	if sc, ok := err.(*ErrorWithStatusCode); ok {
		code = sc.StatusCode
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
