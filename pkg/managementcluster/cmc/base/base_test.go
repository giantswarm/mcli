package base

import "testing"

const (
	path = "management-clusters/test/"
)

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
				ConfigBranch: "auto",
			},

			expected: map[string]string{},
		},
		{
			name: "case 2: one variable",
			template: map[string]string{
				"file1": GetTestTemplate(),
			},
			config: Config{
				ConfigBranch: "auto",
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
				ConfigBranch: "auto",
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

func TestGetBaseConfig(t *testing.T) {
	testCases := []struct {
		name     string
		template map[string]string

		expected  Config
		expectErr bool
	}{
		{
			name:     "case 0: empty",
			template: map[string]string{},

			expectErr: true,
		},
		{
			name: "case 1: catalog patch",
			template: map[string]string{
				path + catalogPatchFile: GetTestPatch(),
			},

			expected: Config{
				CMCBranch:             "main",
				MCAppCollectionBranch: "main",
				BaseDomain:            "test.base.domain.io",
				RegistryDomain:        "registry.domain.test.io",
				MCBBranchSource:       "main",
				ConfigBranch:          "main",
			},
		},
		{
			name: "case 2: cmc branch",
			template: map[string]string{
				path + customBranchCMCFile: GetTestFile(),
				path + catalogPatchFile:    GetTestPatch(),
			},

			expected: Config{
				CMCBranch:             "auto",
				MCAppCollectionBranch: "main",
				BaseDomain:            "test.base.domain.io",
				RegistryDomain:        "registry.domain.test.io",
				MCBBranchSource:       "main",
				ConfigBranch:          "main",
			},
		},
		{
			name: "case 3: config branch",
			template: map[string]string{
				path + customBranchConfigFile: GetTestFile(),
				path + catalogPatchFile:       GetTestPatch(),
			},

			expected: Config{
				CMCBranch:             "main",
				MCAppCollectionBranch: "main",
				BaseDomain:            "test.base.domain.io",
				RegistryDomain:        "registry.domain.test.io",
				MCBBranchSource:       "main",
				ConfigBranch:          "auto",
			},
		},
		{
			name: "case 4: mc app collection branch",
			template: map[string]string{
				path + customBranchCollectionFile: GetTestFile(),
				path + catalogPatchFile:           GetTestPatch(),
			},

			expected: Config{
				CMCBranch:             "main",
				MCAppCollectionBranch: "auto",
				BaseDomain:            "test.base.domain.io",
				RegistryDomain:        "registry.domain.test.io",
				MCBBranchSource:       "main",
				ConfigBranch:          "main",
			},
		},
		{
			name: "case 5: mcb branch source",
			template: map[string]string{
				path + catalogKustomizationFile: GetTestKustomization(),
				path + catalogPatchFile:         GetTestPatch(),
			},

			expected: Config{
				CMCBranch:             "main",
				MCAppCollectionBranch: "main",
				BaseDomain:            "test.base.domain.io",
				RegistryDomain:        "registry.domain.test.io",
				MCBBranchSource:       "hello",
				ConfigBranch:          "main",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := GetBaseConfig(tc.template, "management-clusters/test")
			if err != nil {
				if !tc.expectErr {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if tc.expectErr {
				t.Fatalf("expected error, got nil")
			}

			if actual != tc.expected {
				t.Fatalf("expected %v, got %v", tc.expected, actual)
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
    branch: auto`
}

func GetTestPatch() string {
	return `apiVersion: helm.toolkit.fluxcd.io/v2beta1
kind: HelmRelease
metadata:
  name: appcatalog-default
  namespace: flux-giantswarm
spec:
  values:
    appCatalog:
      config:
        configMap:
          values:
            baseDomain: test.base.domain.io
            managementCluster: test
            provider: capa
            image:
              registry: registry.domain.test.io
            registry:
              domain: registry.domain.test.io`
}

func GetTestKustomization() string {
	return `apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
patches:
  - path: patches/appcatalog-default-patch.yaml
resources:
  - https://github.com/giantswarm/management-cluster-bases//bases/catalogs?ref=hello
  - another-resource.yaml`
}
