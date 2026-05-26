package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoaderAppliesSourcesInOrder(t *testing.T) {
	t.Setenv("SERVICE_PORT", "18030")

	envPath := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(envPath, []byte("SERVICE_PORT=18020\nSERVICE_NAME=from-env-file\nIGNORED=value\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	loader := NewLoader(
		NewMapSource("defaults", Values{"SERVICE_PORT": "18010", "SERVICE_MODE": "prod"}),
		NewEnvFileSource(envPath, []string{"SERVICE_PORT", "SERVICE_NAME"}, false),
		NewOSEnvSource([]string{"SERVICE_PORT"}),
		NewMapSource("flags", Values{"SERVICE_MODE": "dev"}),
	)

	values, err := loader.Load()
	if err != nil {
		t.Fatal(err)
	}

	if got := values["SERVICE_PORT"]; got != "18030" {
		t.Fatalf("SERVICE_PORT = %q, want %q", got, "18030")
	}
	if got := values["SERVICE_MODE"]; got != "dev" {
		t.Fatalf("SERVICE_MODE = %q, want %q", got, "dev")
	}
	if got := values["SERVICE_NAME"]; got != "from-env-file" {
		t.Fatalf("SERVICE_NAME = %q, want %q", got, "from-env-file")
	}
	if _, ok := values["IGNORED"]; ok {
		t.Fatal("unexpected unlisted env key")
	}
}

func TestOptionalEnvFileMissing(t *testing.T) {
	loader := NewLoader(NewEnvFileSource(filepath.Join(t.TempDir(), ".env"), nil, true))

	values, err := loader.Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(values) != 0 {
		t.Fatalf("values = %#v, want empty", values)
	}
}
