package cmc

import (
	"fmt"

	"github.com/giantswarm/mcli/pkg/key"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc/certmanager"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc/clusterapps"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc/coredns"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc/defaultapps"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc/deploykey"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc/issuer"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc/kustomization"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc/mcproxy"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc/provider/capv"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc/provider/capvcd"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc/provider/capz"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc/registry"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc/sops"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc/taylorbot"
)

const (
	SopsFile             = ".sops.yaml"
	ClusterAppsFile      = "cluster-app-manifests.yaml"
	DefaultAppsFile      = "default-apps-manifests.yaml"
	CertManagerFile      = "cert-manager-dns01-secret.yaml"
	KustomizationFile    = "kustomization.yaml"
	TaylorBotFile        = "github-giantswarm-https-credentials.yaml"
	DeployKeyFile        = "giantswarm-clusters-ssh-credentials.yaml"
	VSphereFile          = "vsphere-cloud-config-secret.yaml"
	CloudDirectorFile    = "cloud-director-cloud-config-secret.yaml"
	AzureFile            = "secret-clusteridentity-static-sp.yaml"
	IssuerFile           = "private-cluster-issuer.yaml"
	RegistryFile         = "container-registries-configuration-secret.yaml"
	CoreDNSFile          = "coredns-configmap.yaml"
	DenyNetPolFile       = "deny-all-policies.yaml"
	AllowNetPolFile      = "networkpolicy-egress-to-proxy.yaml"
	SourceControllerFile = "kustomization-post-build-proxy.yaml"
)

func GetCMCFromMap(data map[string]string) *CMC {
	sopsConfig := sops.GetSopsConfig(data[SopsFile])

	path := key.GetCMCPath(sopsConfig.Cluster)

	clusterAppsConfig := clusterapps.GetClusterAppsConfig(data[fmt.Sprintf("%s/%s", path, ClusterAppsFile)])
	defaultAppsConfig := defaultapps.GetDefaultAppsConfig(data[fmt.Sprintf("%s/%s", path, DefaultAppsFile)])
	taylorBotConfig := taylorbot.GetTaylorBotConfig(data[fmt.Sprintf("%s/%s", path, TaylorBotFile)])
	deployKeyConfig := deploykey.GetDeployKeyConfig(data[fmt.Sprintf("%s/%s", path, DeployKeyFile)])
	kustomizationConfig := kustomization.GetKustomizationConfig(data[fmt.Sprintf("%s/%s", path, KustomizationFile)])

	cmc := CMC{
		Cluster:   sopsConfig.Cluster,
		AgePubKey: sopsConfig.AgePubKey,
		ClusterApp: App{
			Name:    clusterAppsConfig.Name,
			Catalog: clusterAppsConfig.Catalog,
			Version: clusterAppsConfig.Version,
			Values:  clusterAppsConfig.Values,
		},
		DefaultApps: App{
			Name:    defaultAppsConfig.Name,
			Catalog: defaultAppsConfig.Catalog,
			Version: defaultAppsConfig.Version,
			Values:  defaultAppsConfig.Values,
		},
		ClusterNamespace: clusterAppsConfig.Namespace,
		Provider: Provider{
			Name: clusterAppsConfig.Provider,
		},
		PrivateCA:             kustomizationConfig.PrivateCA,
		CustomCoreDNS:         kustomizationConfig.CustomCoreDNS,
		DisableDenyAllNetPol:  kustomizationConfig.DisableDenyAllNetPol,
		MCAppsPreventDeletion: clusterAppsConfig.MCAppsPreventDeletion || defaultAppsConfig.MCAppsPreventDeletion,
		TaylorBot: TaylorBot{
			User:  taylorBotConfig.User,
			Token: taylorBotConfig.Token,
		},
		DeployKey: DeployKey{
			Key:        deployKeyConfig.Key,
			Identity:   deployKeyConfig.Identity,
			KnownHosts: deployKeyConfig.KnownHosts,
		},
	}

	if kustomizationConfig.ConfigureContainerRegistries {
		registryConfig := registry.GetRegistryConfig(data[fmt.Sprintf("%s/%s", path, RegistryFile)])
		cmc.ConfigureContainerRegistries = ConfigureContainerRegistries{
			Enabled: true,
			Values:  registryConfig.Values,
		}
	}

	if kustomizationConfig.CertManagerDNSChallenge {
		certManagerConfig := certmanager.GetCertManagerConfig(data[fmt.Sprintf("%s/%s", path, CertManagerFile)])
		cmc.CertManagerDNSChallenge = CertManagerDNSChallenge{
			Enabled:         true,
			Region:          certManagerConfig.Region,
			Role:            certManagerConfig.Role,
			AccessKeyID:     certManagerConfig.AccessKeyID,
			SecretAccessKey: certManagerConfig.SecretAccessKey,
		}
	}

	if kustomizationConfig.MCProxy {
		mcProxy := mcproxy.GetMCProxyConfig(data[fmt.Sprintf("%s/%s", path, AllowNetPolFile)], data[fmt.Sprintf("%s/%s", path, SourceControllerFile)])
		cmc.MCProxy = MCProxy{
			Enabled:  true,
			HostName: mcProxy.HostName,
			Port:     mcProxy.Port,
		}
	}

	if clusterAppsConfig.Provider == key.ProviderVsphere {
		capvConfig := capv.GetCAPVConfig(data[fmt.Sprintf("%s/%s", path, VSphereFile)])
		cmc.Provider.CAPV.CloudConfig = capvConfig.CloudConfig
	} else if clusterAppsConfig.Provider == key.ProviderAzure {
		capzConfig := capz.GetCAPZConfig(data[fmt.Sprintf("%s/%s", path, AzureFile)])
		cmc.Provider.CAPZ.IdentityUA = capzConfig.IdentityUA
		cmc.Provider.CAPZ.IdentitySP = capzConfig.IdentitySP
		cmc.Provider.CAPZ.IdentityStaticSP = capzConfig.IdentityStaticSP
	} else if clusterAppsConfig.Provider == key.ProviderVCD {
		capvcdConfig := capvcd.GetCAPVCDConfig(data[fmt.Sprintf("%s/%s", path, CloudDirectorFile)])
		cmc.Provider.CAPVCD.CloudConfig = capvcdConfig.CloudConfig
	}
	return &cmc
}

