package pushcmc

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/giantswarm/mcli/pkg/managementcluster/cmc"
)

func TestConfig_ReadSecretFlags(t *testing.T) {
	testCases := []struct {
		name                         string
		secretsFolder                string
		provider                     string
		configureContainerRegistries bool
		configureCertManager         bool
		input                        cmc.CMC

		expectErr    bool
		expectOutput SecretFlags
	}{
		{
			name:          "valid secrets folder",
			secretsFolder: "testdata/valid",
			provider:      "capa",
			expectOutput: SecretFlags{
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
				ClusterValues: "cluster-values\n",
			},
		},
		{
			name:          "nonexistent secrets folder",
			secretsFolder: "testdata/nonexistent",
			expectErr:     true,
		},
		{
			name:                         "valid secrets folder with container registries",
			secretsFolder:                "testdata/validcontainerregistries",
			provider:                     "capa",
			configureContainerRegistries: true,
			expectOutput: SecretFlags{
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
				ClusterValues:                  "cluster-values\n",
				ContainerRegistryConfiguration: "container: registries\nconfiguration:\n- 1\n- 2\nand:\n  so: on\n",
			},
		},
		{
			name:                         "secrets folder with nonexistant container registries",
			secretsFolder:                "testdata/invalidcontainerregistries",
			provider:                     "capa",
			configureContainerRegistries: true,
			expectErr:                    true,
		},
		{
			name:          "valid secrets folder with capz provider",
			secretsFolder: "testdata/validcapz",
			provider:      "capz",
			expectOutput: SecretFlags{
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
				ClusterValues: "cluster-values\n",
				Azure: AzureFlags{
					ClientID:     "test-client-id",
					ClientSecret: "test-client-secret",
					TenantID:     "test-tenant-id",
					UAClientID:   "test-ua-client-id",
					UATenantID:   "test-ua-tenant-id",
					UAResourceID: "test-ua-resource-id",
				},
			},
		},
		{
			name:          "invalid secrets folder with capz provider",
			secretsFolder: "testdata/invalidcapz",
			provider:      "capz",
			expectErr:     true,
		},
		{
			name:          "valid secrets folder with capv provider",
			secretsFolder: "testdata/validcapv",
			provider:      "vsphere",
			expectOutput: SecretFlags{
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
				ClusterValues:      "cluster-values\n",
				VSphereCredentials: "cloud: config\n",
			},
		},
		{
			name:          "invalid secrets folder with capv provider",
			secretsFolder: "testdata/invalidcapv",
			provider:      "vsphere",
			expectErr:     true,
		},
		{
			name:          "valid secrets folder with capvcd provider",
			secretsFolder: "testdata/validcapvcd",
			provider:      "cloud-director",
			expectOutput: SecretFlags{
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
				ClusterValues:             "cluster-values\n",
				CloudDirectorRefreshToken: "test-refresh-token",
			},
		},
		{
			name:          "invalid secrets folder with capvcd provider",
			secretsFolder: "testdata/invalidcapvcd",
			provider:      "cloud-director",
			expectErr:     true,
		},
		{
			name:                 "valid secrets folder with cert manager",
			secretsFolder:        "testdata/validcertmanager",
			provider:             "capa",
			configureCertManager: true,
			expectOutput: SecretFlags{
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
				ClusterValues:                     "cluster-values\n",
				CertManagerRoute53Region:          "region",
				CertManagerRoute53Role:            "arn:aws:iam::1234:role/test-role",
				CertManagerRoute53AccessKeyID:     "test-key",
				CertManagerRoute53SecretAccessKey: "test-secret",
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case %d: %s", i, tc.name), func(t *testing.T) {
			c := &Config{
				Flags: CMCFlags{
					ConfigureContainerRegistries: tc.configureContainerRegistries,
					SecretFolder:                 tc.secretsFolder,
					CertManagerDNSChallenge:      tc.configureCertManager,
				},
				Cluster:  "test",
				Provider: tc.provider,
			}

			err := c.ReadSecretFlags()
			if err != nil && !tc.expectErr {
				t.Fatalf("unexpected error: %v", err)
			} else if err == nil && tc.expectErr {
				t.Fatalf("expected error, got nil")
			}
			if !reflect.DeepEqual(c.Flags.Secrets, tc.expectOutput) {
				t.Fatalf("expected %#v, got %#v", tc.expectOutput, c.Flags.Secrets)
			}
		})
	}
}
