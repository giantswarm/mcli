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

type IntegratedDefaultAppsValues struct {
	Global Global `yaml:"global,omitempty"`
}

type Global struct {
	Apps UserConfig `yaml:"apps,omitempty"`
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

func IntegrateDefaultAppsValuesInClusterValues(clusterValues string, c Config) (string, error) {
	log.Debug().Msg("Integrating default apps values in cluster values")

	if !c.PrivateCA &&
		!(key.IsProviderAzure(c.Provider) && !c.PrivateMC) &&
		!c.CertManagerDNSChallenge {
		return clusterValues, nil
	}

	apps := IntegratedDefaultAppsValues{}

	if c.PrivateCA {
		apps.Global.Apps.CertManager.ExtraConfigs = []ExtraConfig{
			{
				Kind: "configmap",
				Name: "cert-manager-user-values",
			},
		}
	}
	if key.IsProviderAzure(c.Provider) {
		if !c.PrivateMC {
			apps.Global.Apps.ExternalDNS.ExtraConfigs = []ExtraConfig{
				{
					Kind: "configmap",
					Name: "external-dns-user-values",
				},
			}
		}
	}
	if c.CertManagerDNSChallenge {
		apps.Global.Apps.CertManager.ExtraConfigs = append(apps.Global.Apps.CertManager.ExtraConfigs, ExtraConfig{
			Kind: "secret",
			Name: key.GetCertManagerSecretName(c.Cluster),
		})
	}

	data, err := key.GetData(apps)
	if err != nil {
		return "", fmt.Errorf("failed to marshal integrated default apps values.\n%w", err)
	}
	clusterValues, err = key.MergeValues(clusterValues, string(data))
	if err != nil {
		return "", fmt.Errorf("failed to merge default apps values in cluster values.\n%w", err)
	}

	return clusterValues, nil
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
