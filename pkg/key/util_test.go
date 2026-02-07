package key

import (
	"strings"
	"testing"
)

func TestGetSchemaHeader(t *testing.T) {
	testCases := []struct {
		name       string
		schemaPath string
		expected   string
	}{
		{
			name:       "relative path",
			schemaPath: "../schema.json",
			expected:   "# yaml-language-server: $schema=../schema.json\n",
		},
		{
			name:       "absolute path",
			schemaPath: "/path/to/schema.json",
			expected:   "# yaml-language-server: $schema=/path/to/schema.json\n",
		},
		{
			name:       "url",
			schemaPath: "https://example.com/schema.json",
			expected:   "# yaml-language-server: $schema=https://example.com/schema.json\n",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := GetSchemaHeader(tc.schemaPath)
			if actual != tc.expected {
				t.Fatalf("expected %q but got %q", tc.expected, actual)
			}
		})
	}
}

func TestPrependSchemaHeader(t *testing.T) {
	testCases := []struct {
		name       string
		data       []byte
		schemaPath string
		expected   string
	}{
		{
			name:       "prepend to yaml",
			data:       []byte("key: value\n"),
			schemaPath: "../schema.json",
			expected:   "# yaml-language-server: $schema=../schema.json\nkey: value\n",
		},
		{
			name:       "prepend to empty data",
			data:       []byte{},
			schemaPath: "../schema.json",
			expected:   "# yaml-language-server: $schema=../schema.json\n",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := PrependSchemaHeader(tc.data, tc.schemaPath)
			if string(actual) != tc.expected {
				t.Fatalf("expected %q but got %q", tc.expected, string(actual))
			}
		})
	}
}

func TestPrependSchemaHeaderStartsWithComment(t *testing.T) {
	data := []byte("key: value\n")
	result := PrependSchemaHeader(data, "../schema.json")
	if !strings.HasPrefix(string(result), "#") {
		t.Fatal("expected result to start with #")
	}
}

func TestGetSecretValue(t *testing.T) {
	testCases := []struct {
		name     string
		key      string
		data     string
		expected string
	}{
		{
			name:     "simple",
			key:      "key",
			data:     "key: value",
			expected: "value",
		},
		{
			name:     "multiline",
			key:      "key",
			data:     "key: |\n  value\n  value",
			expected: "value\nvalue",
		},
		{
			name:     "more keys",
			key:      "key",
			data:     "key: value\nkey2: value2",
			expected: "value",
		},
		{
			name:     "more keys multiline",
			key:      "key",
			data:     "key: |\n  value\n  value\nkey2: value2",
			expected: "value\nvalue",
		},
		{
			name:     "more keys multiline with other keys",
			key:      "key",
			data:     "key: |\n  value\n  value\nkey2: value2\nkey3: |\n  value\n  value",
			expected: "value\nvalue",
		},
		{
			name:     "base64 encoded",
			key:      "key",
			data:     "key: dmFsdWU=",
			expected: "value",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := GetSecretValue(tc.key, tc.data)
			if err != nil {
				t.Fatalf("expected no error but got %v", err)
			}
			if actual != tc.expected {
				t.Fatalf("expected %q but got %q", tc.expected, actual)
			}
		})
	}
}
