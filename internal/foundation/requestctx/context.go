package requestctx

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

const (
	HeaderAuthorization = "Authorization"
	HeaderRequestID     = "X-Request-Id"
	HeaderTraceID       = "X-Trace-Id"
	HeaderParentSpanID  = "X-Parent-Span-Id"
	HeaderServiceName   = "X-Service-Name"
)

type IDs struct {
	RequestID string
	TraceID   string
}

func FromRequest(r *http.Request) IDs {
	return IDs{
		RequestID: ensureID(r.Header.Get(HeaderRequestID), "req"),
		TraceID:   ensureID(r.Header.Get(HeaderTraceID), "trace"),
	}
}

func ApplyHeaders(header http.Header, ids IDs) {
	if ids.RequestID != "" {
		header.Set(HeaderRequestID, ids.RequestID)
	}
	if ids.TraceID != "" {
		header.Set(HeaderTraceID, ids.TraceID)
	}
}

func NewID(prefix string) string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return prefix + "_unavailable"
	}
	return prefix + "_" + hex.EncodeToString(b[:])
}

func ensureID(value string, prefix string) string {
	if value != "" {
		return value
	}
	return NewID(prefix)
}
