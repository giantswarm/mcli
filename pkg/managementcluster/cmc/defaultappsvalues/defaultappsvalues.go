package defaultappsvalues

import (
	"fmt"

	"github.com/giantswarm/mcli/pkg/key"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Cluster                 string
	Provider                string
	CertManagerDNSChallenge bool
	PrivateCA               bool
}

type DefaultAppsValues struct {
	ClusterName       string     `yaml:"clusterName"`
	Organization      string     `yaml:"organization"`
	ManagementCluster string     `yaml:"managementCluster"`
	UserConfig        UserConfig `yaml:"userConfig"`
	Apps              string     `yaml:"apps"`
}

type UserConfig struct {
	CertManager App `yaml:"certManager"`
	ExternalDNS App `yaml:"externalDNS"`
}

type App struct {
	ConfigMap    ConfigMap     `yaml:"configMap"`
	ExtraConfigs []ExtraConfig `yaml:"extraConfigs"`
}

type ConfigMap struct {
	Values string `yaml:"values"`
}

type ExtraConfig struct {
	Kind string `yaml:"kind"`
	Name string `yaml:"name"`
}

func GetDefaultAppsValuesFile(c Config) (string, error) {
	defaultAppsValues := DefaultAppsValues{
		ClusterName:       c.Cluster,
		Organization:      "giantswarm",
		ManagementCluster: c.Cluster,
	}
	if c.PrivateCA {
		defaultAppsValues.UserConfig.CertManager.ConfigMap.Values = GetCertManagerConfig()
	}
	if c.Provider == key.ProviderAzure {
		defaultAppsValues.UserConfig.ExternalDNS.ConfigMap.Values = GetExternalDNSConfig()
	}
	if c.CertManagerDNSChallenge {
		defaultAppsValues.UserConfig.CertManager.ExtraConfigs = append(defaultAppsValues.UserConfig.CertManager.ExtraConfigs, ExtraConfig{
			Kind: "secret",
			Name: fmt.Sprintf("%s--cert-manager-user-secrets", c.Cluster),
		})
	}
	// marshal the object to yaml
	data, err := yaml.Marshal(defaultAppsValues)
	if err != nil {
		return "", fmt.Errorf("failed to marshal default apps values.\n%w", err)
	}
	return string(data), nil
}

func GetCertManagerConfig() string {
	return `controller:
  defaultIssuer:
    name: private-giantswarm`
}

func GetExternalDNSConfig() string {
	return `hostNetwork: true
flavor: capi
provider: azure
clusterID: {{ .Values.clusterName }}
crd:
  install: false
externalDNS:
  namespaceFilter: \"\"
  sources:
  - ingress`
}
