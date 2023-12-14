package pushinstallations

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/giantswarm/mcli/pkg/managementcluster/installations"
)

func TestGetNewInstallationsFromFlags(t *testing.T) {
	var testCases = []struct {
		name    string
		flags   InstallationsFlags
		cluster string

		expected    *installations.Installations
		expectError bool
	}{
		{
			name:        "no flags",
			cluster:     "test",
			flags:       InstallationsFlags{},
			expectError: true,
		},
		{
			name: "some flags",
			flags: InstallationsFlags{
				BaseDomain:    "test.com",
				CMCRepository: "test",
				Team:          "testteam",
			},
			cluster:     "test",
			expectError: true,
		},
		{
			name: "all flags",
			flags: InstallationsFlags{
				BaseDomain:    "test.com",
				CMCRepository: "test",
				Team:          "testteam",
				Provider:      "capz",
			},
			cluster: "test",
			expected: &installations.Installations{
				Base:            "test.com",
				Codename:        "test",
				Customer:        "giantswarm",
				Provider:        "capz-test",
				Pipeline:        "testing",
				CmcRepository:   "test",
				AccountEngineer: "testteam",
			},
			expectError: false,
		},
		{
			name: "all flags with AWS",
			flags: InstallationsFlags{
				BaseDomain:    "test.com",
				CMCRepository: "test",
				Team:          "testteam",
				Provider:      "capa",
				AWS: AWSFlags{
					Region:                 "eu-west-1",
					InstallationAWSAccount: "123456789012",
				},
			},
			cluster: "test",
			expected: &installations.Installations{
				Base:            "test.com",
				Codename:        "test",
				Customer:        "giantswarm",
				Provider:        "capa-test",
				Pipeline:        "testing",
				CmcRepository:   "test",
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
			flags: InstallationsFlags{
				BaseDomain:    "test.com",
				CMCRepository: "test",
				Team:          "testteam",
				Provider:      "capa",
			},
			cluster:     "test",
			expectError: true,
		},
		{
			name: "all flags but missing AWS region",
			flags: InstallationsFlags{
				BaseDomain:    "test.com",
				CMCRepository: "test",
				Team:          "testteam",
				Provider:      "capa",
				AWS: AWSFlags{
					InstallationAWSAccount: "123456789012",
				},
			},
			cluster:     "test",
			expectError: true,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case %d: %s", i, tc.name), func(t *testing.T) {
			installation, err := getNewInstallationsFromFlags(tc.flags, tc.cluster)
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
		flags   InstallationsFlags
		current *installations.Installations

		expected *installations.Installations
	}{
		{
			name:    "no flags, no current",
			current: &installations.Installations{},
			flags:   InstallationsFlags{},

			expected: &installations.Installations{},
		},
		{
			name: "some flags, no current",
			flags: InstallationsFlags{
				BaseDomain:    "test.com",
				CMCRepository: "test",
				Team:          "testteam",
			},
			current: &installations.Installations{},
			expected: &installations.Installations{
				Base:            "test.com",
				CmcRepository:   "test",
				AccountEngineer: "testteam",
			},
		},
		{
			name: "all flags, no current",
			flags: InstallationsFlags{
				BaseDomain:    "test.com",
				CMCRepository: "test",
				Team:          "testteam",
				Provider:      "capz",
			},
			current: &installations.Installations{},
			expected: &installations.Installations{
				Base:            "test.com",
				Customer:        "giantswarm",
				Provider:        "capz",
				Pipeline:        "testing",
				CmcRepository:   "test",
				AccountEngineer: "testteam",
			},
		},
		{
			name: "some flags, current values are set",
			flags: InstallationsFlags{
				BaseDomain:    "test.com",
				CMCRepository: "test",
				Team:          "testteam",
			},
			current: &installations.Installations{
				Base:            "test2.com",
				CmcRepository:   "test2",
				AccountEngineer: "testteam2",
				Customer:        "giantswarm2",
				Provider:        "capv",
				Pipeline:        "stable",
			},
			expected: &installations.Installations{
				Base:            "test.com",
				CmcRepository:   "test",
				AccountEngineer: "testteam",
				Customer:        "giantswarm2",
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
