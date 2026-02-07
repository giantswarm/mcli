package apps

import (
	"reflect"
	"testing"

	"github.com/giantswarm/mcli/pkg/key"
)

func TestGetAppsConfig(t *testing.T) {
	testCases := []struct {
		name string
		file string

		expected    Config
		expectError bool
	}{
		{
			name: "case 0: simple",
			file: getValidClusterApps(),

			expected: Config{
				Name:                         "test",
				Namespace:                    "test-namespace",
				Provider:                     key.ProviderAWS,
				AppName:                      "cluster-aws",
				Catalog:                      "cluster",
				Version:                      "0.1.0",
				Values:                       "global:\n  controlPlane:\n    some: value\n    and:\n      another: value\n    this:\n    - is\n    - a\n    - list\n    oidc:\n      issuerUrl: https://dex.gigantic.io/dex\n      clientId: dex-k8s-authenticator\n      usernameClaim: email\n      groupsClaim: groups\n  providerSpecific:\n    awsClusterRoleIdentityName: \"default\"\n    region: \"eu-west-2\"\n  metadata:\n    preventDeletion: true\n    description: \"test MC\"\n    name: \"test\"\n    organization: \"giantswarm\"\n  nodePools:\n    hello:\n      availabilityZones:\n      - eu-west-2a\n      instanceType: r6i.xlarge\n      minSize: 3\n      maxSize: 6\n      rootVolumeSizeGB: 300\n      customNodeLabels:\n      - label=test\n  podSecurityStandards:\n    enforced: true\n",
				ConfigureContainerRegistries: false,
				MCAppsPreventDeletion:        false,
			},
		},
		{
			name: "case 1: prevent deletion",
			file: getValidClusterAppsPreventDeletion(),

			expected: Config{
				Name:                         "test",
				Namespace:                    "test-namespace",
				Provider:                     key.ProviderAWS,
				AppName:                      "cluster-aws",
				Catalog:                      "cluster",
				Version:                      "0.1.0",
				Values:                       "global:\n  controlPlane:\n    some: value\n    and:\n      another: value\n    this:\n    - is\n    - a\n    - list\n    oidc:\n      issuerUrl: https://dex.gigantic.io/dex\n      clientId: dex-k8s-authenticator\n      usernameClaim: email\n      groupsClaim: groups\n  providerSpecific:\n    awsClusterRoleIdentityName: \"default\"\n    region: \"eu-west-2\"\n  metadata:\n    preventDeletion: true\n    description: \"test MC\"\n    name: \"test\"\n    organization: \"giantswarm\"\n  nodePools:\n    hello:\n      availabilityZones:\n      - eu-west-2a\n      instanceType: r6i.xlarge\n      minSize: 3\n      maxSize: 6\n      rootVolumeSizeGB: 300\n      customNodeLabels:\n      - label=test\n  podSecurityStandards:\n    enforced: true\n",
				ConfigureContainerRegistries: false,
				MCAppsPreventDeletion:        true,
			},
		},
		{
			name: "case 2: configure container registries",
			file: getValidClusterAppsConfigureContainerRegistries(),

			expected: Config{
				Name:                         "test",
				Namespace:                    "test-namespace",
				Provider:                     key.ProviderAWS,
				AppName:                      "cluster-aws",
				Catalog:                      "cluster",
				Version:                      "0.1.0",
				Values:                       "global:\n  controlPlane:\n    some: value\n    and:\n      another: value\n    this:\n    - is\n    - a\n    - list\n    oidc:\n      issuerUrl: https://dex.gigantic.io/dex\n      clientId: dex-k8s-authenticator\n      usernameClaim: email\n      groupsClaim: groups\n  providerSpecific:\n    awsClusterRoleIdentityName: \"default\"\n    region: \"eu-west-2\"\n  metadata:\n    preventDeletion: true\n    description: \"test MC\"\n    name: \"test\"\n    organization: \"giantswarm\"\n  nodePools:\n    hello:\n      availabilityZones:\n      - eu-west-2a\n      instanceType: r6i.xlarge\n      minSize: 3\n      maxSize: 6\n      rootVolumeSizeGB: 300\n      customNodeLabels:\n      - label=test\n  podSecurityStandards:\n    enforced: true\n",
				ConfigureContainerRegistries: true,
				MCAppsPreventDeletion:        false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := GetAppsConfig(tc.file)
			if tc.expectError {
				if err == nil {
					t.Fatalf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error but got %v", err)
				}
				if !reflect.DeepEqual(actual, tc.expected) {
					t.Fatalf("expected %#v but got %#v", tc.expected, actual)
				}
			}
		})
	}
}

