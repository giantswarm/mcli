package cmc

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"filippo.io/age"

	"github.com/giantswarm/mcli/pkg/key"
	"github.com/giantswarm/mcli/pkg/sops"
)

func TestGetMapFromCMC(t *testing.T) {
	var testCases = []struct {
		name        string
		cmc         *CMC
		template    map[string]string
		expectError bool
	}{
		{
			name: "case 0: valid aws CMC",
			cmc: &CMC{
				Cluster:    "cluster",
				BaseDomain: "basedomain.io",
				GitOps: GitOps{
					CMCRepository:         "test-management-clusters",
					CMCBranch:             "cmc-branch",
					MCBBranchSource:       "mcb-branch",
					ConfigBranch:          "config-branch",
					MCAppCollectionBranch: "mc-app-collection-branch",
				},
				ClusterApp: App{
					Name:    "clusterapp-aws",
					Values:  "global:\n  clusterapp: values",
					Version: "clusterappversion",
					Catalog: "clustercatalog",
					AppName: "clusterappname-aws",
				},
				DefaultApps: App{
					Name:    "defaultapp-aws",
					Values:  "defaultappvalues",
					Version: "defaultappversion",
					Catalog: "defaultcatalog",
					AppName: "defaultappname-aws",
				},
				MCAppsPreventDeletion: true,
				PrivateCA:             true,
				ClusterNamespace:      "clusternamespace",
				Provider: Provider{
					Name: key.ProviderAWS,
				},
				TaylorBotToken: "taylorbottoken",
				SSHdeployKey: DeployKey{
					Identity:   "identity",
					Passphrase: "passphrase",
					KnownHosts: "knownhosts",
				},
				CustomerDeployKey: DeployKey{
					Identity:   "customeridentity",
					Passphrase: "customerpassphrase",
					KnownHosts: "customerknownhosts",
				},
				SharedDeployKey: DeployKey{
					Identity:   "sharedidentity",
					Passphrase: "sharedpassphrase",
					KnownHosts: "sharedknownhosts",
				},
				ConfigureContainerRegistries: ConfigureContainerRegistries{
					Enabled: true,
					Values:  "configurecontainerregistriesvalues",
				},
				CertManagerDNSChallenge: CertManagerDNSChallenge{
					Enabled:         true,
					AccessKeyID:     "accesskeyid",
					Region:          "region",
					Role:            "role",
					SecretAccessKey: "secretaccesskey",
				},
				CustomCoreDNS: CustomCoreDNS{
					Enabled: true,
					Values:  "customcorednsvalues",
				},
				DisableDenyAllNetPol: true,
				MCProxy: MCProxy{
					Enabled:  true,
					Hostname: "hostname",
					Port:     "1234",
				},
			},
			template:    GetTestTemplate(),
			expectError: false,
		},
		{
			name: "case 1: valid CMC azure",
			cmc: &CMC{
				BaseDomain: "basedomain.io",
				Cluster:    "cluster",
				GitOps: GitOps{
					CMCRepository:         "test-management-clusters",
					CMCBranch:             "cmc-branch",
					MCBBranchSource:       "mcb-branch",
					ConfigBranch:          "config-branch",
					MCAppCollectionBranch: "mc-app-collection-branch",
				},
				ClusterApp: App{
					Name:    "clusterapp-azure",
					Values:  "global:\n  clusterapp: values\nsubscriptionId: subid\ntenantID: tenantid\nclientID: clientid\nclientSecret: clientsecret\nresourceID: uaresourceid\nuaClientID: uaclientid\nuaTenantID: uatenantid\n",
					Version: "clusterappversion",
					Catalog: "clustercatalog",
					AppName: "clusterappname-azure",
				},
				DefaultApps: App{
					Name:    "defaultapp-azure",
					Values:  "defaultappvalues",
					Version: "defaultappversion",
					Catalog: "defaultcatalog",
					AppName: "defaultappname-azure",
				},
				MCAppsPreventDeletion: true,
				PrivateCA:             true,
				ClusterNamespace:      "clusternamespace",
				Provider: Provider{
					Name: key.ProviderAzure,
					CAPZ: CAPZ{
						UAClientID:     "uaclientid",
						UATenantID:     "tenantid",
						UAResourceID:   "uaresourceid",
						ClientID:       "clientid",
						ClientSecret:   "clientsecret",
						TenantID:       "tenantid",
						SubscriptionID: "subid",
					},
				},
				TaylorBotToken: "taylorbottoken",
				SSHdeployKey: DeployKey{
					Identity:   "identity",
					Passphrase: "passphrase",
					KnownHosts: "knownhosts",
				},
				CustomerDeployKey: DeployKey{
					Identity:   "customeridentity",
					Passphrase: "customerpassphrase",
					KnownHosts: "customerknownhosts",
				},
				SharedDeployKey: DeployKey{
					Identity:   "sharedidentity",
					Passphrase: "sharedpassphrase",
					KnownHosts: "sharedknownhosts",
				},
				ConfigureContainerRegistries: ConfigureContainerRegistries{
					Enabled: true,
					Values:  "configurecontainerregistriesvalues",
				},
				CertManagerDNSChallenge: CertManagerDNSChallenge{
					Enabled:         true,
					AccessKeyID:     "accesskeyid",
					Region:          "region",
					Role:            "role",
					SecretAccessKey: "secretaccesskey",
				},
				CustomCoreDNS: CustomCoreDNS{
					Enabled: true,
					Values:  "customcorednsvalues",
				},
				DisableDenyAllNetPol: true,
				MCProxy: MCProxy{
					Enabled:  true,
					Hostname: "hostname",
					Port:     "1234",
				},
			},
			template:    GetTestTemplate(),
			expectError: false,
		},
		{
			name: "case 2: valid aws CMC integrated default apps",
			cmc: &CMC{
				BaseDomain: "basedomain.io",
				Cluster:    "cluster",
				GitOps: GitOps{
					CMCRepository:         "test-management-clusters",
					CMCBranch:             "cmc-branch",
					MCBBranchSource:       "mcb-branch",
					ConfigBranch:          "config-branch",
					MCAppCollectionBranch: "mc-app-collection-branch",
				},
				ClusterApp: App{
					Name:    "clusterapp-aws",
					Values:  "global:\n  clusterapp: values",
					Version: "clusterappversion",
					Catalog: "clustercatalog",
					AppName: "clusterappname-aws",
				},
				ClusterIntegratesDefaultApps: true,
				MCAppsPreventDeletion:        true,
				PrivateCA:                    true,
				ClusterNamespace:             "clusternamespace",
				Provider: Provider{
					Name: key.ProviderAWS,
				},
				TaylorBotToken: "taylorbottoken",
				SSHdeployKey: DeployKey{
					Identity:   "identity",
					Passphrase: "passphrase",
					KnownHosts: "knownhosts",
				},
				CustomerDeployKey: DeployKey{
					Identity:   "customeridentity",
					Passphrase: "customerpassphrase",
					KnownHosts: "customerknownhosts",
				},
				SharedDeployKey: DeployKey{
					Identity:   "sharedidentity",
					Passphrase: "sharedpassphrase",
					KnownHosts: "sharedknownhosts",
				},
				ConfigureContainerRegistries: ConfigureContainerRegistries{
					Enabled: true,
					Values:  "configurecontainerregistriesvalues",
				},
				CertManagerDNSChallenge: CertManagerDNSChallenge{
					Enabled:         true,
					AccessKeyID:     "accesskeyid",
					Region:          "region",
					Role:            "role",
					SecretAccessKey: "secretaccesskey",
				},
				CustomCoreDNS: CustomCoreDNS{
					Enabled: true,
					Values:  "customcorednsvalues",
				},
				DisableDenyAllNetPol: true,
				MCProxy: MCProxy{
					Enabled:  true,
					Hostname: "hostname",
					Port:     "1234",
				},
			},
			template:    GetTestTemplate(),
			expectError: false,
		},
		{
			name: "case 3: valid CMC azure private MC",
			cmc: &CMC{
				BaseDomain: "basedomain.io",
				Cluster:    "cluster",
				GitOps: GitOps{
					CMCRepository:         "test-management-clusters",
					CMCBranch:             "cmc-branch",
					MCBBranchSource:       "mcb-branch",
					ConfigBranch:          "config-branch",
					MCAppCollectionBranch: "mc-app-collection-branch",
				},
				ClusterApp: App{
					Name:    "clusterapp-azure",
					Values:  "global:\n  clusterapp: values\nsubscriptionId: subid\ntenantID: tenantid\nclientID: clientid\nclientSecret: clientsecret\nresourceID: uaresourceid\nuaClientID: uaclientid\nuaTenantID: uatenantid\n",
					Version: "clusterappversion",
					Catalog: "clustercatalog",
					AppName: "clusterappname-azure",
				},
				MCAppsPreventDeletion:        true,
				PrivateCA:                    true,
				PrivateMC:                    true,
				ClusterIntegratesDefaultApps: true,
				ClusterNamespace:             "clusternamespace",
				Provider: Provider{
					Name: key.ProviderAzure,
					CAPZ: CAPZ{
						UAClientID:     "uaclientid",
						UATenantID:     "tenantid",
						UAResourceID:   "uaresourceid",
						ClientID:       "clientid",
						ClientSecret:   "clientsecret",
						TenantID:       "tenantid",
						SubscriptionID: "subid",
					},
				},
				TaylorBotToken: "taylorbottoken",
				SSHdeployKey: DeployKey{
					Identity:   "identity",
					Passphrase: "passphrase",
					KnownHosts: "knownhosts",
				},
				CustomerDeployKey: DeployKey{
					Identity:   "customeridentity",
					Passphrase: "customerpassphrase",
					KnownHosts: "customerknownhosts",
				},
				SharedDeployKey: DeployKey{
					Identity:   "sharedidentity",
					Passphrase: "sharedpassphrase",
					KnownHosts: "sharedknownhosts",
				},
				ConfigureContainerRegistries: ConfigureContainerRegistries{
					Enabled: true,
					Values:  "configurecontainerregistriesvalues",
				},
				CertManagerDNSChallenge: CertManagerDNSChallenge{
					Enabled:         true,
					AccessKeyID:     "accesskeyid",
					Region:          "region",
					Role:            "role",
					SecretAccessKey: "secretaccesskey",
				},
				CustomCoreDNS: CustomCoreDNS{
					Enabled: true,
					Values:  "customcorednsvalues",
				},
				DisableDenyAllNetPol: true,
				MCProxy: MCProxy{
					Enabled:  true,
					Hostname: "hostname",
					Port:     "1234",
				},
			},
			template:    GetTestTemplate(),
			expectError: false,
		},
		{
			name: "case 4: valid CMC azure private MC",
			cmc: &CMC{
				BaseDomain: "basedomain.io",
				Cluster:    "cluster",
				GitOps: GitOps{
					CMCRepository:         "test-management-clusters",
					CMCBranch:             "cmc-branch",
					MCBBranchSource:       "mcb-branch",
					ConfigBranch:          "config-branch",
					MCAppCollectionBranch: "mc-app-collection-branch",
				},
				ClusterApp: App{
					Name:    "clusterapp-azure",
					Values:  "global:\n  clusterapp: values\nsubscriptionId: subid\ntenantID: tenantid\nclientID: clientid\nclientSecret: clientsecret\nresourceID: uaresourceid\nuaClientID: uaclientid\nuaTenantID: uatenantid\n",
					Version: "clusterappversion",
					Catalog: "clustercatalog",
					AppName: "clusterappname-azure",
				},
				DefaultApps: App{
					Name:    "defaultapp-azure",
					Values:  "defaultappvalues",
					Version: "defaultappversion",
					Catalog: "defaultcatalog",
					AppName: "defaultappname-azure",
				},
				MCAppsPreventDeletion: true,
				PrivateCA:             true,
				PrivateMC:             true,
				ClusterNamespace:      "clusternamespace",
				Provider: Provider{
					Name: key.ProviderAzure,
					CAPZ: CAPZ{
						UAClientID:     "uaclientid",
						UATenantID:     "tenantid",
						UAResourceID:   "uaresourceid",
						ClientID:       "clientid",
						ClientSecret:   "clientsecret",
						TenantID:       "tenantid",
						SubscriptionID: "subid",
					},
				},
				TaylorBotToken: "taylorbottoken",
				SSHdeployKey: DeployKey{
					Identity:   "identity",
					Passphrase: "passphrase",
					KnownHosts: "knownhosts",
				},
				CustomerDeployKey: DeployKey{
					Identity:   "customeridentity",
					Passphrase: "customerpassphrase",
					KnownHosts: "customerknownhosts",
				},
				SharedDeployKey: DeployKey{
					Identity:   "sharedidentity",
					Passphrase: "sharedpassphrase",
					KnownHosts: "sharedknownhosts",
				},
				ConfigureContainerRegistries: ConfigureContainerRegistries{
					Enabled: true,
					Values:  "configurecontainerregistriesvalues",
				},
				CertManagerDNSChallenge: CertManagerDNSChallenge{
					Enabled:         true,
					AccessKeyID:     "accesskeyid",
					Region:          "region",
					Role:            "role",
					SecretAccessKey: "secretaccesskey",
				},
				CustomCoreDNS: CustomCoreDNS{
					Enabled: true,
					Values:  "customcorednsvalues",
				},
				DisableDenyAllNetPol: true,
				MCProxy: MCProxy{
					Enabled:  true,
					Hostname: "hostname",
					Port:     "1234",
				},
			},
			template:    GetTestTemplate(),
			expectError: false,
		},
		{
			name: "case 5: custom registry domain",
			cmc: &CMC{
				BaseDomain:     "basedomain.io",
				RegistryDomain: "registrydomain.io",
				Cluster:        "cluster",
				GitOps: GitOps{
					CMCRepository:         "test-management-clusters",
					CMCBranch:             "cmc-branch",
					MCBBranchSource:       "mcb-branch",
					ConfigBranch:          "config-branch",
					MCAppCollectionBranch: "mc-app-collection-branch",
				},
				ClusterApp: App{
					Name:    "clusterapp-aws",
					Values:  "global:\n  clusterapp: values",
					Version: "clusterappversion",
					Catalog: "clustercatalog",
					AppName: "clusterappname-aws",
				},
				ClusterIntegratesDefaultApps: true,
				MCAppsPreventDeletion:        true,
				PrivateCA:                    true,
				ClusterNamespace:             "clusternamespace",
				Provider: Provider{
					Name: key.ProviderAWS,
				},
				TaylorBotToken: "taylorbottoken",
				SSHdeployKey: DeployKey{
					Identity:   "identity",
					Passphrase: "passphrase",
					KnownHosts: "knownhosts",
				},
				CustomerDeployKey: DeployKey{
					Identity:   "customeridentity",
					Passphrase: "customerpassphrase",
					KnownHosts: "customerknownhosts",
				},
				SharedDeployKey: DeployKey{
					Identity:   "sharedidentity",
					Passphrase: "sharedpassphrase",
					KnownHosts: "sharedknownhosts",
				},
				ConfigureContainerRegistries: ConfigureContainerRegistries{
					Enabled: true,
					Values:  "configurecontainerregistriesvalues",
				},
				CertManagerDNSChallenge: CertManagerDNSChallenge{
					Enabled:         true,
					AccessKeyID:     "accesskeyid",
					Region:          "region",
					Role:            "role",
					SecretAccessKey: "secretaccesskey",
				},
				CustomCoreDNS: CustomCoreDNS{
					Enabled: true,
					Values:  "customcorednsvalues",
				},
				DisableDenyAllNetPol: true,
				MCProxy: MCProxy{
					Enabled:  true,
					Hostname: "hostname",
					Port:     "1234",
				},
			},
			template:    GetTestTemplate(),
			expectError: false,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case %d: %s", i, tc.name), func(t *testing.T) {
			if _, ok := os.LookupEnv("CI"); ok { // we skip this test in CI since it needs sops binary to be present right now
				t.Skip()
			}

			agekey, agepubkey, err := GetTestKeys()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			tc.cmc.AgePubKey = agepubkey
			t.Setenv(sops.EnvAgeKey, agekey)
			err = tc.cmc.SetDefaultAppValues()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			m, err := tc.cmc.GetMap(tc.template)
			if err != nil && !tc.expectError {
				t.Fatalf("unexpected error: %v", err)
			} else if err == nil && tc.expectError {
				t.Fatalf("expected error, got nil")
			}

			result, err := GetCMCFromMap(m, tc.cmc.Cluster, "test-management-clusters")
			if err != nil && !tc.expectError {
				t.Fatalf("unexpected error: %v", err)
			} else if err == nil && tc.expectError {
				t.Fatalf("expected error, got nil")
			}

			if !reflect.DeepEqual(result, tc.cmc) {
				t.Fatalf("expected %v, got %v", tc.cmc, result)
			}

			fmt.Print(m)
		})
	}
}

