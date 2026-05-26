package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
)

// Values stores flattened configuration keys. Service-specific config structs
// can be built from this map after all sources have been merged.
type Values map[string]string

// Source provides one layer of configuration. Sources loaded later override
// earlier values.
type Source interface {
	Name() string
	Load() (Values, error)
}

// Loader applies sources in the documented order: defaults, TOML, .env,
// operating-system environment, then command-line flags.
type Loader struct {
	sources []Source
}

func NewLoader(sources ...Source) Loader {
	return Loader{sources: sources}
}

func (l Loader) Load() (Values, error) {
	merged := Values{}

	for _, source := range l.sources {
		values, err := source.Load()
		if err != nil {
			return nil, fmt.Errorf("load config source %q: %w", source.Name(), err)
		}
		for key, value := range values {
			if key == "" {
				continue
			}
			merged[key] = value
		}
	}

	return merged, nil
}

type MapSource struct {
	sourceName string
	values     Values
}

func NewMapSource(name string, values Values) MapSource {
	copied := Values{}
	for key, value := range values {
		copied[key] = value
	}
	return MapSource{sourceName: name, values: copied}
}

func (s MapSource) Name() string {
	return s.sourceName
}

func (s MapSource) Load() (Values, error) {
	copied := Values{}
	for key, value := range s.values {
		copied[key] = value
	}
	return copied, nil
}

type EnvFileSource struct {
	path     string
	allowed  map[string]struct{}
	optional bool
}

func NewEnvFileSource(path string, keys []string, optional bool) EnvFileSource {
	return EnvFileSource{path: path, allowed: keySet(keys), optional: optional}
}

func (s EnvFileSource) Name() string {
	return ".env"
}

func (s EnvFileSource) Load() (Values, error) {
	file, err := os.Open(s.path)
	if err != nil {
		if s.optional && errors.Is(err, os.ErrNotExist) {
			return Values{}, nil
		}
		return nil, err
	}
	defer file.Close()

	values := Values{}
	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return nil, fmt.Errorf("line %d missing '='", lineNumber)
		}

		key = strings.TrimSpace(key)
		if !isAllowed(s.allowed, key) {
			continue
		}
		values[key] = strings.TrimSpace(value)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return values, nil
}

type OSEnvSource struct {
	allowed map[string]struct{}
}

func NewOSEnvSource(keys []string) OSEnvSource {
	return OSEnvSource{allowed: keySet(keys)}
}

func (s OSEnvSource) Name() string {
	return "environment"
}

func (s OSEnvSource) Load() (Values, error) {
	values := Values{}
	keys := make([]string, 0, len(s.allowed))
	for key := range s.allowed {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		if value, ok := os.LookupEnv(key); ok {
			values[key] = value
		}
	}
	return values, nil
}

type UnsupportedFileSource struct {
	sourceName string
	path       string
	optional   bool
}

func NewTOMLSource(path string, optional bool) UnsupportedFileSource {
	return UnsupportedFileSource{sourceName: "toml", path: path, optional: optional}
}

func (s UnsupportedFileSource) Name() string {
	return s.sourceName
}

func (s UnsupportedFileSource) Load() (Values, error) {
	if s.optional {
		if _, err := os.Stat(s.path); errors.Is(err, os.ErrNotExist) {
			return Values{}, nil
		}
	}
	return nil, fmt.Errorf("%s parsing is not implemented yet for %q", s.sourceName, s.path)
}

func keySet(keys []string) map[string]struct{} {
	allowed := map[string]struct{}{}
	for _, key := range keys {
		if key == "" {
			continue
		}
		allowed[key] = struct{}{}
	}
	return allowed
}

func isAllowed(allowed map[string]struct{}, key string) bool {
	if len(allowed) == 0 {
		return true
	}
	_, ok := allowed[key]
	return ok
}