func TestGetClusterAppsFile(t *testing.T) {
	testCases := []struct {
		name string
		c    Config

		expected    string
		expectError bool
	}{
		{
			name: "case 0: simple",
			c: Config{
				Cluster:                      "test",
				Name:                         "test",
				Namespace:                    "test-namespace",
				Provider:                     key.ProviderAWS,
				AppName:                      "cluster-aws",
				Catalog:                      "cluster",
				Version:                      "0.1.0",
				ConfigureContainerRegistries: false,
				MCAppsPreventDeletion:        false,
				Values:                       "global:\n  controlPlane:\n    some: value\n    and:\n      another: value\n    this:\n    - is\n    - a\n    - list\n    oidc:\n      issuerUrl: https://dex.gigantic.io/dex\n      clientId: dex-k8s-authenticator\n      usernameClaim: email\n      groupsClaim: groups\n  providerSpecific:\n    awsClusterRoleIdentityName: \"default\"\n    region: \"eu-west-2\"\n  metadata:\n    preventDeletion: true\n    description: \"test MC\"\n    name: \"test\"\n    organization: \"giantswarm\"\n  nodePools:\n    hello:\n      availabilityZones:\n      - eu-west-2a\n      instanceType: r6i.xlarge\n      minSize: 3\n      maxSize: 6\n      rootVolumeSizeGB: 300\n      customNodeLabels:\n      - label=test\n  podSecurityStandards:\n    enforced: true\n",
			},

			expected: getValidClusterApps(),
		},
		{
			name: "case 1: prevent deletion",
			c: Config{
				Cluster:                      "test",
				Name:                         "test",
				Namespace:                    "test-namespace",
				Provider:                     key.ProviderAWS,
				AppName:                      "cluster-aws",
				Catalog:                      "cluster",
				Version:                      "0.1.0",
				ConfigureContainerRegistries: false,
				MCAppsPreventDeletion:        true,
				Values:                       "global:\n  controlPlane:\n    some: value\n    and:\n      another: value\n    this:\n    - is\n    - a\n    - list\n    oidc:\n      issuerUrl: https://dex.gigantic.io/dex\n      clientId: dex-k8s-authenticator\n      usernameClaim: email\n      groupsClaim: groups\n  providerSpecific:\n    awsClusterRoleIdentityName: \"default\"\n    region: \"eu-west-2\"\n  metadata:\n    preventDeletion: true\n    description: \"test MC\"\n    name: \"test\"\n    organization: \"giantswarm\"\n  nodePools:\n    hello:\n      availabilityZones:\n      - eu-west-2a\n      instanceType: r6i.xlarge\n      minSize: 3\n      maxSize: 6\n      rootVolumeSizeGB: 300\n      customNodeLabels:\n      - label=test\n  podSecurityStandards:\n    enforced: true\n",
			},
			expected: getValidClusterAppsPreventDeletion(),
		},
		{
			name: "case 2: configure container registries",
			c: Config{
				Cluster:                      "test",
				Name:                         "test",
				Namespace:                    "test-namespace",
				Provider:                     key.ProviderAWS,
				AppName:                      "cluster-aws",
				Catalog:                      "cluster",
				Version:                      "0.1.0",
				ConfigureContainerRegistries: true,
				MCAppsPreventDeletion:        false,
				Values:                       "global:\n  controlPlane:\n    some: value\n    and:\n      another: value\n    this:\n    - is\n    - a\n    - list\n    oidc:\n      issuerUrl: https://dex.gigantic.io/dex\n      clientId: dex-k8s-authenticator\n      usernameClaim: email\n      groupsClaim: groups\n  providerSpecific:\n    awsClusterRoleIdentityName: \"default\"\n    region: \"eu-west-2\"\n  metadata:\n    preventDeletion: true\n    description: \"test MC\"\n    name: \"test\"\n    organization: \"giantswarm\"\n  nodePools:\n    hello:\n      availabilityZones:\n      - eu-west-2a\n      instanceType: r6i.xlarge\n      minSize: 3\n      maxSize: 6\n      rootVolumeSizeGB: 300\n      customNodeLabels:\n      - label=test\n  podSecurityStandards:\n    enforced: true\n",
			},
			expected: getValidClusterAppsConfigureContainerRegistries(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := GetClusterAppsFile(tc.c)
			if tc.expectError {
				if err == nil {
					t.Fatalf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error but got %v", err)
				}
				if actual != tc.expected {
					t.Fatalf("expected %s but got %s", tc.expected, actual)
				}
			}
		})
	}
}

