package httpx

import "testing"

func TestSuccessEnvelope(t *testing.T) {
	envelope := Success(map[string]string{"id": "cfg_001"}, "req_1", "trace_1")

	if envelope.Code != CodeOK {
		t.Fatalf("Code = %q, want %q", envelope.Code, CodeOK)
	}
	if envelope.Message != "success" {
		t.Fatalf("Message = %q, want success", envelope.Message)
	}
	if envelope.RequestID != "req_1" || envelope.TraceID != "trace_1" {
		t.Fatalf("ids = %q/%q", envelope.RequestID, envelope.TraceID)
	}
}

func TestFailureEnvelopeHasNullData(t *testing.T) {
	envelope := Failure(CodeBadRequest, "bad request", "req_1", "trace_1")

	if envelope.Data != nil {
		t.Fatalf("Data = %#v, want nil", envelope.Data)
	}
	if envelope.Code != CodeBadRequest {
		t.Fatalf("Code = %q, want %q", envelope.Code, CodeBadRequest)
	}
}