func (c *CMC) GetMap(cmcTemplate map[string]string) (map[string]string, error) {
	sopsFile, err := c.GetSopsFile()
	if err != nil {
		return nil, fmt.Errorf("failed to get sops file.\n%w", err)
	}
	cmcTemplate[SopsFile] = sopsFile

	path := key.GetCMCPath(c.Cluster)

	clusterAppsFile, err := c.GetClusterAppsFile()
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster apps file.\n%w", err)
	}
	cmcTemplate[fmt.Sprintf("%s/%s", path, ClusterAppsFile)] = clusterAppsFile
	//todo: how do we deal with the values files??

	defaultAppsFile, err := c.GetDefaultAppsFile()
	if err != nil {
		return nil, fmt.Errorf("failed to get default apps file.\n%w", err)
	}
	cmcTemplate[fmt.Sprintf("%s/%s", path, DefaultAppsFile)] = defaultAppsFile

	if c.CertManagerDNSChallenge.Enabled {
		certManagerFile, err := c.GetCertManagerFile()
		if err != nil {
			return nil, fmt.Errorf("failed to get cert-manager file.\n%w", err)
		}
		cmcTemplate[fmt.Sprintf("%s/%s", path, CertManagerFile)] = certManagerFile
	}

	taylorBotFile, err := c.GetTaylorBotFile()
	if err != nil {
		return nil, fmt.Errorf("failed to get taylorbot file.\n%w", err)
	}
	cmcTemplate[fmt.Sprintf("%s/%s", path, TaylorBotFile)] = taylorBotFile

	deployKeyFile, err := c.GetDeployKeyFile()
	if err != nil {
		return nil, fmt.Errorf("failed to get deploy key file.\n%w", err)
	}
	cmcTemplate[fmt.Sprintf("%s/%s", path, DeployKeyFile)] = deployKeyFile

	if c.Provider.Name == key.ProviderVsphere {
		capvFile, err := c.GetCAPVFile()
		if err != nil {
			return nil, fmt.Errorf("failed to get CAPV file.\n%w", err)
		}
		cmcTemplate[fmt.Sprintf("%s/%s", path, VSphereFile)] = capvFile
	} else if c.Provider.Name == key.ProviderAzure {
		capzFile, err := c.GetCAPZFile()
		if err != nil {
			return nil, fmt.Errorf("failed to get CAPZ file.\n%w", err)
		}
		cmcTemplate[fmt.Sprintf("%s/%s", path, AzureFile)] = capzFile
	} else if c.Provider.Name == key.ProviderVCD {
		capvcdFile, err := c.GetCAPVCDFile()
		if err != nil {
			return nil, fmt.Errorf("failed to get CAPVCD file.\n%w", err)
		}
		cmcTemplate[fmt.Sprintf("%s/%s", path, CloudDirectorFile)] = capvcdFile
	}
	if c.PrivateCA {
		issuerFile, err := c.GetIssuerFile()
		if err != nil {
			return nil, fmt.Errorf("failed to get issuer file.\n%w", err)
		}
		cmcTemplate[fmt.Sprintf("%s/%s", path, IssuerFile)] = issuerFile
	}
	if c.ConfigureContainerRegistries.Enabled {
		registryFile, err := c.GetRegistryFile()
		if err != nil {
			return nil, fmt.Errorf("failed to get registry file.\n%w", err)
		}
		cmcTemplate[fmt.Sprintf("%s/%s", path, RegistryFile)] = registryFile
	}
	if c.CustomCoreDNS {
		coreDNSFile, err := c.GetCoreDNSFile()
		if err != nil {
			return nil, fmt.Errorf("failed to get coreDNS file.\n%w", err)
		}
		cmcTemplate[fmt.Sprintf("%s/%s", path, CoreDNSFile)] = coreDNSFile
	}
	if c.DisableDenyAllNetPol {
		// remove deny all network policy entry from map
		delete(cmcTemplate, fmt.Sprintf("%s/%s", path, DenyNetPolFile))
	}
	if c.MCProxy.Enabled {
		allowNetPolFile, err := c.GetAllowNetPolFile()
		if err != nil {
			return nil, fmt.Errorf("failed to get allow netpol file.\n%w", err)
		}
		cmcTemplate[fmt.Sprintf("%s/%s", path, AllowNetPolFile)] = allowNetPolFile

		sourceControllerFile, err := c.GetSourceControllerFile()
		if err != nil {
			return nil, fmt.Errorf("failed to get source controller file.\n%w", err)
		}
		cmcTemplate[fmt.Sprintf("%s/%s", path, SourceControllerFile)] = sourceControllerFile
	}

	return cmcTemplate, nil
}

