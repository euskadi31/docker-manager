package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/xlog"
)

// NotFoundFailure response
func NotFoundFailure(w http.ResponseWriter, r *http.Request) {
	Failure(w, http.StatusNotFound, ErrorMessage{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf(`No route found for \"%s %s\"`, r.Method, r.URL.Path),
	})
}

// FailureFromError write ErrorMessage from error
func FailureFromError(w http.ResponseWriter, status int, err error) {
	xlog.Error(err)

	Failure(w, status, ErrorMessage{
		Code:    status,
		Message: err.Error(),
	})
}

// Failure response
func Failure(w http.ResponseWriter, status int, err ErrorMessage) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)

	body := ErrorResponse{
		Error: err,
	}

	if err := json.NewEncoder(w).Encode(body); err != nil {
		xlog.Error(err)
	}
}

// JSON response
func JSON(w http.ResponseWriter, code int, body interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		FailureFromError(w, http.StatusInternalServerError, err)
	}
}
