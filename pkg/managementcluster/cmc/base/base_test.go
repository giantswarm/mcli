package base

import "testing"

func TestGetBaseFiles(t *testing.T) {
	testCases := []struct {
		name     string
		template map[string]string
		config   Config

		expected map[string]string
	}{
		{
			name:     "case 0: empty",
			template: map[string]string{},
			config:   Config{},

			expected: map[string]string{},
		},
		{
			name: "case 1: no variables",
			template: map[string]string{
				"file1": GetTestFile(),
			},
			config: Config{
				CONFIG_BRANCH: "main",
			},

			expected: map[string]string{},
		},
		{
			name: "case 2: one variable",
			template: map[string]string{
				"file1": GetTestTemplate(),
			},
			config: Config{
				CONFIG_BRANCH: "main",
			},

			expected: map[string]string{
				"file1": GetTestFile(),
			},
		},
		{
			name: "case 3: multiple variables",
			template: map[string]string{
				"file1": GetTestTemplate(),
				"file2": GetTestFile(),
				"file3": GetTestTemplate(),
			},
			config: Config{
				CONFIG_BRANCH: "main",
			},

			expected: map[string]string{
				"file1": GetTestFile(),
				"file3": GetTestFile(),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := GetBaseFiles(tc.config, tc.template)
			if err != nil {
				t.Fatalf("error: %v", err)
			}

			if len(actual) != len(tc.expected) {
				t.Fatalf("expected %d files, got %d", len(tc.expected), len(actual))
			}
			for k, v := range actual {
				if tc.expected[k] != v {
					t.Fatalf("expected %s, got %s", tc.expected[k], v)
				}
			}
		})
	}
}

func GetTestTemplate() string {
	return `apiVersion: source.toolkit.fluxcd.io/v1
kind: GitRepository
metadata:
  name: giantswarm-config
  namespace: flux-giantswarm
spec:
  ref:
    branch: ${CONFIG_BRANCH}`
}

func GetTestFile() string {
	return `apiVersion: source.toolkit.fluxcd.io/v1
kind: GitRepository
metadata:
  name: giantswarm-config
  namespace: flux-giantswarm
spec:
  ref:
    branch: main`
}
