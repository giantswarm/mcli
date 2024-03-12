package cmc

import (
	"bytes"
	"fmt"
	"os"
	"reflect"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"

	"github.com/giantswarm/mcli/pkg/key"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc/defaultappsvalues"
)

type CMC struct {
	AgePubKey                    string                       `yaml:"agePubKey"`
	AgeKey                       string                       `yaml:"ageKey"`
	Cluster                      string                       `yaml:"cluster"`
	ClusterApp                   App                          `yaml:"clusterApp"`
	DefaultApps                  App                          `yaml:"defaultApps"`
	MCAppsPreventDeletion        bool                         `yaml:"mcAppsPreventDeletion"`
	PrivateCA                    bool                         `yaml:"privateCA"`
	ClusterNamespace             string                       `yaml:"clusterNamespace"`
	Provider                     Provider                     `yaml:"provider"`
	TaylorBotToken               string                       `yaml:"taylorBotToken"`
	SSHdeployKey                 DeployKey                    `yaml:"sshDeployKey"`
	CustomerDeployKey            DeployKey                    `yaml:"customerDeployKey"`
	SharedDeployKey              DeployKey                    `yaml:"sharedDeployKey"`
	CertManagerDNSChallenge      CertManagerDNSChallenge      `yaml:"certManagerDNSChallenge"`
	ConfigureContainerRegistries ConfigureContainerRegistries `yaml:"configureContainerRegistries"`
	CustomCoreDNS                bool                         `yaml:"customCoreDNS"`
	DisableDenyAllNetPol         bool                         `yaml:"disableDenyAllNetPol"`
	MCProxy                      MCProxy                      `yaml:"mcProxy"`
}

type App struct {
	Name    string `yaml:"name"`
	Catalog string `yaml:"catalog"`
	Version string `yaml:"version"`
	Values  string `yaml:"values"`
	AppName string `yaml:"appName,omitempty"`
}

type Provider struct {
	Name   string `yaml:"name"`
	CAPV   CAPV   `yaml:"capv,omitempty"`
	CAPZ   CAPZ   `yaml:"capz,omitempty"`
	CAPVCD CAPVCD `yaml:"capvcd,omitempty"`
}

type CAPV struct {
	CloudConfig string `yaml:"cloudConfig"`
}

type CAPZ struct {
	IdentityUA       string `yaml:"identityUA"`
	IdentitySP       string `yaml:"identitySP"`
	IdentityStaticSP string `yaml:"identityStaticSP"`
}

type CAPVCD struct {
	CloudConfig string `yaml:"cloudConfig"`
}

type CertManagerDNSChallenge struct {
	Enabled         bool   `yaml:"enabled"`
	Region          string `yaml:"region,omitempty"`
	Role            string `yaml:"role,omitempty"`
	AccessKeyID     string `yaml:"accessKeyID,omitempty"`
	SecretAccessKey string `yaml:"secretAccessKey,omitempty"`
}

type ConfigureContainerRegistries struct {
	Enabled bool   `yaml:"enabled"`
	Values  string `yaml:"values,omitempty"`
}

type DeployKey struct {
	Passphrase string `yaml:"key"`
	Identity   string `yaml:"identity"`
	KnownHosts string `yaml:"knownHosts"`
}

type MCProxy struct {
	Enabled    bool   `yaml:"enabled"`
	HTTPSProxy string `yaml:"httpsProxy,omitempty"`
}

func GetCMC(data []byte) (*CMC, error) {
	log.Debug().Msg("getting CMC object")
	cmc := CMC{}
	if err := yaml.Unmarshal(data, &cmc); err != nil {
		return nil, fmt.Errorf("failed to unmarshal CMC.\n%w", err)
	}
	return &cmc, nil
}

func GetData(c *CMC) ([]byte, error) {
	log.Debug().Msg("getting data from CMC object")
	w := new(bytes.Buffer)
	encoder := yaml.NewEncoder(w)
	encoder.SetIndent(2)
	err := encoder.Encode(c)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal CMC.\n%w", err)
	}
	return w.Bytes(), nil
}

