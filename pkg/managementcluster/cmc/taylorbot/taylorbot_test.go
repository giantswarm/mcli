package taylorbot

import "testing"

func TestGetTaylorBotToken(t *testing.T) {
	testCases := []struct {
		name string
		file string

		expectErr    bool
		expectOutput string
	}{
		{
			name:         "valid file",
			file:         "apiVersion: v1\ndata:\n    password: somethingsomething\n    url: xyz\n    username: abc\nkind: Secret\nmetadata:\n    creationTimestamp: null\n    name: github-giantswarm-https-credentials\n    namespace: flux-giantswarm\n",
			expectOutput: "somethingsomething",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := GetTaylorBotToken(tc.file)
			if tc.expectErr {
				if err == nil {
					t.Fatalf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error but got %v", err)
				}
				if output != tc.expectOutput {
					t.Fatalf("expected %q but got %q", tc.expectOutput, output)
				}
			}
		})
	}
}