func (c *CMC) GetSopsFile() (string, error) {
	return sops.GetSopsFile(sops.Config{
		Cluster:   c.Cluster,
		AgePubKey: c.AgePubKey,
	})
}

func (c *CMC) GetClusterAppsFile() (string, error) {
	return clusterapps.GetClusterAppsFile(clusterapps.Config{
		Cluster:                      c.Cluster,
		Namespace:                    c.ClusterNamespace,
		ConfigureContainerRegistries: c.ConfigureContainerRegistries.Enabled,
		Provider:                     c.Provider.Name,
		MCAppsPreventDeletion:        c.MCAppsPreventDeletion,
		Name:                         c.ClusterApp.Name,
		Catalog:                      c.ClusterApp.Catalog,
		Version:                      c.ClusterApp.Version,
		Values:                       c.ClusterApp.Values,
	})
}

func (c *CMC) GetDefaultAppsFile() (string, error) {
	return defaultapps.GetDefaultAppsFile(defaultapps.Config{
		Cluster:                 c.Cluster,
		Namespace:               c.ClusterNamespace,
		PrivateCA:               c.PrivateCA,
		Provider:                c.Provider.Name,
		CertManagerDNSChallenge: c.CertManagerDNSChallenge.Enabled,
		MCAppsPreventDeletion:   c.MCAppsPreventDeletion,
		Name:                    c.DefaultApps.Name,
		Catalog:                 c.DefaultApps.Catalog,
		Version:                 c.DefaultApps.Version,
		Values:                  c.DefaultApps.Values,
	})
}

