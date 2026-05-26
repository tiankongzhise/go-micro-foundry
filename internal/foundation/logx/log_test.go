package logx

import "testing"

func TestBaseFields(t *testing.T) {
	fields := BaseFields("config-center", "req_1", "trace_1")

	if fields[FieldServiceName] != "config-center" {
		t.Fatalf("service = %#v", fields[FieldServiceName])
	}
	if fields[FieldRequestID] != "req_1" || fields[FieldTraceID] != "trace_1" {
		t.Fatalf("ids = %#v/%#v", fields[FieldRequestID], fields[FieldTraceID])
	}
}
