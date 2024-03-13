package cmc

import (
	"fmt"

	"github.com/giantswarm/mcli/pkg/key"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc/age"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc/apps"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc/certmanager"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc/coredns"
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
	SopsFile              = ".sops.yaml"
	AgeKeyFile            = "age_secret_keys.yaml"
	ClusterAppsFile       = "cluster-app-manifests.yaml"
	DefaultAppsFile       = "default-apps-manifests.yaml"
	CertManagerFile       = "cert-manager-dns01-secret.yaml"
	KustomizationFile     = "kustomization.yaml"
	TaylorBotFile         = "github-giantswarm-https-credentials.yaml"
	SSHdeployKeyFile      = "giantswarm-clusters-ssh-credentials.yaml"
	CustomerDeployKeyFile = "configs-ssh-credentials.yaml"
	SharedDeployKeyFile   = "shared-configs-ssh-credentials.yaml"
	IssuerFile            = "private-cluster-issuer.yaml"
	RegistryFile          = "container-registries-configuration-secret.yaml"
	CoreDNSFile           = "coredns-configmap.yaml"
	DenyNetPolFile        = "deny-all-policies.yaml"
	AllowNetPolFile       = "networkpolicy-egress-to-proxy.yaml"
	SourceControllerFile  = "kustomization-post-build-proxy.yaml"
)

