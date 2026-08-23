package pushcmc

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/giantswarm/mcli/pkg/managementcluster/cmc"

	"github.com/giantswarm/mcli/pkg/key"
)

const (
	testAppName             = "test-app"
	testValues              = "test-values"
	testName                = "test"
	testBaseDomain          = "test.base.domain.io"
	testBranchSource        = "test-branch-source"
	testConfigBranch        = "test-config-branch"
	testAppCollectionBranch = "test-app-collection-branch"
	testCatalog             = "test-catalog"
	testClusterAppVersion   = "1.2.3"
	testDefaultAppsName     = "test-default-apps"
	testDefaultAppsCatalog  = "test-default-catalog"
	testDefaultAppsVersion  = "3.4.5"
	testAgePubKey           = "test-age-pub-key"
	testTaylorBotToken      = "test-taylor-bot-token"
	testDeployKeyPassphrase = "test-deploy-key-passphrase"
	testDeployKeyIdentity   = "test-deploy-key-identity"
	testDeployKeyKnownHosts = "test-deploy"
	testClusterValues       = "clusterName: test\norganization: giantswarm\nmanagementCluster: test\n"
	testClientID            = "test-client-id"
	testClientSecret        = "test-client"
	testTenantID            = "test-tenant"
	testSubscriptionID      = "test-subscription-id"
	testUAClientID          = "test-ua-client-id"
	testUATenantID          = "test-ua-tenant"
	testUAResourceID        = "test-ua-resource-id"
)

