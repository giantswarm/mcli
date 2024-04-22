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
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc/sopsfile"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc/taylorbot"
	"github.com/giantswarm/mcli/pkg/sops"
)

const (
	SopsFile = ".sops.yaml"
)

func GetCMCFromMap(input map[string]string, cluster string) (*CMC, error) {

	data, err := sops.DecryptDir(input)
	if err != nil {
		return nil, err
	}

	sopsConfig, err := sopsfile.GetSopsConfig(data[SopsFile], cluster)
	if err != nil {
		return nil, fmt.Errorf("failed to get sops config.\n%w", err)
	}

	path := key.GetCMCPath(sopsConfig.Cluster)

	clusterAppsConfig, err := apps.GetAppsConfig(data[fmt.Sprintf("%s/%s", path, kustomization.ClusterAppsFile)])
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster apps config.\n%w", err)
	}
	defaultAppsConfig, err := apps.GetAppsConfig(data[fmt.Sprintf("%s/%s", path, kustomization.DefaultAppsFile)])
	if err != nil {
		return nil, fmt.Errorf("failed to get default apps config.\n%w", err)
	}
	taylorBotToken, err := taylorbot.GetTaylorBotToken(data[fmt.Sprintf("%s/%s", path, kustomization.TaylorBotFile)])
	if err != nil {
		return nil, fmt.Errorf("failed to get taylor bot token.\n%w", err)
	}
	sshDeployKeyConfig, err := deploykey.GetDeployKeyConfig(data[fmt.Sprintf("%s/%s", path, kustomization.SSHdeployKeyFile)])
	if err != nil {
		return nil, fmt.Errorf("failed to get ssh deploy key config.\n%w", err)
	}
	customerDeployKeyConfig, err := deploykey.GetDeployKeyConfig(data[fmt.Sprintf("%s/%s", path, kustomization.CustomerDeployKeyFile)])
	if err != nil {
		return nil, fmt.Errorf("failed to get customer deploy key config.\n%w", err)
	}
	sharedDeployKeyConfig, err := deploykey.GetDeployKeyConfig(data[fmt.Sprintf("%s/%s", path, kustomization.SharedDeployKeyFile)])
	if err != nil {
		return nil, fmt.Errorf("failed to get shared deploy key config.\n%w", err)
	}
	kustomizationConfig, err := kustomization.GetKustomizationConfig(data[fmt.Sprintf("%s/%s", path, kustomization.KustomizationFile)])
	if err != nil {
		return nil, fmt.Errorf("failed to get kustomization config.\n%w", err)
	}

	cmc := CMC{
		Cluster:   sopsConfig.Cluster,
		AgePubKey: sopsConfig.AgePubKey,
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
		registryConfig, err := registry.GetRegistryConfig(data[fmt.Sprintf("%s/%s", path, kustomization.RegistryFile)])
		if err != nil {
			return nil, fmt.Errorf("failed to get registry config.\n%w", err)
		}
		cmc.ConfigureContainerRegistries = ConfigureContainerRegistries{
			Enabled: true,
			Values:  registryConfig,
		}
	}

	if kustomizationConfig.CustomCoreDNS {
		coreDNSConfig, err := coredns.GetCoreDNSValues(data[fmt.Sprintf("%s/%s", path, kustomization.CoreDNSFile)])
		if err != nil {
			return nil, fmt.Errorf("failed to get coreDNS config.\n%w", err)
		}
		cmc.CustomCoreDNS = CustomCoreDNS{
			Enabled: true,
			Values:  coreDNSConfig,
		}
	}

	if kustomizationConfig.CertManagerDNSChallenge {
		certManagerConfig, err := certmanager.GetCertManagerConfig(data[fmt.Sprintf("%s/%s", path, kustomization.CertManagerFile)])
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
		httpsProxy, err := mcproxy.GetHTTPSProxy(data[fmt.Sprintf("%s/%s", path, kustomization.SourceControllerFile)])
		if err != nil {
			return nil, fmt.Errorf("failed to get https proxy.\n%w", err)
		}
		cmc.MCProxy = MCProxy{
			Enabled:  true,
			Hostname: httpsProxy.Hostname,
			Port:     httpsProxy.Port,
		}
	}

	if key.IsProviderVsphere(clusterAppsConfig.Provider) {
		capvConfig, err := capv.GetCAPVConfig(data[fmt.Sprintf("%s/%s", path, kustomization.VsphereCredentialsFile)])
		if err != nil {
			return nil, fmt.Errorf("failed to get CAPV config.\n%w", err)
		}
		cmc.Provider.CAPV.CloudConfig = capvConfig
	} else if key.IsProviderAzure(clusterAppsConfig.Provider) {
		capzConfig, err := capz.GetCAPZConfig(
			data[fmt.Sprintf("%s/%s", path, kustomization.AzureClusterIdentitySPFile)],
			data[fmt.Sprintf("%s/%s", path, kustomization.AzureClusterIdentityUAFile)],
			data[fmt.Sprintf("%s/%s", path, kustomization.AzureSecretClusterIdentityStaticSP)])
		if err != nil {
			return nil, fmt.Errorf("failed to get CAPZ config.\n%w", err)
		}
		cmc.Provider.CAPZ.ClientID = capzConfig.ClientID
		cmc.Provider.CAPZ.ClientSecret = capzConfig.ClientSecret
		cmc.Provider.CAPZ.TenantID = capzConfig.TenantID
		cmc.Provider.CAPZ.UAClientID = capzConfig.UAClientID
		cmc.Provider.CAPZ.UATenantID = capzConfig.UATenantID
		cmc.Provider.CAPZ.UAResourceID = capzConfig.UAResourceID
	} else if key.IsProviderVCD(clusterAppsConfig.Provider) {
		refreshtoken, err := capvcd.GetCAPVCDConfig(data[fmt.Sprintf("%s/%s", path, kustomization.CloudDirectorCredentialsFile)])
		if err != nil {
			return nil, fmt.Errorf("failed to get CAPVCD config.\n%w", err)
		}
		cmc.Provider.CAPVCD.RefreshToken = refreshtoken
	}
	return &cmc, nil
}