func TestGetDefaultAppsFile(t *testing.T) {
	testCases := []struct {
		name string
		c    Config

		expected    string
		expectError bool
	}{
		{
			name: "case 0: simple",
			c: Config{
				Cluster:               "test",
				Name:                  "test-default-apps",
				Namespace:             "test-namespace",
				Provider:              key.ProviderAWS,
				AppName:               "default-apps-aws",
				Catalog:               "cluster",
				Version:               "0.1.0",
				MCAppsPreventDeletion: false,
				Values:                "clusterName: test\norganization: giantswarm\nmanagementCluster: test\n",
			},

			expected: getValidDefaultApps(),
		},
		{
			name: "case 1: prevent deletion",
			c: Config{
				Cluster:               "test",
				Name:                  "test-default-apps",
				Namespace:             "test-namespace",
				Provider:              key.ProviderAWS,
				AppName:               "default-apps-aws",
				Catalog:               "cluster",
				Version:               "0.1.0",
				MCAppsPreventDeletion: true,
				Values:                "clusterName: test\norganization: giantswarm\nmanagementCluster: test\n",
			},

			expected: getValidDefaultAppsPreventDeletion(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := GetDefaultAppsFile(tc.c)
			if tc.expectError {
				if err == nil {
					t.Fatalf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error but got %v", err)
				}
				if actual != tc.expected {
					t.Fatalf("expected %s but got %s", tc.expected, actual)
				}
			}
		})
	}
}

func getValidDefaultApps() string {
	return `---
apiVersion: v1
data:
  values: |
    clusterName: test
    organization: giantswarm
    managementCluster: test
kind: ConfigMap
metadata:
  name: test-default-apps-user-values
  namespace: test-namespace
---
apiVersion: application.giantswarm.io/v1alpha1
kind: App
metadata:
  labels:
    app-operator.giantswarm.io/version: 0.0.0
    giantswarm.io/cluster: test
    giantswarm.io/managed-by: cluster
  name: test-default-apps
  namespace: test-namespace
spec:
  catalog: cluster
  kubeConfig:
    inCluster: true
  name: default-apps-aws
  namespace: test-namespace
  userConfig:
    configMap:
      name: test-default-apps-user-values
      namespace: test-namespace
  version: 0.1.0
  config:
    configMap:
      name: test-cluster-values
      namespace: test-namespace
`
}

func getValidDefaultAppsPreventDeletion() string {
	return `---
apiVersion: v1
data:
  values: |
    clusterName: test
    organization: giantswarm
    managementCluster: test
kind: ConfigMap
metadata:
  labels:
    giantswarm.io/prevent-deletion: "true"
  name: test-default-apps-user-values
  namespace: test-namespace
---
apiVersion: application.giantswarm.io/v1alpha1
kind: App
metadata:
  labels:
    app-operator.giantswarm.io/version: 0.0.0
    giantswarm.io/cluster: test
    giantswarm.io/managed-by: cluster
    giantswarm.io/prevent-deletion: "true"
  name: test-default-apps
  namespace: test-namespace
spec:
  catalog: cluster
  kubeConfig:
    inCluster: true
  name: default-apps-aws
  namespace: test-namespace
  userConfig:
    configMap:
      name: test-default-apps-user-values
      namespace: test-namespace
  version: 0.1.0
  config:
    configMap:
      name: test-cluster-values
      namespace: test-namespace
`
}