func GetCMCFromMap(data map[string]string, cluster string) (*CMC, error) {
	sopsConfig, err := sops.GetSopsConfig(data[SopsFile], cluster)
	if err != nil {
		return nil, fmt.Errorf("failed to get sops config.\n%w", err)
	}

	path := key.GetCMCPath(sopsConfig.Cluster)

	clusterAppsConfig, err := apps.GetClusterAppsConfig(data[fmt.Sprintf("%s/%s", path, ClusterAppsFile)])
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster apps config.\n%w", err)
	}
	defaultAppsConfig, err := apps.GetDefaultAppsConfig(data[fmt.Sprintf("%s/%s", path, DefaultAppsFile)])
	if err != nil {
		return nil, fmt.Errorf("failed to get default apps config.\n%w", err)
	}
	taylorBotToken, err := taylorbot.GetTaylorBotToken(data[fmt.Sprintf("%s/%s", path, TaylorBotFile)])
	if err != nil {
		return nil, fmt.Errorf("failed to get taylor bot token.\n%w", err)
	}
	sshDeployKeyConfig, err := deploykey.GetDeployKeyConfig(data[fmt.Sprintf("%s/%s", path, SSHdeployKeyFile)])
	if err != nil {
		return nil, fmt.Errorf("failed to get ssh deploy key config.\n%w", err)
	}
	customerDeployKeyConfig, err := deploykey.GetDeployKeyConfig(data[fmt.Sprintf("%s/%s", path, CustomerDeployKeyFile)])
	if err != nil {
		return nil, fmt.Errorf("failed to get customer deploy key config.\n%w", err)
	}
	sharedDeployKeyConfig, err := deploykey.GetDeployKeyConfig(data[fmt.Sprintf("%s/%s", path, SharedDeployKeyFile)])
	if err != nil {
		return nil, fmt.Errorf("failed to get shared deploy key config.\n%w", err)
	}
	ageKey, err := age.GetAgeKey(data[fmt.Sprintf("%s/%s", path, AgeKeyFile)], cluster)
	if err != nil {
		return nil, fmt.Errorf("failed to get age key.\n%w", err)
	}
	kustomizationConfig := kustomization.GetKustomizationConfig(data[fmt.Sprintf("%s/%s", path, KustomizationFile)])

	cmc := CMC{
		Cluster:   sopsConfig.Cluster,
		AgePubKey: sopsConfig.AgePubKey,
		AgeKey:    ageKey,
		ClusterApp: App{
			Name:    clusterAppsConfig.Name,
			AppName: clusterAppsConfig.AppName,
			Catalog: clusterAppsConfig.Catalog,
			Version: clusterAppsConfig.Version,
			Values:  clusterAppsConfig.Values,
		},
		DefaultApps: App{
			Name:    defaultAppsConfig.Name,
			AppName: defaultAppsConfig.AppName,
			Catalog: defaultAppsConfig.Catalog,
			Version: defaultAppsConfig.Version,
			Values:  defaultAppsConfig.Values,
		},
		ClusterNamespace: clusterAppsConfig.Namespace,
		Provider: Provider{
			Name: clusterAppsConfig.Provider,
		},
		PrivateCA:             kustomizationConfig.PrivateCA,
		DisableDenyAllNetPol:  kustomizationConfig.DisableDenyAllNetPol,
		MCAppsPreventDeletion: clusterAppsConfig.MCAppsPreventDeletion || defaultAppsConfig.MCAppsPreventDeletion,
		TaylorBotToken:        taylorBotToken,
		SSHdeployKey: DeployKey{
			Passphrase: sshDeployKeyConfig.Passphrase,
			Identity:   sshDeployKeyConfig.Identity,
			KnownHosts: sshDeployKeyConfig.KnownHosts,
		},
		CustomerDeployKey: DeployKey{
			Passphrase: customerDeployKeyConfig.Passphrase,
			Identity:   customerDeployKeyConfig.Identity,
			KnownHosts: customerDeployKeyConfig.KnownHosts,
		},
		SharedDeployKey: DeployKey{
			Passphrase: sharedDeployKeyConfig.Passphrase,
			Identity:   sharedDeployKeyConfig.Identity,
			KnownHosts: sharedDeployKeyConfig.KnownHosts,
		},
	}

	if kustomizationConfig.ConfigureContainerRegistries {
		registryConfig, err := registry.GetRegistryConfig(data[fmt.Sprintf("%s/%s", path, RegistryFile)])
		if err != nil {
			return nil, fmt.Errorf("failed to get registry config.\n%w", err)
		}
		cmc.ConfigureContainerRegistries = ConfigureContainerRegistries{
			Enabled: true,
			Values:  registryConfig,
		}
	}

	if kustomizationConfig.CustomCoreDNS {
		coreDNSConfig, err := coredns.GetCoreDNSValues(data[fmt.Sprintf("%s/%s", path, CoreDNSFile)])
		if err != nil {
			return nil, fmt.Errorf("failed to get coreDNS config.\n%w", err)
		}
		cmc.CustomCoreDNS = CustomCoreDNS{
			Enabled: true,
			Values:  coreDNSConfig,
		}
	}

	if kustomizationConfig.CertManagerDNSChallenge {
		certManagerConfig, err := certmanager.GetCertManagerConfig(data[fmt.Sprintf("%s/%s", path, CertManagerFile)])
		if err != nil {
			return nil, fmt.Errorf("failed to get cert manager config.\n%w", err)
		}
		cmc.CertManagerDNSChallenge = CertManagerDNSChallenge{
			Enabled:         true,
			Region:          certManagerConfig.Region,
			Role:            certManagerConfig.Role,
			AccessKeyID:     certManagerConfig.AccessKeyID,
			SecretAccessKey: certManagerConfig.SecretAccessKey,
		}
	}

	if kustomizationConfig.MCProxy {
		httpsProxy, err := mcproxy.GetHTTPSProxy(data[fmt.Sprintf("%s/%s", path, SourceControllerFile)])
		if err != nil {
			return nil, fmt.Errorf("failed to get https proxy.\n%w", err)
		}
		cmc.MCProxy = MCProxy{
			Enabled:  true,
			Hostname: httpsProxy.Hostname,
			Port:     httpsProxy.Port,
		}
	}

	if clusterAppsConfig.Provider == key.ProviderVsphere {
		capvConfig := capv.GetCAPVConfig(data[fmt.Sprintf("%s/%s", path, key.VsphereCredentialsFile)])
		cmc.Provider.CAPV.CloudConfig = capvConfig.CloudConfig
	} else if clusterAppsConfig.Provider == key.ProviderAzure {
		capzConfig := capz.GetCAPZConfig(
			data[fmt.Sprintf("%s/%s", path, key.AzureClusterIdentitySPFile)],
			data[fmt.Sprintf("%s/%s", path, key.AzureClusterIdentityUAFile)],
			data[fmt.Sprintf("%s/%s", path, key.AzureSecretClusterIdentityStaticSP)])
		cmc.Provider.CAPZ.IdentityUA = capzConfig.IdentityUA
		cmc.Provider.CAPZ.IdentitySP = capzConfig.IdentitySP
		cmc.Provider.CAPZ.IdentityStaticSP = capzConfig.IdentityStaticSP
	} else if clusterAppsConfig.Provider == key.ProviderVCD {
		capvcdConfig := capvcd.GetCAPVCDConfig(data[fmt.Sprintf("%s/%s", path, key.CloudDirectorCredentialsFile)])
		cmc.Provider.CAPVCD.CloudConfig = capvcdConfig.CloudConfig
	}
	return &cmc, nil
}