func GetTestKeys() (string, string, error) {
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		return "", "", err
	}
	return identity.String(), identity.Recipient().String(), nil
}

func GetTestTemplate() map[string]string {
	return map[string]string{
		"management-clusters/cluster/deny-all-policies.yaml":                         "deny-all-policies",
		"management-clusters/cluster/kustomization.yaml":                             GetTestKustomization(),
		"management-clusters/cluster/custom-branch-management-clusters-fleet.yaml":   GetTestRepo("CMC_BRANCH"),
		"management-clusters/cluster/custom-branch-config.yaml":                      GetTestRepo("CONFIG_BRANCH"),
		"management-clusters/cluster/custom-branch-collection.yaml":                  GetTestRepo("MC_APP_COLLECTION_BRANCH"),
		"management-clusters/cluster/catalogs/kustomization.yaml":                    GetTestCatalog(),
		"management-clusters/cluster/catalogs/patches/appcatalog-default-patch.yaml": GetTestPatch(),
	}
}

func GetTestKustomization() string {
	return `apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
patches:
  - path: ../../bases/patches/kustomization-post-build.yaml
  - path: https://raw.githubusercontent.com/giantswarm/management-cluster-bases/${MCB_BRANCH_SOURCE}/extras/vaultless/patch-kustomize-controller.yaml
  - path: sops-secret.yaml
  - path: custom-branch-collection.yaml
  - path: custom-branch-config.yaml
  - path: custom-branch-management-clusters-fleet.yaml

patchesStrategicMerge:
  # This cannot be moved under patches, cos there is a bug in kustomize so that it cannot handle multi document patches
  # (https://github.com/kubernetes-sigs/kustomize/issues/5049, fixed in kustomize v5.2.1)
  - https://raw.githubusercontent.com/giantswarm/management-cluster-bases/${MCB_BRANCH_SOURCE}/extras/vaultless/patch-delete-vault-cronjob.yaml
  - https://raw.githubusercontent.com/giantswarm/management-cluster-bases/${MCB_BRANCH_SOURCE}/extras/flux/patch-remove-psp.yaml
replacements:
  # Changes flux Kustomization path to point to correct subpath of this
  # repository.
  - source:
      kind: ConfigMap
      name: management-cluster-metadata
      namespace: flux-giantswarm
      fieldPath: data.NAME
    targets:
      - select:
          kind: Kustomization
          name: crds
          namespace: flux-giantswarm
        fieldPaths:
          - spec.path
        options:
          delimiter: "/"
          index: 2
      - select:
          kind: Kustomization
          name: flux
          namespace: flux-giantswarm
        fieldPaths:
          - spec.path
        options:
          delimiter: "/"
          index: 2
      - select:
          kind: Kustomization
          name: flux-extras
          namespace: flux-giantswarm
        fieldPaths:
          - spec.path
        options:
          delimiter: "/"
          index: 2
      - select:
          kind: Kustomization
          name: catalogs
          namespace: flux-giantswarm
        fieldPaths:
          - spec.path
        options:
          delimiter: "/"
          index: 2
resources:
  - https://github.com/giantswarm/management-cluster-bases//bases/provider/${PROVIDER}/flux-v2?ref=${MCB_BRANCH_SOURCE}
  - configmap-management-cluster-metadata.yaml
  - cluster-app-manifests.yaml
  - default-apps-manifests.yaml
  - deny-all-policies.yaml`
}

func GetTestRepo(branch string) string {
	return fmt.Sprintf(`apiVersion: source.toolkit.fluxcd.io/v1
kind: GitRepository
metadata:
  name: collection
  namespace: flux-giantswarm
spec:
  ref:
    branch: ${%s}`, branch)
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
            baseDomain: ${BASE_DOMAIN}
            managementCluster: ${INSTALLATION}
            provider: ${PROVIDER}
${CATALOG_REGISTRY_VALUES}`
}

func GetTestCatalog() string {
	return `apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
patches:
  - path: patches/appcatalog-default-patch.yaml
resources:
  - https://github.com/giantswarm/management-cluster-bases//bases/catalogs?ref=${MCB_BRANCH_SOURCE}
  - another-resource.yaml`
}