func getValidClusterApps() string {
	return `---
apiVersion: v1
data:
  values: |
    global:
      controlPlane:
        some: value
        and:
          another: value
        this:
        - is
        - a
        - list
        oidc:
          issuerUrl: https://dex.gigantic.io/dex
          clientId: dex-k8s-authenticator
          usernameClaim: email
          groupsClaim: groups
      providerSpecific:
        awsClusterRoleIdentityName: "default"
        region: "eu-west-2"
      metadata:
        preventDeletion: true
        description: "test MC"
        name: "test"
        organization: "giantswarm"
      nodePools:
        hello:
          availabilityZones:
          - eu-west-2a
          instanceType: r6i.xlarge
          minSize: 3
          maxSize: 6
          rootVolumeSizeGB: 300
          customNodeLabels:
          - label=test
      podSecurityStandards:
        enforced: true
kind: ConfigMap
metadata:
  name: test-user-values
  namespace: test-namespace
---
apiVersion: application.giantswarm.io/v1alpha1
kind: App
metadata:
  labels:
    app-operator.giantswarm.io/version: 0.0.0
    giantswarm.io/cluster: test
  name: test
  namespace: test-namespace
spec:
  catalog: cluster
  kubeConfig:
    inCluster: true
  name: cluster-aws
  namespace: test-namespace
  userConfig:
    configMap:
      name: test-user-values
      namespace: test-namespace
  version: 0.1.0
`
}

func getValidClusterAppsPreventDeletion() string {
	return `---
apiVersion: v1
data:
  values: |
    global:
      controlPlane:
        some: value
        and:
          another: value
        this:
        - is
        - a
        - list
        oidc:
          issuerUrl: https://dex.gigantic.io/dex
          clientId: dex-k8s-authenticator
          usernameClaim: email
          groupsClaim: groups
      providerSpecific:
        awsClusterRoleIdentityName: "default"
        region: "eu-west-2"
      metadata:
        preventDeletion: true
        description: "test MC"
        name: "test"
        organization: "giantswarm"
      nodePools:
        hello:
          availabilityZones:
          - eu-west-2a
          instanceType: r6i.xlarge
          minSize: 3
          maxSize: 6
          rootVolumeSizeGB: 300
          customNodeLabels:
          - label=test
      podSecurityStandards:
        enforced: true
kind: ConfigMap
metadata:
  labels:
    giantswarm.io/prevent-deletion: "true"
  name: test-user-values
  namespace: test-namespace
---
apiVersion: application.giantswarm.io/v1alpha1
kind: App
metadata:
  labels:
    app-operator.giantswarm.io/version: 0.0.0
    giantswarm.io/cluster: test
    giantswarm.io/prevent-deletion: "true"
  name: test
  namespace: test-namespace
spec:
  catalog: cluster
  kubeConfig:
    inCluster: true
  name: cluster-aws
  namespace: test-namespace
  userConfig:
    configMap:
      name: test-user-values
      namespace: test-namespace
  version: 0.1.0
`
}

func getValidClusterAppsConfigureContainerRegistries() string {
	return `---
apiVersion: v1
data:
  values: |
    global:
      controlPlane:
        some: value
        and:
          another: value
        this:
        - is
        - a
        - list
        oidc:
          issuerUrl: https://dex.gigantic.io/dex
          clientId: dex-k8s-authenticator
          usernameClaim: email
          groupsClaim: groups
      providerSpecific:
        awsClusterRoleIdentityName: "default"
        region: "eu-west-2"
      metadata:
        preventDeletion: true
        description: "test MC"
        name: "test"
        organization: "giantswarm"
      nodePools:
        hello:
          availabilityZones:
          - eu-west-2a
          instanceType: r6i.xlarge
          minSize: 3
          maxSize: 6
          rootVolumeSizeGB: 300
          customNodeLabels:
          - label=test
      podSecurityStandards:
        enforced: true
kind: ConfigMap
metadata:
  name: test-user-values
  namespace: test-namespace
---
apiVersion: application.giantswarm.io/v1alpha1
kind: App
metadata:
  labels:
    app-operator.giantswarm.io/version: 0.0.0
    giantswarm.io/cluster: test
  name: test
  namespace: test-namespace
spec:
  catalog: cluster
  extraConfigs:
  - kind: secret
    name: container-registries-configuration
    namespace: default
    priority: 0
  kubeConfig:
    inCluster: true
  name: cluster-aws
  namespace: test-namespace
  userConfig:
    configMap:
      name: test-user-values
      namespace: test-namespace
  version: 0.1.0
`
}