func (c *CMC) GetMap(cmcTemplate map[string]string) (map[string]string, error) {
	var err error
	path := key.GetCMCPath(c.Cluster)

	c.EncodeSecrets()

	cmcTemplate, err = c.GetSops(cmcTemplate, path)
	if err != nil {
		return nil, fmt.Errorf("failed to add sops file.\n%w", err)
	}
	cmcTemplate, err = c.GetPrivateCA(cmcTemplate, path)
	if err != nil {
		return nil, fmt.Errorf("failed to add private CA files.\n%w", err)
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
	cmcTemplate, err = c.GetKustomization(cmcTemplate, path)
	if err != nil {
		return nil, fmt.Errorf("failed to add kustomization file.\n%w", err)
	}
	cmcTemplate, err = c.GetDeployKey(cmcTemplate, path)
	if err != nil {
		return nil, fmt.Errorf("failed to add deploy key files.\n%w", err)
	}
	cmcTemplate, err = c.GetSecrets(cmcTemplate, path)
	if err != nil {
		return nil, fmt.Errorf("failed to add secrets files.\n%w", err)
	}
	cmcTemplate, err = c.GetProviders(cmcTemplate, path)
	if err != nil {
		return nil, fmt.Errorf("failed to add provider files.\n%w", err)
	}

	err = c.DecodeSecrets()
	if err != nil {
		return nil, fmt.Errorf("failed to decode secrets.\n%w", err)
	}

	return cmcTemplate, nil
}

// DeployKey
func (c *CMC) GetDeployKey(cmcTemplate map[string]string, path string) (map[string]string, error) {
	deployKeyMap := map[string]string{}

	deployKeyFile, err := deploykey.GetDeployKeyFile(deploykey.Config{
		Name:       "giantswarm-clusters-ssh-credentials",
		Passphrase: c.SSHdeployKey.Passphrase,
		Identity:   c.SSHdeployKey.Identity,
		KnownHosts: c.SSHdeployKey.KnownHosts,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get ssh deploy key file.\n%w", err)
	}
	deployKeyMap[fmt.Sprintf("%s/%s", path, kustomization.SSHdeployKeyFile)] = deployKeyFile

	deployKeyFile, err = deploykey.GetDeployKeyFile(deploykey.Config{
		Name:       "configs-ssh-credentials",
		Passphrase: c.CustomerDeployKey.Passphrase,
		Identity:   c.CustomerDeployKey.Identity,
		KnownHosts: c.CustomerDeployKey.KnownHosts,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get customer ssh deploy key file.\n%w", err)
	}
	deployKeyMap[fmt.Sprintf("%s/%s", path, kustomization.CustomerDeployKeyFile)] = deployKeyFile

	deployKeyFile, err = deploykey.GetDeployKeyFile(deploykey.Config{
		Name:       "shared-configs-ssh-credentials",
		Passphrase: c.SharedDeployKey.Passphrase,
		Identity:   c.SharedDeployKey.Identity,
		KnownHosts: c.SharedDeployKey.KnownHosts,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get shared ssh deploy key file.\n%w", err)
	}
	deployKeyMap[fmt.Sprintf("%s/%s", path, kustomization.SharedDeployKeyFile)] = deployKeyFile

	deployKeyMap, err = sops.EncryptDir(deployKeyMap, c.AgePubKey)
	if err != nil {
		return nil, err
	}

	for k, v := range deployKeyMap {
		cmcTemplate[k] = v
	}
	return cmcTemplate, nil
}

func (c *CMC) GetSecrets(cmcTemplate map[string]string, path string) (map[string]string, error) {
	secretMap := map[string]string{}

	// Age
	ageFile, err := age.GetAgeFile(age.Config{
		Cluster: c.Cluster,
		AgeKey:  sops.GetAgeKey(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get age file.\n%w", err)
	}
	secretMap[fmt.Sprintf("%s/%s", path, kustomization.AgeKeyFile)] = ageFile

	// TaylorBot
	taylorBotToken, err := taylorbot.GetTaylorBotFile(c.TaylorBotToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get taylorbot file.\n%w", err)
	}
	secretMap[fmt.Sprintf("%s/%s", path, kustomization.TaylorBotFile)] = taylorBotToken

	// CertManager
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
		secretMap[fmt.Sprintf("%s/%s", path, kustomization.CertManagerFile)] = certManagerFile
	} else {
		delete(cmcTemplate, fmt.Sprintf("%s/%s", path, kustomization.CertManagerFile))
	}

	// ConfigureContainerRegistries
	if c.ConfigureContainerRegistries.Enabled {
		registryFile, err := registry.GetRegistryFile(c.ConfigureContainerRegistries.Values)
		if err != nil {
			return nil, fmt.Errorf("failed to get registry file.\n%w", err)
		}
		secretMap[fmt.Sprintf("%s/%s", path, kustomization.RegistryFile)] = registryFile
	} else {
		delete(cmcTemplate, fmt.Sprintf("%s/%s", path, kustomization.RegistryFile))
	}

	secretMap, err = sops.EncryptDir(secretMap, c.AgePubKey)
	if err != nil {
		return nil, err
	}

	for k, v := range secretMap {
		cmcTemplate[k] = v
	}
	return cmcTemplate, nil
}

// Providers
func (c *CMC) GetProviders(cmcTemplate map[string]string, path string) (map[string]string, error) {
	secretMap := map[string]string{}

	if key.IsProviderVsphere(c.Provider.Name) {
		capvFile, err := capv.GetCAPVFile(capv.Config{
			Namespace:   c.ClusterNamespace,
			CloudConfig: c.Provider.CAPV.CloudConfig,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get CAPV file.\n%w", err)
		}
		secretMap[fmt.Sprintf("%s/%s", path, kustomization.VsphereCredentialsFile)] = capvFile
	} else {
		delete(cmcTemplate, fmt.Sprintf("%s/%s", path, kustomization.VsphereCredentialsFile))
	}
	if key.IsProviderAzure(c.Provider.Name) {
		capzconfig := capz.Config{
			Namespace:    c.ClusterNamespace,
			ClientID:     c.Provider.CAPZ.ClientID,
			ClientSecret: c.Provider.CAPZ.ClientSecret,
			TenantID:     c.Provider.CAPZ.TenantID,
			UAClientID:   c.Provider.CAPZ.UAClientID,
			UATenantID:   c.Provider.CAPZ.UATenantID,
			UAResourceID: c.Provider.CAPZ.UAResourceID,
		}
		capzUAFile, err := capz.GetCAPZUAFile(capzconfig)
		if err != nil {
			return nil, fmt.Errorf("failed to get CAPZ file.\n%w", err)
		}
		capzSPFile, err := capz.GetCAPZSPFile(capzconfig)
		if err != nil {
			return nil, fmt.Errorf("failed to get CAPZ file.\n%w", err)
		}
		capzSecret, err := capz.GetCAPZSecret(capzconfig)
		if err != nil {
			return nil, fmt.Errorf("failed to get CAPZ file.\n%w", err)
		}
		secretMap[fmt.Sprintf("%s/%s", path, kustomization.AzureSecretClusterIdentityStaticSP)] = capzSecret
		cmcTemplate[fmt.Sprintf("%s/%s", path, kustomization.AzureClusterIdentitySPFile)] = capzSPFile
		cmcTemplate[fmt.Sprintf("%s/%s", path, kustomization.AzureClusterIdentityUAFile)] = capzUAFile
	} else {
		delete(cmcTemplate, fmt.Sprintf("%s/%s", path, kustomization.AzureSecretClusterIdentityStaticSP))
		delete(cmcTemplate, fmt.Sprintf("%s/%s", path, kustomization.AzureClusterIdentitySPFile))
		delete(cmcTemplate, fmt.Sprintf("%s/%s", path, kustomization.AzureClusterIdentityUAFile))
	}
	if key.IsProviderVCD(c.Provider.Name) {
		capvcdFile, err := capvcd.GetCAPVCDFile(capvcd.Config{
			Namespace:    c.ClusterNamespace,
			RefreshToken: c.Provider.CAPVCD.RefreshToken,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get CAPVCD file.\n%w", err)
		}
		secretMap[fmt.Sprintf("%s/%s", path, kustomization.CloudDirectorCredentialsFile)] = capvcdFile
	} else {
		delete(cmcTemplate, fmt.Sprintf("%s/%s", path, kustomization.CloudDirectorCredentialsFile))
	}

	secretMap, err := sops.EncryptDir(secretMap, c.AgePubKey)
	if err != nil {
		return nil, err
	}

	for k, v := range secretMap {
		cmcTemplate[k] = v
	}

	return cmcTemplate, nil
}

// PrivateCA
func (c *CMC) GetPrivateCA(cmcTemplate map[string]string, path string) (map[string]string, error) {
	if c.PrivateCA {
		issuerfile := issuer.GetIssuerFile()
		cmcTemplate[fmt.Sprintf("%s/%s", path, kustomization.IssuerFile)] = issuerfile
	} else {
		delete(cmcTemplate, fmt.Sprintf("%s/%s", path, kustomization.IssuerFile))
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
		cmcTemplate[fmt.Sprintf("%s/%s", path, kustomization.CoreDNSFile)] = coreDNSFile
	} else {
		delete(cmcTemplate, fmt.Sprintf("%s/%s", path, kustomization.CoreDNSFile))
	}
	return cmcTemplate, nil
}

// DenyNetPol
func (c *CMC) GetDenyNetPol(cmcTemplate map[string]string, path string) (map[string]string, error) {
	if c.DisableDenyAllNetPol {
		// remove deny all network policy entry from map
		delete(cmcTemplate, fmt.Sprintf("%s/%s", path, kustomization.DenyNetPolFile))
	}
	return cmcTemplate, nil
}

// Sops
func (c *CMC) GetSops(cmcTemplate map[string]string, path string) (map[string]string, error) {
	sopsFile, err := sopsfile.GetSopsFile(sopsfile.Config{
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
	cmcTemplate[fmt.Sprintf("%s/%s", path, kustomization.ClusterAppsFile)] = clusterAppsFile
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
	cmcTemplate[fmt.Sprintf("%s/%s", path, kustomization.DefaultAppsFile)] = defaultAppsFile

	return cmcTemplate, nil
}

// MCProxy
func (c *CMC) GetMCProxy(cmcTemplate map[string]string, path string) (map[string]string, error) {
	if c.MCProxy.Enabled {
		allowNetPolFile := mcproxy.GetAllowNetPolFile(mcproxy.Config{
			Hostname: c.MCProxy.Hostname,
			Port:     c.MCProxy.Port,
		})
		cmcTemplate[fmt.Sprintf("%s/%s", path, kustomization.AllowNetPolFile)] = allowNetPolFile

		proxykustomization := mcproxy.GetKustomization(mcproxy.Config{
			Hostname: c.MCProxy.Hostname,
			Port:     c.MCProxy.Port,
		})
		cmcTemplate[fmt.Sprintf("%s/%s", path, kustomization.SourceControllerFile)] = proxykustomization
	} else {
		delete(cmcTemplate, fmt.Sprintf("%s/%s", path, kustomization.AllowNetPolFile))
		delete(cmcTemplate, fmt.Sprintf("%s/%s", path, kustomization.SourceControllerFile))
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
	}, cmcTemplate[fmt.Sprintf("%s/%s", path, kustomization.KustomizationFile)])
	if err != nil {
		return nil, fmt.Errorf("failed to get kustomization file.\n%w", err)
	}
	cmcTemplate[fmt.Sprintf("%s/%s", path, kustomization.KustomizationFile)] = kustomizationFile
	return cmcTemplate, nil
}
