package pushinstallations

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/giantswarm/mcli/pkg/managementcluster/installations"
)

func TestGetNewInstallationsFromFlags(t *testing.T) {
	var testCases = []struct {
		name  string
		flags Config

		expected    *installations.Installations
		expectError bool
	}{
		{
			name:        "no flags",
			flags:       Config{Flags: InstallationsFlags{}},
			expectError: true,
		},
		{
			name: "some flags",
			flags: Config{Flags: InstallationsFlags{
				BaseDomain: "test.com",
				Team:       "testteam",
			},
				CMCRepository: "test"},
			expectError: true,
		},
		{
			name: "all flags",
			flags: Config{Flags: InstallationsFlags{
				BaseDomain:    "test.com",
				Team:          "testteam",
				Customer:      "giantswarm",
				CCRRepository: "test",
			},
				CMCRepository: "test",
				Provider:      "capz",
				Cluster:       "test",
			},
			expected: &installations.Installations{
				Base:            "test.com",
				Codename:        "test",
				Customer:        "giantswarm",
				Provider:        "capz-test",
				Pipeline:        "testing",
				CmcRepository:   "test",
				CcrRepository:   "test",
				AccountEngineer: "testteam",
			},
			expectError: false,
		},
		{
			name: "all flags with AWS",
			flags: Config{Flags: InstallationsFlags{
				BaseDomain:    "test.com",
				Team:          "testteam",
				Customer:      "test",
				CCRRepository: "test",
				AWS: AWSFlags{
					Region:                 "eu-west-1",
					InstallationAWSAccount: "123456789012",
				},
			},
				CMCRepository: "test",
				Provider:      "capa",
				Cluster:       "test"},
			expected: &installations.Installations{
				Base:            "test.com",
				Codename:        "test",
				Customer:        "test",
				Provider:        "capa-test",
				Pipeline:        "testing",
				CmcRepository:   "test",
				CcrRepository:   "test",
				AccountEngineer: "testteam",
				Aws: installations.AwsConfig{
					Region: "eu-west-1",
					HostCluster: installations.HostCluster{
						Account:          "123456789012",
						AdminRoleArn:     "arn:aws:iam::123456789012:role/GiantSwarmAdmin",
						CloudtrailBucket: "",
						GuardDuty:        false,
					},
					GuestCluster: installations.GuestCluster{
						Account:          "123456789012",
						CloudtrailBucket: "",
						GuardDuty:        false,
					},
				},
			},
			expectError: false,
		},
		{
			name: "all flags but missing AWS flags",
			flags: Config{Flags: InstallationsFlags{
				BaseDomain:    "test.com",
				Team:          "testteam",
				CCRRepository: "test",
			},
				CMCRepository: "test",
				Provider:      "capa",
				Cluster:       "test",
			},
			expectError: true,
		},
		{
			name: "all flags but missing AWS region",
			flags: Config{Flags: InstallationsFlags{
				BaseDomain:    "test.com",
				Team:          "testteam",
				CCRRepository: "test",
				AWS: AWSFlags{
					InstallationAWSAccount: "123456789012",
				},
			},
				CMCRepository: "test",
				Provider:      "capa",
				Cluster:       "test"},
			expectError: true,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case %d: %s", i, tc.name), func(t *testing.T) {
			installation, err := getNewInstallationsFromFlags(tc.flags)
			if err != nil && !tc.expectError {
				t.Fatalf("unexpected error: %v", err)
			} else if err == nil && tc.expectError {
				t.Fatalf("expected error, got nil")
			}
			if !reflect.DeepEqual(installation, tc.expected) {
				t.Fatalf("expected %#v, got %#v", tc.expected, installation)
			}
		})
	}
}

func TestOverrideInstallationsWithFlags(t *testing.T) {
	var testCases = []struct {
		name    string
		flags   Config
		current *installations.Installations

		expected *installations.Installations
	}{
		{
			name:    "no flags, no current",
			current: &installations.Installations{},
			flags:   Config{Flags: InstallationsFlags{}},

			expected: &installations.Installations{},
		},
		{
			name: "some flags, no current",
			flags: Config{Flags: InstallationsFlags{
				BaseDomain: "test.com",
				Team:       "testteam",
			},
				CMCRepository: "test"},
			current: &installations.Installations{},
			expected: &installations.Installations{
				Base:            "test.com",
				CmcRepository:   "test",
				AccountEngineer: "testteam",
			},
		},
		{
			name: "all flags, no current",
			flags: Config{Flags: InstallationsFlags{
				BaseDomain:    "test.com",
				Customer:      "giantswarm",
				Team:          "testteam",
				CCRRepository: "test",
			},
				CMCRepository: "test",
				Provider:      "capz",
				Cluster:       "test"},
			current: &installations.Installations{},
			expected: &installations.Installations{
				Codename:        "test",
				Base:            "test.com",
				Customer:        "giantswarm",
				Provider:        "capz",
				CmcRepository:   "test",
				CcrRepository:   "test",
				AccountEngineer: "testteam",
			},
		},
		{
			name: "some flags, current values are set",
			flags: Config{Flags: InstallationsFlags{
				BaseDomain: "test.com",
				Team:       "testteam",
				Customer:   "giantswarm",
			},
				CMCRepository: "test"},
			current: &installations.Installations{
				Codename:        "test2",
				Base:            "test2.com",
				CmcRepository:   "test2",
				CcrRepository:   "test2",
				AccountEngineer: "testteam2",
				Customer:        "giantswarm2",
				Provider:        "capv",
				Pipeline:        "stable",
			},
			expected: &installations.Installations{
				Codename:        "test2",
				Base:            "test.com",
				CmcRepository:   "test",
				CcrRepository:   "test2",
				AccountEngineer: "testteam",
				Customer:        "giantswarm",
				Provider:        "capv",
				Pipeline:        "stable",
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case %d: %s", i, tc.name), func(t *testing.T) {
			installation := overrideInstallationsWithFlags(tc.current, tc.flags)
			if !reflect.DeepEqual(installation, tc.expected) {
				t.Fatalf("expected %#v, got %#v", tc.expected, installation)
			}
		})
	}
}
