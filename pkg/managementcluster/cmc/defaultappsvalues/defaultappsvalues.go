package defaultappsvalues

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/giantswarm/mcli/pkg/key"
)

const (
	IdentityClientIDKey = "identityClientID"
)

type Config struct {
	Cluster                 string
	Provider                string
	CertManagerDNSChallenge bool
	PrivateCA               bool
	PrivateMC               bool
	SubscriptionID          string
	IdentityClientID        string
}

type DefaultAppsValues struct {
	ClusterName       string     `yaml:"clusterName"`
	Organization      string     `yaml:"organization"`
	ManagementCluster string     `yaml:"managementCluster"`
	UserConfig        UserConfig `yaml:"userConfig,omitempty"`
	Apps              string     `yaml:"apps,omitempty"`
	SubscriptionID    string     `yaml:"subscriptionID,omitempty"`
	IdentityClientID  string     `yaml:"identityClientID,omitempty"`
}

type UserConfig struct {
	CertManager App `yaml:"certManager,omitempty"`
	ExternalDNS App `yaml:"externalDNS,omitempty"`
}

type App struct {
	ConfigMap    ConfigMap     `yaml:"configMap,omitempty"`
	ExtraConfigs []ExtraConfig `yaml:"extraConfigs,omitempty"`
}

type ConfigMap struct {
	Values string `yaml:"values,omitempty"`
}

type ExtraConfig struct {
	Kind      string `yaml:"kind,omitempty"`
	Name      string `yaml:"name,omitempty"`
	Namespace string `yaml:"namespace,omitempty"`
}

func GetDefaultAppsValuesFile(c Config) (string, error) {
	log.Debug().Msg("Creating default apps values")

	defaultAppsValues := DefaultAppsValues{
		ClusterName:       c.Cluster,
		Organization:      "giantswarm",
		ManagementCluster: c.Cluster,
	}
	if c.PrivateCA {
		defaultAppsValues.UserConfig.CertManager.ConfigMap.Values = GetCertManagerConfig()
	}
	if key.IsProviderAzure(c.Provider) {
		if !c.PrivateMC {
			defaultAppsValues.UserConfig.ExternalDNS.ConfigMap.Values = GetExternalDNSConfig()
		} else {
			defaultAppsValues.SubscriptionID = c.SubscriptionID
			defaultAppsValues.IdentityClientID = c.IdentityClientID
		}
	}
	if c.CertManagerDNSChallenge {
		defaultAppsValues.UserConfig.CertManager.ExtraConfigs = append(defaultAppsValues.UserConfig.CertManager.ExtraConfigs, ExtraConfig{
			Kind:      "secret",
			Namespace: "org-giantswarm",
			Name:      key.GetCertManagerSecretName(c.Cluster),
		})
	}
	// marshal the object to yaml
	data, err := key.GetData(defaultAppsValues)
	if err != nil {
		return "", fmt.Errorf("failed to marshal default apps values.\n%w", err)
	}
	return string(data), nil
}

func IsPrivateMC(file string) bool {
	v, err := key.GetValue(IdentityClientIDKey, file)
	if err == nil && v != "" {
		return true
	}
	return false
}

func GetCertManagerConfig() string {
	return `controller:
  defaultIssuer:
    name: private-giantswarm`
}

func GetExternalDNSConfig() string {
	return `flavor: capi
provider: azure
clusterID: {{ .Values.clusterName }}
crd:
  install: false
externalDNS:
  namespaceFilter: \"\"
  sources:
  - ingress`
}