func (c *CMC) GetMap(cmcTemplate map[string]string) (map[string]string, error) {
	var err error
	path := key.GetCMCPath(c.Cluster)

	cmcTemplate, err = c.GetSops(cmcTemplate, path)
	if err != nil {
		return nil, fmt.Errorf("failed to add sops file.\n%w", err)
	}
	cmcTemplate, err = c.GetProviders(cmcTemplate, path)
	if err != nil {
		return nil, fmt.Errorf("failed to add provider files.\n%w", err)
	}
	cmcTemplate, err = c.GetPrivateCA(cmcTemplate, path)
	if err != nil {
		return nil, fmt.Errorf("failed to add private CA files.\n%w", err)
	}
	cmcTemplate, err = c.GetConfigureContainerRegistries(cmcTemplate, path)
	if err != nil {
		return nil, fmt.Errorf("failed to add configure container registries files.\n%w", err)
	}
	cmcTemplate, err = c.GetCustomCoreDNS(cmcTemplate, path)
	if err != nil {
		return nil, fmt.Errorf("failed to add custom coreDNS files.\n%w", err)
	}
	cmcTemplate, err = c.GetDenyNetPol(cmcTemplate, path)
	if err != nil {
		return nil, fmt.Errorf("failed to remove deny netpol files.\n%w", err)
	}
	cmcTemplate, err = c.GetMCProxy(cmcTemplate, path)
	if err != nil {
		return nil, fmt.Errorf("failed to add mcproxy files.\n%w", err)
	}
	cmcTemplate, err = c.GetClusterApps(cmcTemplate, path)
	if err != nil {
		return nil, fmt.Errorf("failed to add cluster apps files.\n%w", err)
	}
	cmcTemplate, err = c.GetDefaultApps(cmcTemplate, path)
	if err != nil {
		return nil, fmt.Errorf("failed to add default apps files.\n%w", err)
	}
	cmcTemplate, err = c.GetCertManager(cmcTemplate, path)
	if err != nil {
		return nil, fmt.Errorf("failed to add cert manager files.\n%w", err)
	}
	cmcTemplate, err = c.GetTaylorBot(cmcTemplate, path)
	if err != nil {
		return nil, fmt.Errorf("failed to add taylor bot files.\n%w", err)
	}
	cmcTemplate, err = c.GetDeployKey(cmcTemplate, path)
	if err != nil {
		return nil, fmt.Errorf("failed to add deploy key files.\n%w", err)
	}
	cmcTemplate, err = c.GetKustomization(cmcTemplate, path)
	if err != nil {
		return nil, fmt.Errorf("failed to add kustomization file.\n%w", err)
	}
	cmcTemplate, err = c.GetAge(cmcTemplate, path)
	if err != nil {
		return nil, fmt.Errorf("failed to add age file.\n%w", err)
	}

	return cmcTemplate, nil
}

