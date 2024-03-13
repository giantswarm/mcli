package pushcmc

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/giantswarm/mcli/pkg/managementcluster/cmc"
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
				ClusterAppName: "test-app",
				Secrets: SecretFlags{
					ClusterValues: "test-values",
				},
			}},
			expectError: true,
		},
		{
			name: "all flags",
			flags: Config{
				Provider: "capa",
				Cluster:  "test",
				Flags: CMCFlags{
					ClusterAppName:     "test-app",
					ClusterAppCatalog:  "test-catalog",
					ClusterAppVersion:  "1.2.3",
					DefaultAppsName:    "test-default-apps",
					DefaultAppsCatalog: "test-default-catalog",
					DefaultAppsVersion: "3.4.5",
					ClusterNamespace:   "test",
					AgePubKey:          "test-age-pub-key",
					TaylorBotToken:     "test-taylor-bot-token",
					Secrets: SecretFlags{
						ClusterValues: "test-values",
						SSHDeployKey: DeployKey{
							Passphrase: "test-deploy-key-passphrase",
							Identity:   "test-deploy-key-identity",
							KnownHosts: "test-deploy",
						},
						CustomerDeployKey: DeployKey{
							Passphrase: "test-deploy-key-passphrase",
							Identity:   "test-deploy-key-identity",
							KnownHosts: "test-deploy",
						},
						SharedDeployKey: DeployKey{
							Passphrase: "test-deploy-key-passphrase",
							Identity:   "test-deploy-key-identity",
							KnownHosts: "test-deploy",
						},
					},
				},
			},
			expected: &cmc.CMC{
				Provider: cmc.Provider{
					Name: "capa",
				},
				Cluster: "test",
				ClusterApp: cmc.App{
					Name:    "test-app",
					Catalog: "test-catalog",
					Version: "1.2.3",
					Values:  "test-values",
				},
				DefaultApps: cmc.App{
					Name:    "test-default-apps",
					Catalog: "test-default-catalog",
					Version: "3.4.5",
				},
				AgePubKey:        "test-age-pub-key",
				TaylorBotToken:   "test-taylor-bot-token",
				ClusterNamespace: "test",
				SSHdeployKey: cmc.DeployKey{
					Passphrase: "test-deploy-key-passphrase",
					Identity:   "test-deploy-key-identity",
					KnownHosts: "test-deploy",
				},
			},
		},
		{
			name: "all flags with azure",
			flags: Config{
				Provider: "capz",
				Cluster:  "test",
				Flags: CMCFlags{
					ClusterAppName:     "test-app",
					ClusterAppCatalog:  "test-catalog",
					ClusterAppVersion:  "1.2.3",
					DefaultAppsName:    "test-default-apps",
					DefaultAppsCatalog: "test-default-catalog",
					DefaultAppsVersion: "3.4.5",
					ClusterNamespace:   "test",
					AgePubKey:          "test-age-pub-key",
					TaylorBotToken:     "test-taylor-bot-token",
					Secrets: SecretFlags{
						ClusterValues: "test-values",
						SSHDeployKey: DeployKey{
							Passphrase: "test-deploy-key-passphrase",
							Identity:   "test-deploy-key-identity",
							KnownHosts: "test-deploy",
						},
						CustomerDeployKey: DeployKey{
							Passphrase: "test-deploy-key-passphrase",
							Identity:   "test-deploy-key-identity",
							KnownHosts: "test-deploy",
						},
						SharedDeployKey: DeployKey{
							Passphrase: "test-deploy-key-passphrase",
							Identity:   "test-deploy-key-identity",
							KnownHosts: "test-deploy",
						},
						AzureClusterIdentityUA:             "something",
						AzureClusterIdentitySP:             "somethingmore",
						AzureSecretClusterIdentityStaticSP: "somethingmoremore",
					},
				}},
			expected: &cmc.CMC{
				Provider: cmc.Provider{
					Name: "capz",
					CAPZ: cmc.CAPZ{
						IdentityUA:       "something",
						IdentitySP:       "somethingmore",
						IdentityStaticSP: "somethingmoremore",
					},
				},
				Cluster: "test",
				ClusterApp: cmc.App{
					Name:    "test-app",
					Catalog: "test-catalog",
					Version: "1.2.3",
					Values:  "test-values",
				},
				DefaultApps: cmc.App{
					Name:    "test-default-apps",
					Catalog: "test-default-catalog",
					Version: "3.4.5",
				},
				AgePubKey:        "test-age-pub-key",
				TaylorBotToken:   "test-taylor-bot-token",
				ClusterNamespace: "test",
				SSHdeployKey: cmc.DeployKey{
					Passphrase: "test-deploy-key-passphrase",
					Identity:   "test-deploy-key-identity",
					KnownHosts: "test-deploy",
				},
				DisableDenyAllNetPol: true,
			},
		},
		{
			name: "all flags with azure and credentials are missing",
			flags: Config{
				Provider: "capz",
				Cluster:  "test",
				Flags: CMCFlags{
					ClusterAppName:     "test-app",
					ClusterAppCatalog:  "test-catalog",
					ClusterAppVersion:  "1.2.3",
					DefaultAppsName:    "test-default-apps",
					DefaultAppsCatalog: "test-default-catalog",
					DefaultAppsVersion: "3.4.5",
					ClusterNamespace:   "test",
					AgePubKey:          "test-age-pub-key",
					TaylorBotToken:     "test-taylor-bot-token",
					Secrets: SecretFlags{
						ClusterValues: "test-values",
						SSHDeployKey: DeployKey{
							Passphrase: "test-deploy-key-passphrase",
							Identity:   "test-deploy-key-identity",
							KnownHosts: "test-deploy",
						},
						CustomerDeployKey: DeployKey{
							Passphrase: "test-deploy-key-passphrase",
							Identity:   "test-deploy-key-identity",
							KnownHosts: "test-deploy",
						},
						SharedDeployKey: DeployKey{
							Passphrase: "test-deploy-key-passphrase",
							Identity:   "test-deploy-key-identity",
							KnownHosts: "test-deploy",
						},
					},
				}},
			expectError: true,
		},
		{
			name: "CertManager DNS challenge enabled",
			flags: Config{
				Provider: "capa",
				Cluster:  "test",
				Flags: CMCFlags{
					ClusterAppName:          "test-app",
					ClusterAppCatalog:       "test-catalog",
					ClusterAppVersion:       "1.2.3",
					DefaultAppsName:         "test-default-apps",
					DefaultAppsCatalog:      "test-default-catalog",
					DefaultAppsVersion:      "3.4.5",
					ClusterNamespace:        "test",
					AgePubKey:               "test-age-pub-key",
					TaylorBotToken:          "test-taylor-bot-token",
					CertManagerDNSChallenge: true,
					Secrets: SecretFlags{
						ClusterValues: "test-values",
						SSHDeployKey: DeployKey{
							Passphrase: "test-deploy-key-passphrase",
							Identity:   "test-deploy-key-identity",
							KnownHosts: "test-deploy",
						},
						CustomerDeployKey: DeployKey{
							Passphrase: "test-deploy-key-passphrase",
							Identity:   "test-deploy-key-identity",
							KnownHosts: "test-deploy",
						},
						SharedDeployKey: DeployKey{
							Passphrase: "test-deploy-key-passphrase",
							Identity:   "test-deploy-key-identity",
							KnownHosts: "test-deploy",
						},
						CertManagerRoute53Region:          "us-west-2",
						CertManagerRoute53Role:            "cert-manager-role",
						CertManagerRoute53AccessKeyID:     "access-key-id",
						CertManagerRoute53SecretAccessKey: "secret-access-key",
					},
				},
			},
			expected: &cmc.CMC{
				Provider: cmc.Provider{
					Name: "capa",
				},
				Cluster: "test",
				ClusterApp: cmc.App{
					Name:    "test-app",
					Catalog: "test-catalog",
					Version: "1.2.3",
					Values:  "test-values",
				},
				DefaultApps: cmc.App{
					Name:    "test-default-apps",
					Catalog: "test-default-catalog",
					Version: "3.4.5",
				},
				AgePubKey:        "test-age-pub-key",
				TaylorBotToken:   "test-taylor-bot-token",
				ClusterNamespace: "test",
				SSHdeployKey: cmc.DeployKey{
					Passphrase: "test-deploy-key-passphrase",
					Identity:   "test-deploy-key-identity",
					KnownHosts: "test-deploy",
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
				Provider: "vsphere",
				Cluster:  "test",
				Flags: CMCFlags{
					ClusterAppName:     "test-app",
					ClusterAppCatalog:  "test-catalog",
					ClusterAppVersion:  "1.2.3",
					DefaultAppsName:    "test-default-apps",
					DefaultAppsCatalog: "test-default-catalog",
					DefaultAppsVersion: "3.4.5",
					ClusterNamespace:   "test",
					AgePubKey:          "test-age-pub-key",
					TaylorBotToken:     "test-taylor-bot-token",
					Secrets: SecretFlags{
						ClusterValues: "test-values",
						SSHDeployKey: DeployKey{
							Passphrase: "test-deploy-key-passphrase",
							Identity:   "test-deploy-key-identity",
							KnownHosts: "test-deploy",
						},
						CustomerDeployKey: DeployKey{
							Passphrase: "test-deploy-key-passphrase",
							Identity:   "test-deploy-key-identity",
							KnownHosts: "test-deploy",
						},
						SharedDeployKey: DeployKey{
							Passphrase: "test-deploy-key-passphrase",
							Identity:   "test-deploy-key-identity",
							KnownHosts: "test-deploy",
						},
						VSphereCredentials: "test-vsphere-credentials",
					},
				},
			},
			expected: &cmc.CMC{
				Provider: cmc.Provider{
					Name: "vsphere",
					CAPV: cmc.CAPV{
						CloudConfig: "test-vsphere-credentials",
					},
				},
				Cluster: "test",
				ClusterApp: cmc.App{
					Name:    "test-app",
					Catalog: "test-catalog",
					Version: "1.2.3",
					Values:  "test-values",
				},
				DefaultApps: cmc.App{
					Name:    "test-default-apps",
					Catalog: "test-default-catalog",
					Version: "3.4.5",
				},
				AgePubKey:        "test-age-pub-key",
				TaylorBotToken:   "test-taylor-bot-token",
				ClusterNamespace: "test",
				SSHdeployKey: cmc.DeployKey{
					Passphrase: "test-deploy-key-passphrase",
					Identity:   "test-deploy-key-identity",
					KnownHosts: "test-deploy",
				},
			},
		},
		{
			name: "Provider vcd",
			flags: Config{
				Provider: "cloud-director",
				Cluster:  "test",
				Flags: CMCFlags{
					ClusterAppName:     "test-app",
					ClusterAppCatalog:  "test-catalog",
					ClusterAppVersion:  "1.2.3",
					DefaultAppsName:    "test-default-apps",
					DefaultAppsCatalog: "test-default-catalog",
					DefaultAppsVersion: "3.4.5",
					ClusterNamespace:   "test",
					AgePubKey:          "test-age-pub-key",
					TaylorBotToken:     "test-taylor-bot-token",
					Secrets: SecretFlags{
						ClusterValues: "test-values",
						SSHDeployKey: DeployKey{
							Passphrase: "test-deploy-key-passphrase",
							Identity:   "test-deploy-key-identity",
							KnownHosts: "test-deploy",
						},
						CustomerDeployKey: DeployKey{
							Passphrase: "test-deploy-key-passphrase",
							Identity:   "test-deploy-key-identity",
							KnownHosts: "test-deploy",
						},
						SharedDeployKey: DeployKey{
							Passphrase: "test-deploy-key-passphrase",
							Identity:   "test-deploy-key-identity",
							KnownHosts: "test-deploy",
						},
						CloudDirectorCredentials: "test-vcd-credentials",
					},
				},
			},
			expected: &cmc.CMC{
				Provider: cmc.Provider{
					Name: "cloud-director",
					CAPVCD: cmc.CAPVCD{
						CloudConfig: "test-vcd-credentials",
					},
				},
				Cluster: "test",
				ClusterApp: cmc.App{
					Name:    "test-app",
					Catalog: "test-catalog",
					Version: "1.2.3",
					Values:  "test-values",
				},
				DefaultApps: cmc.App{
					Name:    "test-default-apps",
					Catalog: "test-default-catalog",
					Version: "3.4.5",
				},
				AgePubKey:        "test-age-pub-key",
				TaylorBotToken:   "test-taylor-bot-token",
				ClusterNamespace: "test",
				SSHdeployKey: cmc.DeployKey{
					Passphrase: "test-deploy-key-passphrase",
					Identity:   "test-deploy-key-identity",
					KnownHosts: "test-deploy",
				},
			},
		},
		{
			name: "Configure container registries enabled",
			flags: Config{
				Provider: "capa",
				Cluster:  "test",
				Flags: CMCFlags{
					ClusterAppName:               "test-app",
					ClusterAppCatalog:            "test-catalog",
					ClusterAppVersion:            "1.2.3",
					DefaultAppsName:              "test-default-apps",
					DefaultAppsCatalog:           "test-default-catalog",
					DefaultAppsVersion:           "3.4.5",
					ClusterNamespace:             "test",
					AgePubKey:                    "test-age-pub-key",
					TaylorBotToken:               "test-taylor-bot-token",
					ConfigureContainerRegistries: true,
					Secrets: SecretFlags{
						ClusterValues: "test-values",
						SSHDeployKey: DeployKey{
							Passphrase: "test-deploy-key-passphrase",
							Identity:   "test-deploy-key-identity",
							KnownHosts: "test-deploy",
						},
						CustomerDeployKey: DeployKey{
							Passphrase: "test-deploy-key-passphrase",
							Identity:   "test-deploy-key-identity",
							KnownHosts: "test-deploy",
						},
						SharedDeployKey: DeployKey{
							Passphrase: "test-deploy-key-passphrase",
							Identity:   "test-deploy-key-identity",
							KnownHosts: "test-deploy",
						},
						ContainerRegistryConfiguration: "test-container-registry-configuration",
					},
				},
			},
			expected: &cmc.CMC{
				Provider: cmc.Provider{
					Name: "capa",
				},
				Cluster: "test",
				ClusterApp: cmc.App{
					Name:    "test-app",
					Catalog: "test-catalog",
					Version: "1.2.3",
					Values:  "test-values",
				},
				DefaultApps: cmc.App{
					Name:    "test-default-apps",
					Catalog: "test-default-catalog",
					Version: "3.4.5",
				},
				AgePubKey:        "test-age-pub-key",
				TaylorBotToken:   "test-taylor-bot-token",
				ClusterNamespace: "test",
				SSHdeployKey: cmc.DeployKey{
					Passphrase: "test-deploy-key-passphrase",
					Identity:   "test-deploy-key-identity",
					KnownHosts: "test-deploy",
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
				Provider: "capa",
				Cluster:  "test",
				Flags: CMCFlags{
					ClusterAppName:     "test-app",
					ClusterAppCatalog:  "test-catalog",
					ClusterAppVersion:  "1.2.3",
					DefaultAppsName:    "test-default-apps",
					DefaultAppsCatalog: "test-default-catalog",
					DefaultAppsVersion: "3.4.5",
					ClusterNamespace:   "test",
					AgePubKey:          "test-age-pub-key",
					TaylorBotToken:     "test-taylor-bot-token",
					MCProxyEnabled:     true,
					MCHTTPSProxy:       "test-mc-https-proxy",
					Secrets: SecretFlags{
						ClusterValues: "test-values",
						SSHDeployKey: DeployKey{
							Passphrase: "test-deploy-key-passphrase",
							Identity:   "test-deploy-key-identity",
							KnownHosts: "test-deploy",
						},
						CustomerDeployKey: DeployKey{
							Passphrase: "test-deploy-key-passphrase",
							Identity:   "test-deploy-key-identity",
							KnownHosts: "test-deploy",
						},
						SharedDeployKey: DeployKey{
							Passphrase: "test-deploy-key-passphrase",
							Identity:   "test-deploy-key-identity",
							KnownHosts: "test-deploy",
						},
					},
				},
			},
			expected: &cmc.CMC{
				Provider: cmc.Provider{
					Name: "capa",
				},
				Cluster: "test",
				ClusterApp: cmc.App{
					Name:    "test-app",
					Catalog: "test-catalog",
					Version: "1.2.3",
					Values:  "test-values",
				},
				DefaultApps: cmc.App{
					Name:    "test-default-apps",
					Catalog: "test-default-catalog",
					Version: "3.4.5",
				},
				AgePubKey:        "test-age-pub-key",
				TaylorBotToken:   "test-taylor-bot-token",
				ClusterNamespace: "test",
				SSHdeployKey: cmc.DeployKey{
					Passphrase: "test-deploy-key-passphrase",
					Identity:   "test-deploy-key-identity",
					KnownHosts: "test-deploy",
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
				Provider: "vsphere",
				Cluster:  "test",
				Flags: CMCFlags{
					ClusterAppName:     "test-app",
					ClusterAppCatalog:  "test-catalog",
					ClusterAppVersion:  "1.2.3",
					DefaultAppsName:    "test-default-apps",
					DefaultAppsCatalog: "test-default-catalog",
					DefaultAppsVersion: "3.4.5",
					ClusterNamespace:   "test",
					AgePubKey:          "test-age-pub-key",
					TaylorBotToken:     "test-taylor-bot-token",
					Secrets: SecretFlags{
						ClusterValues: "test-values",
						SSHDeployKey: DeployKey{
							Passphrase: "test-deploy-key-passphrase",
							Identity:   "test-deploy-key-identity",
							KnownHosts: "test-deploy",
						},
						CustomerDeployKey: DeployKey{
							Passphrase: "test-deploy-key-passphrase",
							Identity:   "test-deploy-key-identity",
							KnownHosts: "test-deploy",
						},
						SharedDeployKey: DeployKey{
							Passphrase: "test-deploy-key-passphrase",
							Identity:   "test-deploy-key-identity",
							KnownHosts: "test-deploy",
						},
					},
				},
			},
			expectError: true,
		},
		{
			name: "Provider vcd missing configuration",
			flags: Config{
				Provider: "cloud-director",
				Cluster:  "test",
				Flags: CMCFlags{
					ClusterAppName:     "test-app",
					ClusterAppCatalog:  "test-catalog",
					ClusterAppVersion:  "1.2.3",
					DefaultAppsName:    "test-default-apps",
					DefaultAppsCatalog: "test-default-catalog",
					DefaultAppsVersion: "3.4.5",
					ClusterNamespace:   "test",
					AgePubKey:          "test-age-pub-key",
					TaylorBotToken:     "test-taylor-bot-token",
					Secrets: SecretFlags{
						ClusterValues: "test-values",
						SSHDeployKey: DeployKey{
							Passphrase: "test-deploy-key-passphrase",
							Identity:   "test-deploy-key-identity",
							KnownHosts: "test-deploy",
						},
						CustomerDeployKey: DeployKey{
							Passphrase: "test-deploy-key-passphrase",
							Identity:   "test-deploy-key-identity",
							KnownHosts: "test-deploy",
						},
						SharedDeployKey: DeployKey{
							Passphrase: "test-deploy-key-passphrase",
							Identity:   "test-deploy-key-identity",
							KnownHosts: "test-deploy",
						},
					},
				},
			},
			expectError: true,
		},
		{
			name: "Configure container registries enabled missing configuration",
			flags: Config{
				Provider: "capa",
				Cluster:  "test",
				Flags: CMCFlags{
					ClusterAppName:               "test-app",
					ClusterAppCatalog:            "test-catalog",
					ClusterAppVersion:            "1.2.3",
					DefaultAppsName:              "test-default-apps",
					DefaultAppsCatalog:           "test-default-catalog",
					DefaultAppsVersion:           "3.4.5",
					ClusterNamespace:             "test",
					AgePubKey:                    "test-age-pub-key",
					TaylorBotToken:               "test-taylor-bot-token",
					ConfigureContainerRegistries: true,
					Secrets: SecretFlags{
						ClusterValues: "test-values",
						SSHDeployKey: DeployKey{
							Passphrase: "test-deploy-key-passphrase",
							Identity:   "test-deploy-key-identity",
							KnownHosts: "test-deploy",
						},
						CustomerDeployKey: DeployKey{
							Passphrase: "test-deploy-key-passphrase",
							Identity:   "test-deploy-key-identity",
							KnownHosts: "test-deploy",
						},
						SharedDeployKey: DeployKey{
							Passphrase: "test-deploy-key-passphrase",
							Identity:   "test-deploy-key-identity",
							KnownHosts: "test-deploy",
						},
					},
				},
			},
			expectError: true,
		},
		{
			name: "MC proxy enabled missing configuration",
			flags: Config{
				Provider: "capa",
				Cluster:  "test",
				Flags: CMCFlags{
					ClusterAppName:     "test-app",
					ClusterAppCatalog:  "test-catalog",
					ClusterAppVersion:  "1.2.3",
					DefaultAppsName:    "test-default-apps",
					DefaultAppsCatalog: "test-default-catalog",
					DefaultAppsVersion: "3.4.5",
					ClusterNamespace:   "test",
					AgePubKey:          "test-age-pub-key",
					TaylorBotToken:     "test-taylor-bot-token",
					MCProxyEnabled:     true,
					Secrets: SecretFlags{
						ClusterValues: "test-values",
						SSHDeployKey: DeployKey{
							Passphrase: "test-deploy-key-passphrase",
							Identity:   "test-deploy-key-identity",
							KnownHosts: "test-deploy",
						},
						CustomerDeployKey: DeployKey{
							Passphrase: "test-deploy-key-passphrase",
							Identity:   "test-deploy-key-identity",
							KnownHosts: "test-deploy",
						},
						SharedDeployKey: DeployKey{
							Passphrase: "test-deploy-key-passphrase",
							Identity:   "test-deploy-key-identity",
							KnownHosts: "test-deploy",
						},
					},
				},
			},
			expectError: true,
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
				t.Fatalf("expected %#v, got %#v", tc.expected, installation)
			}
		})
	}
}