func TestGetNewCMCFromFlags(t *testing.T) {
	var testCases = []struct {
		name  string
		flags Config

		expected    *cmc.CMC
		expectError bool
	}{
		{
			name:        "no flags",
			flags:       Config{Flags: CMCFlags{}},
			expectError: true,
		},
		{
			name: "some flags",
			flags: Config{Flags: CMCFlags{
				ClusterAppName: testAppName,
				Secrets: SecretFlags{
					ClusterValues: testValues,
				},
			}},
			expectError: true,
		},
		{
			name: "all flags",
			flags: Config{
				Provider:   key.ProviderAWS,
				Cluster:    testName,
				BaseDomain: testBaseDomain,
				Flags: CMCFlags{
					MCBBranchSource:       testBranchSource,
					ConfigBranch:          testConfigBranch,
					MCAppCollectionBranch: testAppCollectionBranch,
					ClusterAppName:        testAppName,
					ClusterAppCatalog:     testCatalog,
					ClusterAppVersion:     testClusterAppVersion,
					DefaultAppsName:       testDefaultAppsName,
					DefaultAppsCatalog:    testDefaultAppsCatalog,
					DefaultAppsVersion:    testDefaultAppsVersion,
					ClusterNamespace:      testName,
					AgePubKey:             testAgePubKey,
					TaylorBotToken:        testTaylorBotToken,
					Secrets: SecretFlags{
						ClusterValues: testValues,
						SSHDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						CustomerDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						SharedDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
					},
				},
			},
			expected: &cmc.CMC{
				GitOps: cmc.GitOps{
					MCBBranchSource:       testBranchSource,
					ConfigBranch:          testConfigBranch,
					MCAppCollectionBranch: testAppCollectionBranch,
				},
				BaseDomain: testBaseDomain,
				Provider: cmc.Provider{
					Name: key.ProviderAWS,
				},
				Cluster: testName,
				ClusterApp: cmc.App{
					Name:    testAppName,
					AppName: testName,
					Catalog: testCatalog,
					Version: testClusterAppVersion,
					Values:  testValues,
				},
				DefaultApps: cmc.App{
					Name:    testDefaultAppsName,
					AppName: testDefaultAppsName,
					Catalog: testDefaultAppsCatalog,
					Version: testDefaultAppsVersion,
					Values:  testClusterValues,
				},
				AgePubKey: testAgePubKey,

				TaylorBotToken:   testTaylorBotToken,
				ClusterNamespace: testName,
				SSHdeployKey: cmc.DeployKey{
					Passphrase: testDeployKeyPassphrase,
					Identity:   testDeployKeyIdentity,
					KnownHosts: testDeployKeyKnownHosts,
				},
				CustomerDeployKey: cmc.DeployKey{
					Passphrase: testDeployKeyPassphrase,
					Identity:   testDeployKeyIdentity,
					KnownHosts: testDeployKeyKnownHosts,
				},
				SharedDeployKey: cmc.DeployKey{
					Passphrase: testDeployKeyPassphrase,
					Identity:   testDeployKeyIdentity,
					KnownHosts: testDeployKeyKnownHosts,
				},
			},
		},
		{
			name: "all flags with azure",
			flags: Config{
				Provider:   key.ProviderAzure,
				Cluster:    testName,
				BaseDomain: testBaseDomain,
				Flags: CMCFlags{
					MCBBranchSource:       testBranchSource,
					ConfigBranch:          testConfigBranch,
					MCAppCollectionBranch: testAppCollectionBranch,
					ClusterAppName:        testAppName,
					ClusterAppCatalog:     testCatalog,
					ClusterAppVersion:     testClusterAppVersion,
					DefaultAppsName:       testDefaultAppsName,
					DefaultAppsCatalog:    testDefaultAppsCatalog,
					DefaultAppsVersion:    testDefaultAppsVersion,
					ClusterNamespace:      testName,
					AgePubKey:             testAgePubKey,

					TaylorBotToken: testTaylorBotToken,
					Secrets: SecretFlags{
						ClusterValues: testValues,
						SSHDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						CustomerDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						SharedDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						Azure: AzureFlags{
							ClientID:       testClientID,
							ClientSecret:   testClientSecret,
							TenantID:       testTenantID,
							SubscriptionID: testSubscriptionID,
							UAClientID:     testUAClientID,
							UATenantID:     testUATenantID,
							UAResourceID:   testUAResourceID,
						},
					},
				},
			},
			expected: &cmc.CMC{
				GitOps: cmc.GitOps{
					MCBBranchSource:       testBranchSource,
					ConfigBranch:          testConfigBranch,
					MCAppCollectionBranch: testAppCollectionBranch,
				},
				BaseDomain: testBaseDomain,
				Provider: cmc.Provider{
					Name: key.ProviderAzure,
					CAPZ: cmc.CAPZ{
						UAClientID:     testUAClientID,
						UATenantID:     testUATenantID,
						UAResourceID:   testUAResourceID,
						ClientID:       testClientID,
						ClientSecret:   testClientSecret,
						TenantID:       testTenantID,
						SubscriptionID: testSubscriptionID,
					},
				},
				Cluster: testName,
				ClusterApp: cmc.App{
					Name:    testAppName,
					AppName: testName,
					Catalog: testCatalog,
					Version: testClusterAppVersion,
					Values:  testValues,
				},
				DefaultApps: cmc.App{
					Name:    testDefaultAppsName,
					AppName: testDefaultAppsName,
					Catalog: testDefaultAppsCatalog,
					Version: testDefaultAppsVersion,
					Values:  "clusterName: test\norganization: giantswarm\nmanagementCluster: test\nuserConfig:\n  externalDNS:\n    configMap:\n      values: |-\n        flavor: capi\n        provider: azure\n        clusterID: {{ .Values.clusterName }}\n        crd:\n          install: false\n        externalDNS:\n          namespaceFilter: \\\"\\\"\n          sources:\n          - ingress\n",
				},
				AgePubKey: testAgePubKey,

				TaylorBotToken:   testTaylorBotToken,
				ClusterNamespace: testName,
				SSHdeployKey: cmc.DeployKey{
					Passphrase: testDeployKeyPassphrase,
					Identity:   testDeployKeyIdentity,
					KnownHosts: testDeployKeyKnownHosts,
				},
				CustomerDeployKey: cmc.DeployKey{
					Passphrase: testDeployKeyPassphrase,
					Identity:   testDeployKeyIdentity,
					KnownHosts: testDeployKeyKnownHosts,
				},
				SharedDeployKey: cmc.DeployKey{
					Passphrase: testDeployKeyPassphrase,
					Identity:   testDeployKeyIdentity,
					KnownHosts: testDeployKeyKnownHosts,
				},
				DisableDenyAllNetPol: true,
			},
		},
		{
			name: "all flags with azure and credentials are missing",
			flags: Config{
				BaseDomain: testBaseDomain,
				Provider:   key.ProviderAzure,
				Cluster:    testName,
				Flags: CMCFlags{
					MCBBranchSource:       testBranchSource,
					ConfigBranch:          testConfigBranch,
					MCAppCollectionBranch: testAppCollectionBranch,
					ClusterAppName:        testAppName,
					ClusterAppCatalog:     testCatalog,
					ClusterAppVersion:     testClusterAppVersion,
					DefaultAppsName:       testDefaultAppsName,
					DefaultAppsCatalog:    testDefaultAppsCatalog,
					DefaultAppsVersion:    testDefaultAppsVersion,
					ClusterNamespace:      testName,
					AgePubKey:             testAgePubKey,

					TaylorBotToken: testTaylorBotToken,
					Secrets: SecretFlags{
						ClusterValues: testValues,
						SSHDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						CustomerDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						SharedDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
					},
				},
			},
			expectError: true,
		},
		{
			name: "CertManager DNS challenge enabled",
			flags: Config{
				BaseDomain: testBaseDomain,
				Provider:   key.ProviderAWS,
				Cluster:    testName,
				Flags: CMCFlags{
					MCBBranchSource:       testBranchSource,
					ConfigBranch:          testConfigBranch,
					MCAppCollectionBranch: testAppCollectionBranch,
					ClusterAppName:        testAppName,
					ClusterAppCatalog:     testCatalog,
					ClusterAppVersion:     testClusterAppVersion,
					DefaultAppsName:       testDefaultAppsName,
					DefaultAppsCatalog:    testDefaultAppsCatalog,
					DefaultAppsVersion:    testDefaultAppsVersion,
					ClusterNamespace:      testName,
					AgePubKey:             testAgePubKey,

					TaylorBotToken:          testTaylorBotToken,
					CertManagerDNSChallenge: true,
					Secrets: SecretFlags{
						ClusterValues: testValues,
						SSHDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						CustomerDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						SharedDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						CertManagerRoute53Region:          "us-west-2",
						CertManagerRoute53Role:            "cert-manager-role",
						CertManagerRoute53AccessKeyID:     "access-key-id",
						CertManagerRoute53SecretAccessKey: "secret-access-key",
					},
				},
			},
			expected: &cmc.CMC{
				GitOps: cmc.GitOps{
					MCBBranchSource:       testBranchSource,
					ConfigBranch:          testConfigBranch,
					MCAppCollectionBranch: testAppCollectionBranch,
				},
				BaseDomain: testBaseDomain,
				Provider: cmc.Provider{
					Name: key.ProviderAWS,
				},
				Cluster: testName,
				ClusterApp: cmc.App{
					Name:    testAppName,
					AppName: testName,
					Catalog: testCatalog,
					Version: testClusterAppVersion,
					Values:  testValues,
				},
				DefaultApps: cmc.App{
					Name:    testDefaultAppsName,
					AppName: testDefaultAppsName,
					Catalog: testDefaultAppsCatalog,
					Version: testDefaultAppsVersion,
					Values:  "clusterName: test\norganization: giantswarm\nmanagementCluster: test\nuserConfig:\n  certManager:\n    extraConfigs:\n      - kind: secret\n        name: test-cert-manager-user-secrets\n        namespace: org-giantswarm\n",
				},
				AgePubKey: testAgePubKey,

				TaylorBotToken:   testTaylorBotToken,
				ClusterNamespace: testName,
				SSHdeployKey: cmc.DeployKey{
					Passphrase: testDeployKeyPassphrase,
					Identity:   testDeployKeyIdentity,
					KnownHosts: testDeployKeyKnownHosts,
				},
				CustomerDeployKey: cmc.DeployKey{
					Passphrase: testDeployKeyPassphrase,
					Identity:   testDeployKeyIdentity,
					KnownHosts: testDeployKeyKnownHosts,
				},
				SharedDeployKey: cmc.DeployKey{
					Passphrase: testDeployKeyPassphrase,
					Identity:   testDeployKeyIdentity,
					KnownHosts: testDeployKeyKnownHosts,
				},
				CertManagerDNSChallenge: cmc.CertManagerDNSChallenge{
					Enabled:         true,
					Region:          "us-west-2",
					Role:            "cert-manager-role",
					AccessKeyID:     "access-key-id",
					SecretAccessKey: "secret-access-key",
				},
			},
		},
		{
			name: "Provider vsphere",
			flags: Config{
				BaseDomain: testBaseDomain,
				Provider:   key.ProviderVsphere,
				Cluster:    testName,
				Flags: CMCFlags{
					MCBBranchSource:       testBranchSource,
					ConfigBranch:          testConfigBranch,
					MCAppCollectionBranch: testAppCollectionBranch,
					ClusterAppName:        testAppName,
					ClusterAppCatalog:     testCatalog,
					ClusterAppVersion:     testClusterAppVersion,
					DefaultAppsName:       testDefaultAppsName,
					DefaultAppsCatalog:    testDefaultAppsCatalog,
					DefaultAppsVersion:    testDefaultAppsVersion,
					ClusterNamespace:      testName,
					AgePubKey:             testAgePubKey,

					TaylorBotToken: testTaylorBotToken,
					Secrets: SecretFlags{
						ClusterValues: testValues,
						SSHDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						CustomerDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						SharedDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						VSphereCredentials: "test-vsphere-credentials",
					},
				},
			},
			expected: &cmc.CMC{
				GitOps: cmc.GitOps{
					MCBBranchSource:       testBranchSource,
					ConfigBranch:          testConfigBranch,
					MCAppCollectionBranch: testAppCollectionBranch,
				},
				BaseDomain: testBaseDomain,
				Provider: cmc.Provider{
					Name: key.ProviderVsphere,
					CAPV: cmc.CAPV{
						CloudConfig: "test-vsphere-credentials",
					},
				},
				Cluster: testName,
				ClusterApp: cmc.App{
					Name:    testAppName,
					AppName: testName,
					Catalog: testCatalog,
					Version: testClusterAppVersion,
					Values:  testValues,
				},
				DefaultApps: cmc.App{
					Name:    testDefaultAppsName,
					AppName: testDefaultAppsName,
					Catalog: testDefaultAppsCatalog,
					Version: testDefaultAppsVersion,
					Values:  testClusterValues,
				},
				AgePubKey: testAgePubKey,

				TaylorBotToken:   testTaylorBotToken,
				ClusterNamespace: testName,
				SSHdeployKey: cmc.DeployKey{
					Passphrase: testDeployKeyPassphrase,
					Identity:   testDeployKeyIdentity,
					KnownHosts: testDeployKeyKnownHosts,
				},
				CustomerDeployKey: cmc.DeployKey{
					Passphrase: testDeployKeyPassphrase,
					Identity:   testDeployKeyIdentity,
					KnownHosts: testDeployKeyKnownHosts,
				},
				SharedDeployKey: cmc.DeployKey{
					Passphrase: testDeployKeyPassphrase,
					Identity:   testDeployKeyIdentity,
					KnownHosts: testDeployKeyKnownHosts,
				},
			},
		},
		{
			name: "Provider vcd",
			flags: Config{
				BaseDomain: testBaseDomain,
				Provider:   key.ProviderVCD,
				Cluster:    testName,
				Flags: CMCFlags{
					MCBBranchSource:       testBranchSource,
					ConfigBranch:          testConfigBranch,
					MCAppCollectionBranch: testAppCollectionBranch,
					ClusterAppName:        testAppName,
					ClusterAppCatalog:     testCatalog,
					ClusterAppVersion:     testClusterAppVersion,
					DefaultAppsName:       testDefaultAppsName,
					DefaultAppsCatalog:    testDefaultAppsCatalog,
					DefaultAppsVersion:    testDefaultAppsVersion,
					ClusterNamespace:      testName,
					AgePubKey:             testAgePubKey,

					TaylorBotToken: testTaylorBotToken,
					Secrets: SecretFlags{
						ClusterValues: testValues,
						SSHDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						CustomerDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						SharedDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						CloudDirectorRefreshToken: "test-vcd-credentials",
					},
				},
			},
			expected: &cmc.CMC{
				GitOps: cmc.GitOps{
					MCBBranchSource:       testBranchSource,
					ConfigBranch:          testConfigBranch,
					MCAppCollectionBranch: testAppCollectionBranch,
				},
				BaseDomain: testBaseDomain,
				Provider: cmc.Provider{
					Name: key.ProviderVCD,
					CAPVCD: cmc.CAPVCD{
						RefreshToken: "test-vcd-credentials",
					},
				},
				Cluster: testName,
				ClusterApp: cmc.App{
					Name:    testAppName,
					AppName: testName,
					Catalog: testCatalog,
					Version: testClusterAppVersion,
					Values:  testValues,
				},
				DefaultApps: cmc.App{
					Name:    testDefaultAppsName,
					AppName: testDefaultAppsName,
					Catalog: testDefaultAppsCatalog,
					Version: testDefaultAppsVersion,
					Values:  testClusterValues,
				},
				AgePubKey: testAgePubKey,

				TaylorBotToken:   testTaylorBotToken,
				ClusterNamespace: testName,
				SSHdeployKey: cmc.DeployKey{
					Passphrase: testDeployKeyPassphrase,
					Identity:   testDeployKeyIdentity,
					KnownHosts: testDeployKeyKnownHosts,
				},
				CustomerDeployKey: cmc.DeployKey{
					Passphrase: testDeployKeyPassphrase,
					Identity:   testDeployKeyIdentity,
					KnownHosts: testDeployKeyKnownHosts,
				},
				SharedDeployKey: cmc.DeployKey{
					Passphrase: testDeployKeyPassphrase,
					Identity:   testDeployKeyIdentity,
					KnownHosts: testDeployKeyKnownHosts,
				},
			},
		},
		{
			name: "Configure container registries enabled",
			flags: Config{
				BaseDomain: testBaseDomain,
				Provider:   key.ProviderAWS,
				Cluster:    testName,
				Flags: CMCFlags{
					MCBBranchSource:       testBranchSource,
					ConfigBranch:          testConfigBranch,
					MCAppCollectionBranch: testAppCollectionBranch,
					ClusterAppName:        testAppName,
					ClusterAppCatalog:     testCatalog,
					ClusterAppVersion:     testClusterAppVersion,
					DefaultAppsName:       testDefaultAppsName,
					DefaultAppsCatalog:    testDefaultAppsCatalog,
					DefaultAppsVersion:    testDefaultAppsVersion,
					ClusterNamespace:      testName,
					AgePubKey:             testAgePubKey,

					TaylorBotToken:               testTaylorBotToken,
					ConfigureContainerRegistries: true,
					Secrets: SecretFlags{
						ClusterValues: testValues,
						SSHDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						CustomerDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						SharedDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						ContainerRegistryConfiguration: "test-container-registry-configuration",
					},
				},
			},
			expected: &cmc.CMC{
				GitOps: cmc.GitOps{
					MCBBranchSource:       testBranchSource,
					ConfigBranch:          testConfigBranch,
					MCAppCollectionBranch: testAppCollectionBranch,
				},
				BaseDomain: testBaseDomain,
				Provider: cmc.Provider{
					Name: key.ProviderAWS,
				},
				Cluster: testName,
				ClusterApp: cmc.App{
					Name:    testAppName,
					AppName: testName,
					Catalog: testCatalog,
					Version: testClusterAppVersion,
					Values:  testValues,
				},
				DefaultApps: cmc.App{
					Name:    testDefaultAppsName,
					AppName: testDefaultAppsName,
					Catalog: testDefaultAppsCatalog,
					Version: testDefaultAppsVersion,
					Values:  testClusterValues,
				},
				AgePubKey: testAgePubKey,

				TaylorBotToken:   testTaylorBotToken,
				ClusterNamespace: testName,
				SSHdeployKey: cmc.DeployKey{
					Passphrase: testDeployKeyPassphrase,
					Identity:   testDeployKeyIdentity,
					KnownHosts: testDeployKeyKnownHosts,
				},
				CustomerDeployKey: cmc.DeployKey{
					Passphrase: testDeployKeyPassphrase,
					Identity:   testDeployKeyIdentity,
					KnownHosts: testDeployKeyKnownHosts,
				},
				SharedDeployKey: cmc.DeployKey{
					Passphrase: testDeployKeyPassphrase,
					Identity:   testDeployKeyIdentity,
					KnownHosts: testDeployKeyKnownHosts,
				},
				ConfigureContainerRegistries: cmc.ConfigureContainerRegistries{
					Enabled: true,
					Values:  "test-container-registry-configuration",
				},
			},
		},
		{
			name: "MC proxy enabled",
			flags: Config{
				BaseDomain: testBaseDomain,
				Provider:   key.ProviderAWS,
				Cluster:    testName,
				Flags: CMCFlags{
					MCBBranchSource:       testBranchSource,
					ConfigBranch:          testConfigBranch,
					MCAppCollectionBranch: testAppCollectionBranch,
					ClusterAppName:        testAppName,
					ClusterAppCatalog:     testCatalog,
					ClusterAppVersion:     testClusterAppVersion,
					DefaultAppsName:       testDefaultAppsName,
					DefaultAppsCatalog:    testDefaultAppsCatalog,
					DefaultAppsVersion:    testDefaultAppsVersion,
					ClusterNamespace:      testName,
					AgePubKey:             testAgePubKey,

					TaylorBotToken: testTaylorBotToken,
					MCProxyEnabled: true,
					MCHTTPSProxy:   "http://test-mc-https-proxy:443",
					Secrets: SecretFlags{
						ClusterValues: testValues,
						SSHDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						CustomerDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						SharedDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
					},
				},
			},
			expected: &cmc.CMC{
				GitOps: cmc.GitOps{
					MCBBranchSource:       testBranchSource,
					ConfigBranch:          testConfigBranch,
					MCAppCollectionBranch: testAppCollectionBranch,
				},
				BaseDomain: testBaseDomain,
				Provider: cmc.Provider{
					Name: key.ProviderAWS,
				},
				Cluster: testName,
				ClusterApp: cmc.App{
					Name:    testAppName,
					AppName: testName,
					Catalog: testCatalog,
					Version: testClusterAppVersion,
					Values:  testValues,
				},
				DefaultApps: cmc.App{
					Name:    testDefaultAppsName,
					AppName: testDefaultAppsName,
					Catalog: testDefaultAppsCatalog,
					Version: testDefaultAppsVersion,
					Values:  testClusterValues,
				},
				AgePubKey: testAgePubKey,

				TaylorBotToken:   testTaylorBotToken,
				ClusterNamespace: testName,
				SSHdeployKey: cmc.DeployKey{
					Passphrase: testDeployKeyPassphrase,
					Identity:   testDeployKeyIdentity,
					KnownHosts: testDeployKeyKnownHosts,
				},
				CustomerDeployKey: cmc.DeployKey{
					Passphrase: testDeployKeyPassphrase,
					Identity:   testDeployKeyIdentity,
					KnownHosts: testDeployKeyKnownHosts,
				},
				SharedDeployKey: cmc.DeployKey{
					Passphrase: testDeployKeyPassphrase,
					Identity:   testDeployKeyIdentity,
					KnownHosts: testDeployKeyKnownHosts,
				},
				MCProxy: cmc.MCProxy{
					Enabled:  true,
					Hostname: "test-mc-https-proxy",
					Port:     "443",
				},
			},
		},
		{
			name: "Provider vsphere missing configuration",
			flags: Config{
				BaseDomain: testBaseDomain,
				Provider:   key.ProviderVsphere,
				Cluster:    testName,
				Flags: CMCFlags{
					MCBBranchSource:       testBranchSource,
					ConfigBranch:          testConfigBranch,
					MCAppCollectionBranch: testAppCollectionBranch,
					ClusterAppName:        testAppName,
					ClusterAppCatalog:     testCatalog,
					ClusterAppVersion:     testClusterAppVersion,
					DefaultAppsName:       testDefaultAppsName,
					DefaultAppsCatalog:    testDefaultAppsCatalog,
					DefaultAppsVersion:    testDefaultAppsVersion,
					ClusterNamespace:      testName,
					AgePubKey:             testAgePubKey,

					TaylorBotToken: testTaylorBotToken,
					Secrets: SecretFlags{
						ClusterValues: testValues,
						SSHDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						CustomerDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						SharedDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
					},
				},
			},
			expectError: true,
		},
		{
			name: "Provider vcd missing configuration",
			flags: Config{
				BaseDomain: testBaseDomain,
				Provider:   key.ProviderVCD,
				Cluster:    testName,
				Flags: CMCFlags{
					MCBBranchSource:       testBranchSource,
					ConfigBranch:          testConfigBranch,
					MCAppCollectionBranch: testAppCollectionBranch,
					ClusterAppName:        testAppName,
					ClusterAppCatalog:     testCatalog,
					ClusterAppVersion:     testClusterAppVersion,
					DefaultAppsName:       testDefaultAppsName,
					DefaultAppsCatalog:    testDefaultAppsCatalog,
					DefaultAppsVersion:    testDefaultAppsVersion,
					ClusterNamespace:      testName,
					AgePubKey:             testAgePubKey,

					TaylorBotToken: testTaylorBotToken,
					Secrets: SecretFlags{
						ClusterValues: testValues,
						SSHDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						CustomerDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						SharedDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
					},
				},
			},
			expectError: true,
		},
		{
			name: "Configure container registries enabled missing configuration",
			flags: Config{
				BaseDomain: testBaseDomain,
				Provider:   key.ProviderAWS,
				Cluster:    testName,
				Flags: CMCFlags{
					MCBBranchSource:       testBranchSource,
					ConfigBranch:          testConfigBranch,
					MCAppCollectionBranch: testAppCollectionBranch,
					ClusterAppName:        testAppName,
					ClusterAppCatalog:     testCatalog,
					ClusterAppVersion:     testClusterAppVersion,
					DefaultAppsName:       testDefaultAppsName,
					DefaultAppsCatalog:    testDefaultAppsCatalog,
					DefaultAppsVersion:    testDefaultAppsVersion,
					ClusterNamespace:      testName,
					AgePubKey:             testAgePubKey,

					TaylorBotToken:               testTaylorBotToken,
					ConfigureContainerRegistries: true,
					Secrets: SecretFlags{
						ClusterValues: testValues,
						SSHDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						CustomerDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						SharedDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
					},
				},
			},
			expectError: true,
		},
		{
			name: "MC proxy enabled missing configuration",
			flags: Config{
				BaseDomain: testBaseDomain,
				Provider:   key.ProviderAWS,
				Cluster:    testName,
				Flags: CMCFlags{
					MCBBranchSource:       testBranchSource,
					ConfigBranch:          testConfigBranch,
					MCAppCollectionBranch: testAppCollectionBranch,
					ClusterAppName:        testAppName,
					ClusterAppCatalog:     testCatalog,
					ClusterAppVersion:     testClusterAppVersion,
					DefaultAppsName:       testDefaultAppsName,
					DefaultAppsCatalog:    testDefaultAppsCatalog,
					DefaultAppsVersion:    testDefaultAppsVersion,
					ClusterNamespace:      testName,
					AgePubKey:             testAgePubKey,

					TaylorBotToken: testTaylorBotToken,
					MCProxyEnabled: true,
					Secrets: SecretFlags{
						ClusterValues: testValues,
						SSHDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						CustomerDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						SharedDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
					},
				},
			},
			expectError: true,
		},
		{
			name: "azure private MC",
			flags: Config{
				BaseDomain: testBaseDomain,
				Provider:   key.ProviderAzure,
				Cluster:    testName,
				Flags: CMCFlags{
					MCBBranchSource:       testBranchSource,
					ConfigBranch:          testConfigBranch,
					MCAppCollectionBranch: testAppCollectionBranch,
					ClusterAppName:        testAppName,
					ClusterAppCatalog:     testCatalog,
					ClusterAppVersion:     testClusterAppVersion,
					DefaultAppsName:       testDefaultAppsName,
					DefaultAppsCatalog:    testDefaultAppsCatalog,
					DefaultAppsVersion:    testDefaultAppsVersion,
					ClusterNamespace:      testName,
					PrivateMC:             true,
					AgePubKey:             testAgePubKey,

					TaylorBotToken: testTaylorBotToken,
					Secrets: SecretFlags{
						ClusterValues: testValues,
						SSHDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						CustomerDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						SharedDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						Azure: AzureFlags{
							ClientID:       testClientID,
							ClientSecret:   testClientSecret,
							TenantID:       testTenantID,
							SubscriptionID: testSubscriptionID,
							UAClientID:     testUAClientID,
							UATenantID:     testUATenantID,
							UAResourceID:   testUAResourceID,
						},
					},
				},
			},
			expected: &cmc.CMC{
				GitOps: cmc.GitOps{
					MCBBranchSource:       testBranchSource,
					ConfigBranch:          testConfigBranch,
					MCAppCollectionBranch: testAppCollectionBranch,
				},
				BaseDomain: testBaseDomain,
				Provider: cmc.Provider{
					Name: key.ProviderAzure,
					CAPZ: cmc.CAPZ{
						UAClientID:     testUAClientID,
						UATenantID:     testUATenantID,
						UAResourceID:   testUAResourceID,
						ClientID:       testClientID,
						ClientSecret:   testClientSecret,
						TenantID:       testTenantID,
						SubscriptionID: testSubscriptionID,
					},
				},
				Cluster: testName,
				ClusterApp: cmc.App{
					Name:    testAppName,
					AppName: testName,
					Catalog: testCatalog,
					Version: testClusterAppVersion,
					Values:  testValues,
				},
				DefaultApps: cmc.App{
					Name:    testDefaultAppsName,
					AppName: testDefaultAppsName,
					Catalog: testDefaultAppsCatalog,
					Version: testDefaultAppsVersion,
					Values:  "clusterName: test\norganization: giantswarm\nmanagementCluster: test\nsubscriptionID: test-subscription-id\nidentityClientID: test-ua-client-id\n",
				},
				AgePubKey: testAgePubKey,

				TaylorBotToken:   testTaylorBotToken,
				ClusterNamespace: testName,
				SSHdeployKey: cmc.DeployKey{
					Passphrase: testDeployKeyPassphrase,
					Identity:   testDeployKeyIdentity,
					KnownHosts: testDeployKeyKnownHosts,
				},
				CustomerDeployKey: cmc.DeployKey{
					Passphrase: testDeployKeyPassphrase,
					Identity:   testDeployKeyIdentity,
					KnownHosts: testDeployKeyKnownHosts,
				},
				SharedDeployKey: cmc.DeployKey{
					Passphrase: testDeployKeyPassphrase,
					Identity:   testDeployKeyIdentity,
					KnownHosts: testDeployKeyKnownHosts,
				},
				DisableDenyAllNetPol: true,
				PrivateMC:            true,
			},
		},
		{
			name: "integrated default apps",
			flags: Config{
				BaseDomain: testBaseDomain,
				Provider:   key.ProviderAWS,
				Cluster:    testName,
				Flags: CMCFlags{
					MCBBranchSource:              testBranchSource,
					ConfigBranch:                 testConfigBranch,
					MCAppCollectionBranch:        testAppCollectionBranch,
					ClusterAppName:               testAppName,
					ClusterAppCatalog:            testCatalog,
					ClusterAppVersion:            testClusterAppVersion,
					ClusterIntegratesDefaultApps: true,
					ClusterNamespace:             testName,
					AgePubKey:                    testAgePubKey,

					TaylorBotToken: testTaylorBotToken,
					Secrets: SecretFlags{
						ClusterValues: testValues,
						SSHDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						CustomerDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						SharedDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
					},
				},
			},
			expected: &cmc.CMC{
				GitOps: cmc.GitOps{
					MCBBranchSource:       testBranchSource,
					ConfigBranch:          testConfigBranch,
					MCAppCollectionBranch: testAppCollectionBranch,
				},
				BaseDomain: testBaseDomain,
				Provider: cmc.Provider{
					Name: key.ProviderAWS,
				},
				Cluster: testName,
				ClusterApp: cmc.App{
					Name:    testAppName,
					AppName: testName,
					Catalog: testCatalog,
					Version: testClusterAppVersion,
					Values:  testValues,
				},
				ClusterIntegratesDefaultApps: true,
				AgePubKey:                    testAgePubKey,

				TaylorBotToken:   testTaylorBotToken,
				ClusterNamespace: testName,
				SSHdeployKey: cmc.DeployKey{
					Passphrase: testDeployKeyPassphrase,
					Identity:   testDeployKeyIdentity,
					KnownHosts: testDeployKeyKnownHosts,
				},
				CustomerDeployKey: cmc.DeployKey{
					Passphrase: testDeployKeyPassphrase,
					Identity:   testDeployKeyIdentity,
					KnownHosts: testDeployKeyKnownHosts,
				},
				SharedDeployKey: cmc.DeployKey{
					Passphrase: testDeployKeyPassphrase,
					Identity:   testDeployKeyIdentity,
					KnownHosts: testDeployKeyKnownHosts,
				},
			},
		},
		{
			name: "private mc with azure",
			flags: Config{
				BaseDomain: testBaseDomain,
				Provider:   key.ProviderAzure,
				Cluster:    testName,
				Flags: CMCFlags{
					MCBBranchSource:              testBranchSource,
					ConfigBranch:                 testConfigBranch,
					MCAppCollectionBranch:        testAppCollectionBranch,
					ClusterAppName:               testAppName,
					ClusterAppCatalog:            testCatalog,
					ClusterAppVersion:            testClusterAppVersion,
					ClusterIntegratesDefaultApps: true,
					PrivateMC:                    true,
					ClusterNamespace:             testName,
					AgePubKey:                    testAgePubKey,

					TaylorBotToken: testTaylorBotToken,
					Secrets: SecretFlags{
						ClusterValues: testValues,
						SSHDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						CustomerDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						SharedDeployKey: DeployKey{
							Passphrase: testDeployKeyPassphrase,
							Identity:   testDeployKeyIdentity,
							KnownHosts: testDeployKeyKnownHosts,
						},
						Azure: AzureFlags{
							ClientID:       testClientID,
							ClientSecret:   testClientSecret,
							TenantID:       testTenantID,
							SubscriptionID: testSubscriptionID,
							UAClientID:     testUAClientID,
							UATenantID:     testUATenantID,
							UAResourceID:   testUAResourceID,
						},
					},
				},
			},
			expected: &cmc.CMC{
				GitOps: cmc.GitOps{
					MCBBranchSource:       testBranchSource,
					ConfigBranch:          testConfigBranch,
					MCAppCollectionBranch: testAppCollectionBranch,
				},
				BaseDomain: testBaseDomain,
				Provider: cmc.Provider{
					Name: key.ProviderAzure,
					CAPZ: cmc.CAPZ{
						UAClientID:     testUAClientID,
						UATenantID:     testUATenantID,
						UAResourceID:   testUAResourceID,
						ClientID:       testClientID,
						ClientSecret:   testClientSecret,
						TenantID:       testTenantID,
						SubscriptionID: testSubscriptionID,
					},
				},
				Cluster: testName,
				ClusterApp: cmc.App{
					Name:    testAppName,
					AppName: testName,
					Catalog: testCatalog,
					Version: testClusterAppVersion,
					Values:  testValues,
				},
				ClusterIntegratesDefaultApps: true,
				PrivateMC:                    true,
				AgePubKey:                    testAgePubKey,

				TaylorBotToken:   testTaylorBotToken,
				ClusterNamespace: testName,
				SSHdeployKey: cmc.DeployKey{
					Passphrase: testDeployKeyPassphrase,
					Identity:   testDeployKeyIdentity,
					KnownHosts: testDeployKeyKnownHosts,
				},
				CustomerDeployKey: cmc.DeployKey{
					Passphrase: testDeployKeyPassphrase,
					Identity:   testDeployKeyIdentity,
					KnownHosts: testDeployKeyKnownHosts,
				},
				SharedDeployKey: cmc.DeployKey{
					Passphrase: testDeployKeyPassphrase,
					Identity:   testDeployKeyIdentity,
					KnownHosts: testDeployKeyKnownHosts,
				},
				DisableDenyAllNetPol: true,
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case %d: %s", i, tc.name), func(t *testing.T) {
			installation, err := getNewCMCFromFlags(tc.flags)
			if err != nil && !tc.expectError {
				t.Fatalf("unexpected error: %v", err)
			} else if err == nil && tc.expectError {
				t.Fatalf("expected error, got nil")
			}
			if !reflect.DeepEqual(installation, tc.expected) {
				t.Fatalf("expected %v, got %v", tc.expected, installation)
			}
		})
	}
}