// DeployKey
func (c *CMC) GetDeployKey(cmcTemplate map[string]string, path string) (map[string]string, error) {
	deployKeyFile, err := deploykey.GetDeployKeyFile(deploykey.Config{
		Name:       "giantswarm-clusters-ssh-credentials",
		Passphrase: c.SSHdeployKey.Passphrase,
		Identity:   c.SSHdeployKey.Identity,
		KnownHosts: c.SSHdeployKey.KnownHosts,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get ssh deploy key file.\n%w", err)
	}
	cmcTemplate[fmt.Sprintf("%s/%s", path, SSHdeployKeyFile)] = deployKeyFile

	deployKeyFile, err = deploykey.GetDeployKeyFile(deploykey.Config{
		Name:       "configs-ssh-credentials",
		Passphrase: c.CustomerDeployKey.Passphrase,
		Identity:   c.CustomerDeployKey.Identity,
		KnownHosts: c.CustomerDeployKey.KnownHosts,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get customer ssh deploy key file.\n%w", err)
	}
	cmcTemplate[fmt.Sprintf("%s/%s", path, CustomerDeployKeyFile)] = deployKeyFile

	deployKeyFile, err = deploykey.GetDeployKeyFile(deploykey.Config{
		Name:       "shared-configs-ssh-credentials",
		Passphrase: c.SharedDeployKey.Passphrase,
		Identity:   c.SharedDeployKey.Identity,
		KnownHosts: c.SharedDeployKey.KnownHosts,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get shared ssh deploy key file.\n%w", err)
	}
	cmcTemplate[fmt.Sprintf("%s/%s", path, SharedDeployKeyFile)] = deployKeyFile
	return cmcTemplate, nil
}

// Age
func (c *CMC) GetAge(cmcTemplate map[string]string, path string) (map[string]string, error) {
	ageFile, err := age.GetAgeFile(age.Config{
		Cluster: c.Cluster,
		AgeKey:  c.AgeKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get age file.\n%w", err)
	}
	cmcTemplate[fmt.Sprintf("%s/%s", path, AgeKeyFile)] = ageFile
	return cmcTemplate, nil
}

// Providers
func (c *CMC) GetProviders(cmcTemplate map[string]string, path string) (map[string]string, error) {
	if c.Provider.Name == key.ProviderVsphere {
		capvFile, err := capv.GetCAPVFile(capv.Config{
			Namespace:   c.ClusterNamespace,
			CloudConfig: c.Provider.CAPV.CloudConfig,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get CAPV file.\n%w", err)
		}
		cmcTemplate[fmt.Sprintf("%s/%s", path, key.VsphereCredentialsFile)] = capvFile
	} else {
		delete(cmcTemplate, fmt.Sprintf("%s/%s", path, key.VsphereCredentialsFile))
	}
	if c.Provider.Name == key.ProviderAzure {
		capzFile, err := capz.GetCAPZFile(capz.Config{
			Namespace:        c.ClusterNamespace,
			IdentityUA:       c.Provider.CAPZ.IdentityUA,
			IdentitySP:       c.Provider.CAPZ.IdentitySP,
			IdentityStaticSP: c.Provider.CAPZ.IdentityStaticSP,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get CAPZ file.\n%w", err)
		}
		cmcTemplate[fmt.Sprintf("%s/%s", path, key.AzureSecretClusterIdentityStaticSP)] = capzFile
		cmcTemplate[fmt.Sprintf("%s/%s", path, key.AzureClusterIdentitySPFile)] = c.Provider.CAPZ.IdentitySP
		cmcTemplate[fmt.Sprintf("%s/%s", path, key.AzureClusterIdentityUAFile)] = c.Provider.CAPZ.IdentityUA
	} else {
		delete(cmcTemplate, fmt.Sprintf("%s/%s", path, key.AzureSecretClusterIdentityStaticSP))
		delete(cmcTemplate, fmt.Sprintf("%s/%s", path, key.AzureClusterIdentitySPFile))
		delete(cmcTemplate, fmt.Sprintf("%s/%s", path, key.AzureClusterIdentityUAFile))
	}
	if c.Provider.Name == key.ProviderVCD {
		capvcdFile, err := capvcd.GetCAPVCDFile(capvcd.Config{
			Namespace:   c.ClusterNamespace,
			CloudConfig: c.Provider.CAPVCD.CloudConfig,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get CAPVCD file.\n%w", err)
		}
		cmcTemplate[fmt.Sprintf("%s/%s", path, key.CloudDirectorCredentialsFile)] = capvcdFile
	} else {
		delete(cmcTemplate, fmt.Sprintf("%s/%s", path, key.CloudDirectorCredentialsFile))
	}
	return cmcTemplate, nil
}

// PrivateCA
func (c *CMC) GetPrivateCA(cmcTemplate map[string]string, path string) (map[string]string, error) {
	if c.PrivateCA {
		issuerfile := issuer.GetIssuerFile()
		cmcTemplate[fmt.Sprintf("%s/%s", path, IssuerFile)] = issuerfile
	} else {
		delete(cmcTemplate, fmt.Sprintf("%s/%s", path, IssuerFile))
	}
	return cmcTemplate, nil
}

// CustomCoreDNS
func (c *CMC) GetCustomCoreDNS(cmcTemplate map[string]string, path string) (map[string]string, error) {
	if c.CustomCoreDNS.Enabled {
		coreDNSFile, err := coredns.GetCoreDNSFile(c.CustomCoreDNS.Values)
		if err != nil {
			return nil, fmt.Errorf("failed to get coreDNS file.\n%w", err)
		}
		cmcTemplate[fmt.Sprintf("%s/%s", path, CoreDNSFile)] = coreDNSFile
	} else {
		delete(cmcTemplate, fmt.Sprintf("%s/%s", path, CoreDNSFile))
	}
	return cmcTemplate, nil
}

// DenyNetPol
func (c *CMC) GetDenyNetPol(cmcTemplate map[string]string, path string) (map[string]string, error) {
	if c.DisableDenyAllNetPol {
		// remove deny all network policy entry from map
		delete(cmcTemplate, fmt.Sprintf("%s/%s", path, DenyNetPolFile))
	}
	return cmcTemplate, nil
}

// Sops
func (c *CMC) GetSops(cmcTemplate map[string]string, path string) (map[string]string, error) {
	sopsFile, err := sops.GetSopsFile(sops.Config{
		Cluster:   c.Cluster,
		AgePubKey: c.AgePubKey,
	}, cmcTemplate[SopsFile])
	if err != nil {
		return nil, fmt.Errorf("failed to get sops file.\n%w", err)
	}
	cmcTemplate[SopsFile] = sopsFile
	return cmcTemplate, nil
}

// ClusterApps
func (c *CMC) GetClusterApps(cmcTemplate map[string]string, path string) (map[string]string, error) {
	clusterAppsFile, err := apps.GetClusterAppsFile(apps.Config{
		Cluster:                      c.Cluster,
		Namespace:                    c.ClusterNamespace,
		ConfigureContainerRegistries: c.ConfigureContainerRegistries.Enabled,
		Provider:                     c.Provider.Name,
		MCAppsPreventDeletion:        c.MCAppsPreventDeletion,
		Name:                         c.ClusterApp.Name,
		AppName:                      c.ClusterApp.AppName,
		Catalog:                      c.ClusterApp.Catalog,
		Version:                      c.ClusterApp.Version,
		Values:                       c.ClusterApp.Values,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster apps file.\n%w", err)
	}
	cmcTemplate[fmt.Sprintf("%s/%s", path, ClusterAppsFile)] = clusterAppsFile
	//todo: how do we deal with the values files??

	return cmcTemplate, nil
}

// DefaultApps
func (c *CMC) GetDefaultApps(cmcTemplate map[string]string, path string) (map[string]string, error) {
	defaultAppsFile, err := apps.GetDefaultAppsFile(apps.Config{
		Cluster:               c.Cluster,
		Namespace:             c.ClusterNamespace,
		Provider:              c.Provider.Name,
		MCAppsPreventDeletion: c.MCAppsPreventDeletion,
		Name:                  c.DefaultApps.Name,
		AppName:               c.DefaultApps.AppName,
		Catalog:               c.DefaultApps.Catalog,
		Version:               c.DefaultApps.Version,
		Values:                c.DefaultApps.Values,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get default apps file.\n%w", err)
	}
	cmcTemplate[fmt.Sprintf("%s/%s", path, DefaultAppsFile)] = defaultAppsFile

	return cmcTemplate, nil
}

// CertManager
func (c *CMC) GetCertManager(cmcTemplate map[string]string, path string) (map[string]string, error) {
	if c.CertManagerDNSChallenge.Enabled {
		certManagerFile, err := certmanager.GetCertManagerFile(certmanager.Config{
			Cluster:          c.Cluster,
			ClusterNamespace: c.ClusterNamespace,
			Region:           c.CertManagerDNSChallenge.Region,
			Role:             c.CertManagerDNSChallenge.Role,
			AccessKeyID:      c.CertManagerDNSChallenge.AccessKeyID,
			SecretAccessKey:  c.CertManagerDNSChallenge.SecretAccessKey,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get cert-manager file.\n%w", err)
		}
		cmcTemplate[fmt.Sprintf("%s/%s", path, CertManagerFile)] = certManagerFile
	} else {
		delete(cmcTemplate, fmt.Sprintf("%s/%s", path, CertManagerFile))
	}
	return cmcTemplate, nil
}

// TaylorBot
func (c *CMC) GetTaylorBot(cmcTemplate map[string]string, path string) (map[string]string, error) {
	taylorBotToken, err := taylorbot.GetTaylorBotToken(c.TaylorBotToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get taylorbot file.\n%w", err)
	}
	cmcTemplate[fmt.Sprintf("%s/%s", path, TaylorBotFile)] = taylorBotToken
	return cmcTemplate, nil
}

// MCProxy
func (c *CMC) GetMCProxy(cmcTemplate map[string]string, path string) (map[string]string, error) {
	if c.MCProxy.Enabled {
		allowNetPolFile := mcproxy.GetAllowNetPolFile(mcproxy.Config{
			Hostname: c.MCProxy.Hostname,
			Port:     c.MCProxy.Port,
		})
		cmcTemplate[fmt.Sprintf("%s/%s", path, AllowNetPolFile)] = allowNetPolFile

		kustomization := mcproxy.GetKustomization(mcproxy.Config{
			Hostname: c.MCProxy.Hostname,
			Port:     c.MCProxy.Port,
		})
		cmcTemplate[fmt.Sprintf("%s/%s", path, SourceControllerFile)] = kustomization
	} else {
		delete(cmcTemplate, fmt.Sprintf("%s/%s", path, AllowNetPolFile))
		delete(cmcTemplate, fmt.Sprintf("%s/%s", path, SourceControllerFile))
	}
	return cmcTemplate, nil
}

// ConfigureContainerRegistries
func (c *CMC) GetConfigureContainerRegistries(cmcTemplate map[string]string, path string) (map[string]string, error) {
	if c.ConfigureContainerRegistries.Enabled {
		registryFile, err := registry.GetRegistryFile(c.ConfigureContainerRegistries.Values)
		if err != nil {
			return nil, fmt.Errorf("failed to get registry file.\n%w", err)
		}
		cmcTemplate[fmt.Sprintf("%s/%s", path, RegistryFile)] = registryFile
	} else {
		delete(cmcTemplate, fmt.Sprintf("%s/%s", path, RegistryFile))
	}
	return cmcTemplate, nil
}

// Kustomization
func (c *CMC) GetKustomization(cmcTemplate map[string]string, path string) (map[string]string, error) {
	kustomizationFile, err := kustomization.GetKustomizationFile(kustomization.Config{
		CertManagerDNSChallenge:      c.CertManagerDNSChallenge.Enabled,
		Provider:                     c.Provider.Name,
		PrivateCA:                    c.PrivateCA,
		ConfigureContainerRegistries: c.ConfigureContainerRegistries.Enabled,
		CustomCoreDNS:                c.CustomCoreDNS.Enabled,
		DisableDenyAllNetPol:         c.DisableDenyAllNetPol,
		MCProxy:                      c.MCProxy.Enabled,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get kustomization file.\n%w", err)
	}
	cmcTemplate[fmt.Sprintf("%s/%s", path, KustomizationFile)] = kustomizationFile
	return cmcTemplate, nil
}
