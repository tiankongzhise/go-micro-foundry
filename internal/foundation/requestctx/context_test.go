package requestctx

import (
	"net/http"
	"strings"
	"testing"
)

func TestFromRequestPreservesProvidedIDs(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://example.test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set(HeaderRequestID, "req_existing")
	req.Header.Set(HeaderTraceID, "trace_existing")

	ids := FromRequest(req)
	if ids.RequestID != "req_existing" || ids.TraceID != "trace_existing" {
		t.Fatalf("ids = %#v", ids)
	}
}

func TestFromRequestGeneratesMissingIDs(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://example.test", nil)
	if err != nil {
		t.Fatal(err)
	}

	ids := FromRequest(req)
	if !strings.HasPrefix(ids.RequestID, "req_") {
		t.Fatalf("RequestID = %q", ids.RequestID)
	}
	if !strings.HasPrefix(ids.TraceID, "trace_") {
		t.Fatalf("TraceID = %q", ids.TraceID)
	}
}
