package pushinstallations

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/giantswarm/mcli/pkg/managementcluster/installations"

	"github.com/giantswarm/mcli/pkg/key"
)

const (
	testTeam         = "testteam"
	testDomain       = "test.com"
	testName         = "test"
	testCustomer     = "giantswarm"
	testStableBranch = "stable"
	testingCatalog   = "testing"
	testName2        = "test2"
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
				Team: testTeam,
			},
				BaseDomain:    testDomain,
				CMCRepository: testName},
			expectError: true,
		},
		{
			name: "all flags",
			flags: Config{Flags: InstallationsFlags{
				Team:          testTeam,
				Customer:      testCustomer,
				CCRRepository: testName,
				Pipeline:      testStableBranch,
			},
				BaseDomain:    testDomain,
				CMCRepository: testName,
				Provider:      key.ProviderAzure,
				Cluster:       testName,
			},
			expected: &installations.Installations{
				Base:            testDomain,
				Codename:        testName,
				Customer:        testCustomer,
				Provider:        "capz-test",
				Pipeline:        testStableBranch,
				CmcRepository:   testName,
				CcrRepository:   testName,
				AccountEngineer: testTeam,
			},
			expectError: false,
		},
		{
			name: "all flags with AWS",
			flags: Config{Flags: InstallationsFlags{
				Team:          testTeam,
				Customer:      testName,
				CCRRepository: testName,
				Pipeline:      testingCatalog,
				AWS: AWSFlags{
					Region:                 "eu-west-1",
					InstallationAWSAccount: "123456789012",
				},
			},
				BaseDomain:    testDomain,
				CMCRepository: testName,
				Provider:      key.ProviderAWS,
				Cluster:       testName},
			expected: &installations.Installations{
				Base:            testDomain,
				Codename:        testName,
				Customer:        testName,
				Provider:        "capa-test",
				Pipeline:        testingCatalog,
				CmcRepository:   testName,
				CcrRepository:   testName,
				AccountEngineer: testTeam,
				Aws: installations.AwsConfig{
					Region: "eu-west-1",
					HostCluster: installations.HostCluster{
						Account:          "123456789012",
						AdminRoleArn:     "arn:aws:iam::123456789012:role/GiantSwarmAdmin",
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
				Team:          testTeam,
				CCRRepository: testName,
				Pipeline:      testingCatalog,
			},
				BaseDomain:    testDomain,
				CMCRepository: testName,
				Provider:      key.ProviderAWS,
				Cluster:       testName,
			},
			expectError: true,
		},
		{
			name: "all flags but missing AWS region",
			flags: Config{Flags: InstallationsFlags{
				Team:          testTeam,
				CCRRepository: testName,
				Pipeline:      testingCatalog,
				AWS: AWSFlags{
					InstallationAWSAccount: "123456789012",
				},
			},
				BaseDomain:    testDomain,
				CMCRepository: testName,
				Provider:      key.ProviderAWS,
				Cluster:       testName},
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
				Team: testTeam,
			},
				BaseDomain:    testDomain,
				CMCRepository: testName},
			current: &installations.Installations{},
			expected: &installations.Installations{
				Base:            testDomain,
				CmcRepository:   testName,
				AccountEngineer: testTeam,
			},
		},
		{
			name: "all flags, no current",
			flags: Config{Flags: InstallationsFlags{
				Customer:      testCustomer,
				Team:          testTeam,
				CCRRepository: testName,
			},
				BaseDomain:    testDomain,
				CMCRepository: testName,
				Provider:      key.ProviderAzure,
				Cluster:       testName},
			current: &installations.Installations{},
			expected: &installations.Installations{
				Codename:        testName,
				Base:            testDomain,
				Customer:        testCustomer,
				Provider:        key.ProviderAzure,
				CmcRepository:   testName,
				CcrRepository:   testName,
				AccountEngineer: testTeam,
			},
		},
		{
			name: "some flags, current values are set",
			flags: Config{Flags: InstallationsFlags{
				Team:     testTeam,
				Customer: testCustomer,
			},
				BaseDomain:    testDomain,
				CMCRepository: testName},
			current: &installations.Installations{
				Codename:        testName2,
				Base:            "test2.com",
				CmcRepository:   testName2,
				CcrRepository:   testName2,
				AccountEngineer: "testteam2",
				Customer:        "giantswarm2",
				Provider:        "capv",
				Pipeline:        testStableBranch,
			},
			expected: &installations.Installations{
				Codename:        testName2,
				Base:            testDomain,
				CmcRepository:   testName,
				CcrRepository:   testName2,
				AccountEngineer: testTeam,
				Customer:        testCustomer,
				Provider:        "capv",
				Pipeline:        testStableBranch,
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
