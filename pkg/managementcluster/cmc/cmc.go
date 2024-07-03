package cmc

import (
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
	Cluster                      string                       `yaml:"cluster"`
	ClusterApp                   App                          `yaml:"clusterApp"`
	DefaultApps                  App                          `yaml:"defaultApps,omitempty"`
	ClusterIntegratesDefaultApps bool                         `yaml:"clusterIntegratesDefaultApps"`
	MCAppsPreventDeletion        bool                         `yaml:"mcAppsPreventDeletion"`
	PrivateCA                    bool                         `yaml:"privateCA"`
	PrivateMC                    bool                         `yaml:"privateMC"`
	ClusterNamespace             string                       `yaml:"clusterNamespace"`
	Provider                     Provider                     `yaml:"provider"`
	TaylorBotToken               string                       `yaml:"taylorBotToken"`
	SSHdeployKey                 DeployKey                    `yaml:"sshDeployKey"`
	CustomerDeployKey            DeployKey                    `yaml:"customerDeployKey"`
	SharedDeployKey              DeployKey                    `yaml:"sharedDeployKey"`
	CertManagerDNSChallenge      CertManagerDNSChallenge      `yaml:"certManagerDNSChallenge"`
	ConfigureContainerRegistries ConfigureContainerRegistries `yaml:"configureContainerRegistries"`
	CustomCoreDNS                CustomCoreDNS                `yaml:"customCoreDNS"`
	DisableDenyAllNetPol         bool                         `yaml:"disableDenyAllNetPol"`
	MCProxy                      MCProxy                      `yaml:"mcProxy"`
	BaseDomain                   string                       `yaml:"baseDomain,omitempty"`
	RegistryDomain               string                       `yaml:"registryDomain,omitempty"`
	GitOps                       GitOps                       `yaml:"gitOps,omitempty"`
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

type GitOps struct {
	CMCRepository         string `yaml:"cmcRepository"`
	CMCBranch             string `yaml:"cmcBranch"`
	MCBBranchSource       string `yaml:"mcbBranchSource"`
	ConfigBranch          string `yaml:"configBranch"`
	MCAppCollectionBranch string `yaml:"mcAppCollectionBranch"`
}

type CAPV struct {
	CloudConfig string `yaml:"cloudConfig"`
}

type CAPZ struct {
	UAClientID     string `yaml:"uaClientID"`
	UATenantID     string `yaml:"uaTenantID"`
	UAResourceID   string `yaml:"uaResourceID"`
	ClientID       string `yaml:"clientID"`
	ClientSecret   string `yaml:"clientSecret"`
	TenantID       string `yaml:"tenantID"`
	SubscriptionID string `yaml:"subscriptionID"`
}

type CAPVCD struct {
	RefreshToken string `yaml:"refreshToken"`
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

type CustomCoreDNS struct {
	Enabled bool   `yaml:"enabled"`
	Values  string `yaml:"values,omitempty"`
}

type DeployKey struct {
	Passphrase string `yaml:"key"`
	Identity   string `yaml:"identity"`
	KnownHosts string `yaml:"knownHosts"`
}

type MCProxy struct {
	Enabled  bool   `yaml:"enabled"`
	Hostname string `yaml:"hostname,omitempty"`
	Port     string `yaml:"port,omitempty"`
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

	return key.GetData(c)
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
	if override.BaseDomain != "" {
		cmc.BaseDomain = override.BaseDomain
	}
	if override.RegistryDomain != "" {
		cmc.RegistryDomain = override.RegistryDomain
	}
	if override.GitOps.CMCRepository != "" {
		cmc.GitOps.CMCRepository = override.GitOps.CMCRepository
	}
	if override.GitOps.CMCBranch != "" {
		cmc.GitOps.CMCBranch = override.GitOps.CMCBranch
	}
	if override.GitOps.MCBBranchSource != "" {
		cmc.GitOps.MCBBranchSource = override.GitOps.MCBBranchSource
	}
	if override.GitOps.ConfigBranch != "" {
		cmc.GitOps.ConfigBranch = override.GitOps.ConfigBranch
	}
	if override.GitOps.MCAppCollectionBranch != "" {
		cmc.GitOps.MCAppCollectionBranch = override.GitOps.MCAppCollectionBranch
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
	if override.ClusterIntegratesDefaultApps {
		cmc.ClusterIntegratesDefaultApps = override.ClusterIntegratesDefaultApps
	}
	if override.MCAppsPreventDeletion {
		cmc.MCAppsPreventDeletion = override.MCAppsPreventDeletion
	}
	if override.PrivateCA {
		cmc.PrivateCA = override.PrivateCA
	}
	if override.PrivateMC {
		cmc.PrivateMC = override.PrivateMC
	}
	if override.ClusterNamespace != "" {
		cmc.ClusterNamespace = override.ClusterNamespace
	}
	if !override.ClusterIntegratesDefaultApps {
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
	} else {
		cmc.DefaultApps = App{}
	}
	if override.Provider.Name != "" {
		cmc.Provider.Name = override.Provider.Name
		if key.IsProviderVsphere(override.Provider.Name) {
			if override.Provider.CAPV.CloudConfig != "" {
				cmc.Provider.CAPV.CloudConfig = override.Provider.CAPV.CloudConfig
			}
		} else if key.IsProviderAzure(override.Provider.Name) {
			if override.Provider.CAPZ.UAClientID != "" {
				cmc.Provider.CAPZ.UAClientID = override.Provider.CAPZ.UAClientID
			}
			if override.Provider.CAPZ.UATenantID != "" {
				cmc.Provider.CAPZ.UATenantID = override.Provider.CAPZ.UATenantID
			}
			if override.Provider.CAPZ.UAResourceID != "" {
				cmc.Provider.CAPZ.UAResourceID = override.Provider.CAPZ.UAResourceID
			}
			if override.Provider.CAPZ.ClientID != "" {
				cmc.Provider.CAPZ.ClientID = override.Provider.CAPZ.ClientID
			}
			if override.Provider.CAPZ.ClientSecret != "" {
				cmc.Provider.CAPZ.ClientSecret = override.Provider.CAPZ.ClientSecret
			}
			if override.Provider.CAPZ.TenantID != "" {
				cmc.Provider.CAPZ.TenantID = override.Provider.CAPZ.TenantID
			}
			if override.Provider.CAPZ.SubscriptionID != "" {
				cmc.Provider.CAPZ.SubscriptionID = override.Provider.CAPZ.SubscriptionID
			}
		} else if key.IsProviderVCD(override.Provider.Name) {
			if override.Provider.CAPVCD.RefreshToken != "" {
				cmc.Provider.CAPVCD.RefreshToken = override.Provider.CAPVCD.RefreshToken
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
	if override.CustomCoreDNS.Enabled {
		cmc.CustomCoreDNS.Values = override.CustomCoreDNS.Values
	}
	if override.DisableDenyAllNetPol {
		cmc.DisableDenyAllNetPol = override.DisableDenyAllNetPol
	}
	if override.MCProxy.Enabled {
		cmc.MCProxy.Enabled = override.MCProxy.Enabled
		if override.MCProxy.Hostname != "" {
			cmc.MCProxy.Hostname = override.MCProxy.Hostname
		}
		if override.MCProxy.Port != "" {
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
	if c.ClusterNamespace == "" {
		return fmt.Errorf("cluster namespace is empty")
	}
	if c.BaseDomain == "" {
		return fmt.Errorf("base domain is empty")
	}
	if c.GitOps.CMCRepository == "" {
		return fmt.Errorf("gitops cmc repository is empty")
	}
	if c.GitOps.CMCBranch == "" {
		return fmt.Errorf("gitops cmc branch is empty")
	}
	if c.GitOps.MCBBranchSource == "" {
		return fmt.Errorf("gitops mcb branch source is empty")
	}
	if c.GitOps.ConfigBranch == "" {
		return fmt.Errorf("gitops config branch is empty")
	}
	if c.GitOps.MCAppCollectionBranch == "" {
		return fmt.Errorf("gitops mc app collection branch is empty")
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
	if c.Provider.Name == "" {
		return fmt.Errorf("provider is empty")
	}
	if !c.ClusterIntegratesDefaultApps {
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
	}
	if key.IsProviderVsphere(c.Provider.Name) {
		if c.Provider.CAPV.CloudConfig == "" {
			return fmt.Errorf("provider vsphere cloud config is empty")
		}
	} else if key.IsProviderAzure(c.Provider.Name) {
		if c.Provider.CAPZ.UAClientID == "" {
			return fmt.Errorf("provider azure ua client id is empty")
		}
		if c.Provider.CAPZ.UATenantID == "" {
			return fmt.Errorf("provider azure ua tenant id is empty")
		}
		if c.Provider.CAPZ.UAResourceID == "" {
			return fmt.Errorf("provider azure ua resource id is empty")
		}
		if c.Provider.CAPZ.ClientID == "" {
			return fmt.Errorf("provider azure client id is empty")
		}
		if c.Provider.CAPZ.ClientSecret == "" {
			return fmt.Errorf("provider azure client secret is empty")
		}
		if c.Provider.CAPZ.TenantID == "" {
			return fmt.Errorf("provider azure tenant id is empty")
		}
		if c.Provider.CAPZ.SubscriptionID == "" {
			return fmt.Errorf("provider azure subscription id is empty")
		}
	} else if key.IsProviderVCD(c.Provider.Name) {
		if c.Provider.CAPVCD.RefreshToken == "" {
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
	if c.CustomCoreDNS.Enabled {
		if c.CustomCoreDNS.Values == "" {
			return fmt.Errorf("custom core dns values is empty")
		}
	}
	if c.MCProxy.Enabled {
		if c.MCProxy.Hostname == "" {
			return fmt.Errorf("mc proxy hostname is empty")
		}
		if c.MCProxy.Port == "" {
			return fmt.Errorf("mc proxy port is empty")
		}
	}
	return nil
}

func (c *CMC) Equals(desired *CMC) bool {
	return reflect.DeepEqual(c, desired)
}

func (c *CMC) SetDefaultAppValues() error {
	config := defaultappsvalues.Config{
		Cluster:                 c.Cluster,
		PrivateCA:               c.PrivateCA,
		PrivateMC:               c.PrivateMC,
		Provider:                c.Provider.Name,
		CertManagerDNSChallenge: c.CertManagerDNSChallenge.Enabled,
	}
	if key.IsProviderAzure(config.Provider) && config.PrivateMC {
		config.IdentityClientID = c.Provider.CAPZ.UAClientID
		config.SubscriptionID = c.Provider.CAPZ.SubscriptionID
	}

	if c.ClusterIntegratesDefaultApps {
		values, err := defaultappsvalues.IntegrateDefaultAppsValuesInClusterValues(c.ClusterApp.Values, config)
		if err != nil {
			return fmt.Errorf("failed to integrate default apps values in cluster values.\n%w", err)
		}
		c.ClusterApp.Values = values
	} else {
		values, err := defaultappsvalues.GetDefaultAppsValuesFile(config)
		if err != nil {
			return fmt.Errorf("failed to get default apps values.\n%w", err)
		}
		c.DefaultApps.Values = values
	}
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
