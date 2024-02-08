package repositories

import (
	"reflect"
	"testing"
)

func TestGetDataFromRepositories(t *testing.T) {
	var testCases = []struct {
		name        string
		repos       []Repo
		expected    string
		expectedErr bool
	}{
		{
			name: "one repository",
			repos: []Repo{
				{
					Name:          "test",
					ComponentType: "configuration",
					Gen: Gen{
						Flavours: []string{"generic"},
						Language: "generic",
					},
					Replace: map[string]bool{
						"dependabotRemove": true,
					},
				},
			},
			expected: "- name: test\n  componentType: configuration\n  gen:\n    flavours:\n      - generic\n    language: generic\n  replace:\n    dependabotRemove: true\n",
		},
		{
			name: "two repositories",
			repos: []Repo{
				{
					Name:          "test",
					ComponentType: "configuration",
					Gen: Gen{
						Flavours: []string{"generic"},
						Language: "generic",
					},
					Replace: map[string]bool{
						"dependabotRemove": true,
					},
				},
				{
					Name:          "test2",
					ComponentType: "configuration",
					Gen: Gen{
						Flavours: []string{"generic"},
						Language: "generic",
					},
					Replace: map[string]bool{
						"dependabotRemove": true,
					},
				},
			},
			expected: "- name: test\n  componentType: configuration\n  gen:\n    flavours:\n      - generic\n    language: generic\n  replace:\n    dependabotRemove: true\n- name: test2\n  componentType: configuration\n  gen:\n    flavours:\n      - generic\n    language: generic\n  replace:\n    dependabotRemove: true\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := GetData(tc.repos)
			if err != nil && !tc.expectedErr {
				t.Fatalf("unexpected error: %v", err)
			} else if err == nil && tc.expectedErr {
				t.Fatalf("expected error, got nil")
			}
			if string(data) != tc.expected {
				t.Fatalf("expected %s, got %s", tc.expected, string(data))
			}
		})
	}
}

func TestGetReposFromData(t *testing.T) {
	var testCases = []struct {
		name        string
		data        string
		expected    []Repo
		expectedErr bool
	}{
		{
			name: "one repository",
			data: "- name: test\n  componentType: configuration\n  gen:\n    flavours:\n      - generic\n    language: generic\n  replace:\n    dependabotRemove: true\n",
			expected: []Repo{
				{
					Name:          "test",
					ComponentType: "configuration",
					Gen: Gen{
						Flavours: []string{"generic"},
						Language: "generic",
					},
					Replace: map[string]bool{
						"dependabotRemove": true,
					},
				},
			},
		},
		{
			name: "two repositories",
			data: "- name: test\n  componentType: configuration\n  gen:\n    flavours:\n      - generic\n    language: generic\n  replace:\n    dependabotRemove: true\n- name: test2\n  componentType: configuration\n  gen:\n    flavours:\n      - generic\n    language: generic\n  replace:\n    dependabotRemove: true\n",
			expected: []Repo{
				{
					Name:          "test",
					ComponentType: "configuration",
					Gen: Gen{
						Flavours: []string{"generic"},
						Language: "generic",
					},
					Replace: map[string]bool{
						"dependabotRemove": true,
					},
				},
				{
					Name:          "test2",
					ComponentType: "configuration",
					Gen: Gen{
						Flavours: []string{"generic"},
						Language: "generic",
					},
					Replace: map[string]bool{
						"dependabotRemove": true,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repos, err := GetRepos([]byte(tc.data))
			if err != nil && !tc.expectedErr {
				t.Fatalf("unexpected error: %v", err)
			} else if err == nil && tc.expectedErr {
				t.Fatalf("expected error, got nil")
			}
			if !reflect.DeepEqual(repos, tc.expected) {
				t.Fatalf("expected %v, got %v", tc.expected, repos)
			}
		})
	}
}
