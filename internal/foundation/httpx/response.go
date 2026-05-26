package httpx

import (
	"encoding/json"
	"net/http"
)

const (
	CodeOK            = "OK"
	CodeBadRequest    = "COMMON_BAD_REQUEST"
	CodeUnauthorized  = "COMMON_UNAUTHORIZED"
	CodeForbidden     = "COMMON_FORBIDDEN"
	CodeNotFound      = "COMMON_NOT_FOUND"
	CodeConflict      = "COMMON_CONFLICT"
	CodeTooMany       = "COMMON_TOO_MANY_REQUESTS"
	CodeInternalError = "COMMON_INTERNAL_ERROR"
)

type Envelope struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Data      any    `json:"data"`
	RequestID string `json:"request_id"`
	TraceID   string `json:"trace_id"`
}

type Page[T any] struct {
	Items    []T `json:"items"`
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}

func Success(data any, requestID string, traceID string) Envelope {
	return Envelope{
		Code:      CodeOK,
		Message:   "success",
		Data:      data,
		RequestID: requestID,
		TraceID:   traceID,
	}
}

func Failure(code string, message string, requestID string, traceID string) Envelope {
	return Envelope{
		Code:      code,
		Message:   message,
		Data:      nil,
		RequestID: requestID,
		TraceID:   traceID,
	}
}

func WriteJSON(w http.ResponseWriter, status int, envelope Envelope) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(envelope)
}
