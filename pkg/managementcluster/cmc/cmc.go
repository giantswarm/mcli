package cmc

import (
	"bytes"
	"fmt"
	"os"
	"reflect"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"

	"github.com/giantswarm/mcli/pkg/key"
)

type CMC struct {
	AgePubKey                    string                       `yaml:"agePubKey"`
	Cluster                      string                       `yaml:"cluster"`
	ClusterApp                   App                          `yaml:"clusterApp"`
	DefaultApps                  App                          `yaml:"defaultApp"`
	MCAppsPreventDeletion        bool                         `yaml:"mcAppsPreventDeletion"`
	PrivateCA                    bool                         `yaml:"privateCA"`
	ClusterNamespace             string                       `yaml:"clusterNamespace"`
	Provider                     Provider                     `yaml:"provider"`
	TaylorBot                    TaylorBot                    `yaml:"taylorBot"`
	DeployKey                    DeployKey                    `yaml:"deployKey"`
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

type TaylorBot struct {
	User  string `yaml:"user"`
	Token string `yaml:"token"`
}

type DeployKey struct {
	Key        string `yaml:"key"`
	Identity   string `yaml:"identity"`
	KnownHosts string `yaml:"knownHosts"`
}

type MCProxy struct {
	Enabled  bool   `yaml:"enabled"`
	HostName string `yaml:"hostName,omitempty"`
	Port     int    `yaml:"port,omitempty"`
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
	fmt.Print(string(data))
	return nil
}

func (c *CMC) Override(override *CMC) *CMC {
	cmc := *c
	if override.AgePubKey != "" {
		cmc.AgePubKey = override.AgePubKey
	}
	if override.Cluster != "" {
		cmc.Cluster = override.Cluster
	}
	if override.ClusterApp.Name != "" {
		cmc.ClusterApp.Name = override.ClusterApp.Name
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
	if override.TaylorBot.User != "" {
		cmc.TaylorBot.User = override.TaylorBot.User
	}
	if override.TaylorBot.Token != "" {
		cmc.TaylorBot.Token = override.TaylorBot.Token
	}
	if override.DeployKey.Key != "" {
		cmc.DeployKey.Key = override.DeployKey.Key
	}
	if override.DeployKey.Identity != "" {
		cmc.DeployKey.Identity = override.DeployKey.Identity
	}
	if override.DeployKey.KnownHosts != "" {
		cmc.DeployKey.KnownHosts = override.DeployKey.KnownHosts
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
		if override.MCProxy.HostName != "" {
			cmc.MCProxy.HostName = override.MCProxy.HostName
		}
		if override.MCProxy.Port != 0 {
			cmc.MCProxy.Port = override.MCProxy.Port
		}
	}
	return &cmc
}

func (c *CMC) Validate() error {
	if c.AgePubKey == "" {
		return fmt.Errorf("age public key is empty")
	}
	if c.Cluster == "" {
		return fmt.Errorf("cluster is empty")
	}
	if c.ClusterApp.Name == "" {
		return fmt.Errorf("cluster app name is empty")
	}
	if c.ClusterApp.Catalog == "" {
		return fmt.Errorf("cluster app catalog is empty")
	}
	if c.ClusterApp.Version == "" {
		return fmt.Errorf("cluster app version is empty")
	}
	if c.DefaultApps.Name == "" {
		return fmt.Errorf("default app name is empty")
	}
	if c.DefaultApps.Catalog == "" {
		return fmt.Errorf("default app catalog is empty")
	}
	if c.DefaultApps.Version == "" {
		return fmt.Errorf("default app version is empty")
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
	if c.TaylorBot.User == "" {
		return fmt.Errorf("taylor bot user is empty")
	}
	if c.TaylorBot.Token == "" {
		return fmt.Errorf("taylor bot token is empty")
	}
	if c.DeployKey.Key == "" {
		return fmt.Errorf("deploy key is empty")
	}
	if c.DeployKey.Identity == "" {
		return fmt.Errorf("deploy key identity is empty")
	}
	if c.DeployKey.KnownHosts == "" {
		return fmt.Errorf("deploy key known hosts is empty")
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
		if c.MCProxy.HostName == "" {
			return fmt.Errorf("mc proxy host name is empty")
		}
		if c.MCProxy.Port == 0 {
			return fmt.Errorf("mc proxy port is empty")
		}
	}
	return nil
}

func (c *CMC) Equals(desired *CMC) bool {
	return reflect.DeepEqual(c, desired)
}

func GetCMCFromFile(file string) (*CMC, error) {
	log.Debug().Msg(fmt.Sprintf("getting CMC object from file %s", file))
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read CMC file %s.\n%w", file, err)
	}
	return GetCMC(data)
}