func (c *CMC) GetCertManagerFile() (string, error) {
	return certmanager.GetCertManagerFile(certmanager.Config{
		Region:          c.CertManagerDNSChallenge.Region,
		Role:            c.CertManagerDNSChallenge.Role,
		AccessKeyID:     c.CertManagerDNSChallenge.AccessKeyID,
		SecretAccessKey: c.CertManagerDNSChallenge.SecretAccessKey,
	})
}

func (c *CMC) GetTaylorBotFile() (string, error) {
	return taylorbot.GetTaylorBotFile(taylorbot.Config{
		User:      c.TaylorBot.User,
		Token:     c.TaylorBot.Token,
		AgePubKey: c.AgePubKey,
	})
}

func (c *CMC) GetDeployKeyFile() (string, error) {
	return deploykey.GetDeployKeyFile(deploykey.Config{
		Key:        c.DeployKey.Key,
		Identity:   c.DeployKey.Identity,
		KnownHosts: c.DeployKey.KnownHosts,
		AgePubKey:  c.AgePubKey,
	})
}

func (c *CMC) GetCAPVFile() (string, error) {
	return capv.GetCAPVFile(capv.Config{
		Namespace:   c.ClusterNamespace,
		CloudConfig: c.Provider.CAPV.CloudConfig,
		AgePubKey:   c.AgePubKey,
	})
}

func (c *CMC) GetCAPZFile() (string, error) {
	return capz.GetCAPZFile(capz.Config{
		Namespace:        c.ClusterNamespace,
		IdentityUA:       c.Provider.CAPZ.IdentityUA,
		IdentitySP:       c.Provider.CAPZ.IdentitySP,
		IdentityStaticSP: c.Provider.CAPZ.IdentityStaticSP,
		AgePubKey:        c.AgePubKey,
	})
}

func (c *CMC) GetCAPVCDFile() (string, error) {
	return capvcd.GetCAPVCDFile(capvcd.Config{
		Namespace:   c.ClusterNamespace,
		CloudConfig: c.Provider.CAPVCD.CloudConfig,
		AgePubKey:   c.AgePubKey,
	})
}

func (c *CMC) GetIssuerFile() (string, error) {
	return issuer.GetIssuerFile()
}

func (c *CMC) GetRegistryFile() (string, error) {
	return registry.GetRegistryFile(registry.Config{
		Values:    c.ConfigureContainerRegistries.Values,
		AgePubKey: c.AgePubKey,
	})
}

func (c *CMC) GetCoreDNSFile() (string, error) {
	return coredns.GetCoreDNSFile()
}

func (c *CMC) GetAllowNetPolFile() (string, error) {
	return mcproxy.GetAllowNetPolFile(mcproxy.Config{
		HostName: c.MCProxy.HostName,
		Port:     c.MCProxy.Port,
	})
}

func (c *CMC) GetSourceControllerFile() (string, error) {
	return mcproxy.GetSourceControllerFile(mcproxy.Config{
		HostName: c.MCProxy.HostName,
		Port:     c.MCProxy.Port,
	})
}

func (c *CMC) GetKustomizationFile() (string, error) {
	return kustomization.GetKustomizationFile(kustomization.Config{
		CertManagerDNSChallenge:      c.CertManagerDNSChallenge.Enabled,
		Provider:                     c.Provider.Name,
		PrivateCA:                    c.PrivateCA,
		ConfigureContainerRegistries: c.ConfigureContainerRegistries.Enabled,
		CustomCoreDNS:                c.CustomCoreDNS,
		DisableDenyAllNetPol:         c.DisableDenyAllNetPol,
		MCProxy:                      c.MCProxy.Enabled,
	})
}