func (c *CMC) Print() error {
	data, err := GetData(c)
	if err != nil {
		return err
	}
	log.Debug().Msg("printing CMC object")
	// TODO hide secrets
	fmt.Print(string(data))
	return nil
}

func (c *CMC) Override(override *CMC) *CMC {
	cmc := *c
	if override.AgePubKey != "" {
		cmc.AgePubKey = override.AgePubKey
	}
	if override.AgeKey != "" {
		cmc.AgeKey = override.AgeKey
	}
	if override.Cluster != "" {
		cmc.Cluster = override.Cluster
	}
	if override.ClusterApp.Name != "" {
		cmc.ClusterApp.Name = override.ClusterApp.Name
	}
	if override.ClusterApp.AppName != "" {
		cmc.ClusterApp.AppName = override.ClusterApp.AppName
	}
	if override.ClusterApp.Catalog != "" {
		cmc.ClusterApp.Catalog = override.ClusterApp.Catalog
	}
	if override.ClusterApp.Version != "" {
		cmc.ClusterApp.Version = override.ClusterApp.Version
	}
	if override.ClusterApp.Values != "" {
		cmc.ClusterApp.Values = override.ClusterApp.Values
	}
	if override.DefaultApps.Name != "" {
		cmc.DefaultApps.Name = override.DefaultApps.Name
	}
	if override.DefaultApps.AppName != "" {
		cmc.DefaultApps.AppName = override.DefaultApps.AppName
	}
	if override.DefaultApps.Catalog != "" {
		cmc.DefaultApps.Catalog = override.DefaultApps.Catalog
	}
	if override.DefaultApps.Version != "" {
		cmc.DefaultApps.Version = override.DefaultApps.Version
	}
	if override.DefaultApps.Values != "" {
		cmc.DefaultApps.Values = override.DefaultApps.Values
	}
	if override.MCAppsPreventDeletion {
		cmc.MCAppsPreventDeletion = override.MCAppsPreventDeletion
	}
	if override.PrivateCA {
		cmc.PrivateCA = override.PrivateCA
	}
	if override.ClusterNamespace != "" {
		cmc.ClusterNamespace = override.ClusterNamespace
	}
	if override.Provider.Name != "" {
		cmc.Provider.Name = override.Provider.Name
		if override.Provider.Name == key.ProviderVsphere {
			if override.Provider.CAPV.CloudConfig != "" {
				cmc.Provider.CAPV.CloudConfig = override.Provider.CAPV.CloudConfig
			}
		} else if override.Provider.Name == key.ProviderAzure {
			if override.Provider.CAPZ.IdentityUA != "" {
				cmc.Provider.CAPZ.IdentityUA = override.Provider.CAPZ.IdentityUA
			}
			if override.Provider.CAPZ.IdentitySP != "" {
				cmc.Provider.CAPZ.IdentitySP = override.Provider.CAPZ.IdentitySP
			}
			if override.Provider.CAPZ.IdentityStaticSP != "" {
				cmc.Provider.CAPZ.IdentityStaticSP = override.Provider.CAPZ.IdentityStaticSP
			}
		} else if override.Provider.Name == key.ProviderVCD {
			if override.Provider.CAPVCD.CloudConfig != "" {
				cmc.Provider.CAPVCD.CloudConfig = override.Provider.CAPVCD.CloudConfig
			}
		}
	}
	if override.TaylorBotToken != "" {
		cmc.TaylorBotToken = override.TaylorBotToken
	}
	if override.SSHdeployKey.Passphrase != "" {
		cmc.SSHdeployKey.Passphrase = override.SSHdeployKey.Passphrase
	}
	if override.SSHdeployKey.Identity != "" {
		cmc.SSHdeployKey.Identity = override.SSHdeployKey.Identity
	}
	if override.SSHdeployKey.KnownHosts != "" {
		cmc.SSHdeployKey.KnownHosts = override.SSHdeployKey.KnownHosts
	}
	if override.CustomerDeployKey.Passphrase != "" {
		cmc.CustomerDeployKey.Passphrase = override.CustomerDeployKey.Passphrase
	}
	if override.CustomerDeployKey.Identity != "" {
		cmc.CustomerDeployKey.Identity = override.CustomerDeployKey.Identity
	}
	if override.CustomerDeployKey.KnownHosts != "" {
		cmc.CustomerDeployKey.KnownHosts = override.CustomerDeployKey.KnownHosts
	}
	if override.SharedDeployKey.Passphrase != "" {
		cmc.SharedDeployKey.Passphrase = override.SharedDeployKey.Passphrase
	}
	if override.SharedDeployKey.Identity != "" {
		cmc.SharedDeployKey.Identity = override.SharedDeployKey.Identity
	}
	if override.SharedDeployKey.KnownHosts != "" {
		cmc.SharedDeployKey.KnownHosts = override.SharedDeployKey.KnownHosts
	}
	if override.CertManagerDNSChallenge.Enabled {
		cmc.CertManagerDNSChallenge.Enabled = override.CertManagerDNSChallenge.Enabled
		if override.CertManagerDNSChallenge.Region != "" {
			cmc.CertManagerDNSChallenge.Region = override.CertManagerDNSChallenge.Region
		}
		if override.CertManagerDNSChallenge.Role != "" {
			cmc.CertManagerDNSChallenge.Role = override.CertManagerDNSChallenge.Role
		}
		if override.CertManagerDNSChallenge.AccessKeyID != "" {
			cmc.CertManagerDNSChallenge.AccessKeyID = override.CertManagerDNSChallenge.AccessKeyID
		}
		if override.CertManagerDNSChallenge.SecretAccessKey != "" {
			cmc.CertManagerDNSChallenge.SecretAccessKey = override.CertManagerDNSChallenge.SecretAccessKey
		}
	}
	if override.ConfigureContainerRegistries.Enabled {
		cmc.ConfigureContainerRegistries.Enabled = override.ConfigureContainerRegistries.Enabled
		if override.ConfigureContainerRegistries.Values != "" {
			cmc.ConfigureContainerRegistries.Values = override.ConfigureContainerRegistries.Values
		}
	}
	if override.CustomCoreDNS {
		cmc.CustomCoreDNS = override.CustomCoreDNS
	}
	if override.DisableDenyAllNetPol {
		cmc.DisableDenyAllNetPol = override.DisableDenyAllNetPol
	}
	if override.MCProxy.Enabled {
		cmc.MCProxy.Enabled = override.MCProxy.Enabled
		if override.MCProxy.HTTPSProxy != "" {
			cmc.MCProxy.HTTPSProxy = override.MCProxy.HTTPSProxy
		}
	}
	return &cmc
}

