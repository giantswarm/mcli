package key

import "testing"

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
