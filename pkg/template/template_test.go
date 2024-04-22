package template

import (
	"testing"
)

func TestExecute(t *testing.T) {
	var testCases = []struct {
		name     string
		tmplFile string
		data     any

		expected    string
		expectError bool
	}{
		{
			name:        "case 0",
			tmplFile:    "Hello {{ .Name }}!\n",
			data:        map[string]string{"Name": "John Doe"},
			expected:    "Hello John Doe!\n",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := Execute(tc.tmplFile, tc.data)

			if tc.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if actual != tc.expected {
					t.Fatalf("expected %s, got %s", tc.expected, actual)
				}
			}
		})
	}
}
