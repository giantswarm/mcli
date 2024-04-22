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
				Cluster: "cluster",
				ClusterApp: App{
					Name:    "clusterapp-aws",
					Values:  "clusterappvalues",
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
			template:    map[string]string{},
			expectError: false,
		},
		{
			name: "case 1: valid CMC azure",
			cmc: &CMC{
				Cluster: "cluster",
				ClusterApp: App{
					Name:    "clusterapp-azure",
					Values:  "clusterappvalues",
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
						UAClientID:   "uaclientid",
						UATenantID:   "tenantid",
						UAResourceID: "uaresourceid",
						ClientID:     "clientid",
						ClientSecret: "clientsecret",
						TenantID:     "tenantid",
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
			template:    map[string]string{},
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

			m, err := tc.cmc.GetMap(tc.template)
			if err != nil && !tc.expectError {
				t.Fatalf("unexpected error: %v", err)
			} else if err == nil && tc.expectError {
				t.Fatalf("expected error, got nil")
			}

			result, err := GetCMCFromMap(m, tc.cmc.Cluster)
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