func (c *CMC) Validate() error {
	if c.AgePubKey == "" {
		return fmt.Errorf("age public key is empty")
	}
	if c.AgeKey == "" {
		return fmt.Errorf("age key is empty")
	}
	if c.Cluster == "" {
		return fmt.Errorf("cluster is empty")
	}
	if c.ClusterApp.Name == "" {
		return fmt.Errorf("cluster app name is empty")
	}
	if c.ClusterApp.AppName == "" {
		return fmt.Errorf("cluster app app name is empty")
	}
	if c.ClusterApp.Catalog == "" {
		return fmt.Errorf("cluster app catalog is empty")
	}
	if c.ClusterApp.Version == "" {
		return fmt.Errorf("cluster app version is empty")
	}
	if c.ClusterApp.Values == "" {
		return fmt.Errorf("cluster app values is empty")
	}
	if c.DefaultApps.Name == "" {
		return fmt.Errorf("default app name is empty")
	}
	if c.DefaultApps.AppName == "" {
		return fmt.Errorf("default app app name is empty")
	}
	if c.DefaultApps.Catalog == "" {
		return fmt.Errorf("default app catalog is empty")
	}
	if c.DefaultApps.Version == "" {
		return fmt.Errorf("default app version is empty")
	}
	if c.DefaultApps.Values == "" {
		return fmt.Errorf("default app values is empty")
	}
	if c.Provider.Name == "" {
		return fmt.Errorf("provider is empty")
	}
	if c.Provider.Name == key.ProviderVsphere {
		if c.Provider.CAPV.CloudConfig == "" {
			return fmt.Errorf("provider vsphere cloud config is empty")
		}
	} else if c.Provider.Name == key.ProviderAzure {
		if c.Provider.CAPZ.IdentityUA == "" {
			return fmt.Errorf("provider azure identity user-assigned is empty")
		}
		if c.Provider.CAPZ.IdentitySP == "" {
			return fmt.Errorf("provider azure identity service principal is empty")
		}
		if c.Provider.CAPZ.IdentityStaticSP == "" {
			return fmt.Errorf("provider azure identity static service principal is empty")
		}
	} else if c.Provider.Name == key.ProviderVCD {
		if c.Provider.CAPVCD.CloudConfig == "" {
			return fmt.Errorf("provider vcd cloud config is empty")
		}
	}
	if c.TaylorBotToken == "" {
		return fmt.Errorf("taylor bot token is empty")
	}
	if c.SSHdeployKey.Passphrase == "" {
		return fmt.Errorf("ssh deploy key passphrase is empty")
	}
	if c.SSHdeployKey.Identity == "" {
		return fmt.Errorf("ssh deploy key identity is empty")
	}
	if c.SSHdeployKey.KnownHosts == "" {
		return fmt.Errorf("ssh deploy key known hosts is empty")
	}
	if c.CustomerDeployKey.Passphrase == "" {
		return fmt.Errorf("customer deploy key passphrase is empty")
	}
	if c.CustomerDeployKey.Identity == "" {
		return fmt.Errorf("customer deploy key identity is empty")
	}
	if c.CustomerDeployKey.KnownHosts == "" {
		return fmt.Errorf("customer deploy key known hosts is empty")
	}
	if c.SharedDeployKey.Passphrase == "" {
		return fmt.Errorf("shared deploy key passphrase is empty")
	}
	if c.SharedDeployKey.Identity == "" {
		return fmt.Errorf("shared deploy key identity is empty")
	}
	if c.SharedDeployKey.KnownHosts == "" {
		return fmt.Errorf("shared deploy key known hosts is empty")
	}
	if c.CertManagerDNSChallenge.Enabled {
		if c.CertManagerDNSChallenge.Region == "" {
			return fmt.Errorf("cert manager dns challenge region is empty")
		}
		if c.CertManagerDNSChallenge.Role == "" {
			return fmt.Errorf("cert manager dns challenge role is empty")
		}
		if c.CertManagerDNSChallenge.AccessKeyID == "" {
			return fmt.Errorf("cert manager dns challenge access key id is empty")
		}
		if c.CertManagerDNSChallenge.SecretAccessKey == "" {
			return fmt.Errorf("cert manager dns challenge secret access key is empty")
		}
	}
	if c.ConfigureContainerRegistries.Enabled {
		if c.ConfigureContainerRegistries.Values == "" {
			return fmt.Errorf("configure container registries values is empty")
		}
	}
	if c.MCProxy.Enabled {
		if c.MCProxy.HTTPSProxy == "" {
			return fmt.Errorf("mc proxy https proxy is empty")
		}
	}
	return nil
}

func (c *CMC) Equals(desired *CMC) bool {
	return reflect.DeepEqual(c, desired)
}

func (c *CMC) SetDefaultAppValues() error {
	values, err := defaultappsvalues.GetDefaultAppsValuesFile(defaultappsvalues.Config{
		Cluster:                 c.Cluster,
		PrivateCA:               c.PrivateCA,
		Provider:                c.Provider.Name,
		CertManagerDNSChallenge: c.CertManagerDNSChallenge.Enabled,
	})
	if err != nil {
		return fmt.Errorf("failed to get default apps values.\n%w", err)
	}
	c.DefaultApps.Values = values
	return nil
}

func GetCMCFromFile(file string) (*CMC, error) {
	log.Debug().Msg(fmt.Sprintf("getting CMC object from file %s", file))
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read CMC file %s.\n%w", file, err)
	}
	return GetCMC(data)
}
