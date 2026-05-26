package logx

import "log/slog"

const (
	FieldServiceName = "service"
	FieldRequestID   = "request_id"
	FieldTraceID     = "trace_id"
	FieldAction      = "action"
	FieldResource    = "resource"
	FieldResult      = "result"
)

type Fields map[string]any

func BaseFields(serviceName string, requestID string, traceID string) Fields {
	return Fields{
		FieldServiceName: serviceName,
		FieldRequestID:   requestID,
		FieldTraceID:     traceID,
	}
}

func Attrs(fields Fields) []slog.Attr {
	attrs := make([]slog.Attr, 0, len(fields))
	for key, value := range fields {
		attrs = append(attrs, slog.Any(key, value))
	}
	return attrs
}
