package pushcmc

import (
	"context"
	"fmt"
	"strings"

	"github.com/giantswarm/mcli/pkg/github"
	"github.com/giantswarm/mcli/pkg/key"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Cluster       string
	Github        *github.Github
	CMCRepository string
	CMCBranch     string
	Provider      string
	Input         *cmc.CMC
	Flags         CMCFlags
}

type CMCFlags struct {
	Secrets                      SecretFlags
	SecretFolder                 string
	AgePubKey                    string
	AgeKey                       string
	TaylorBotToken               string
	MCAppsPreventDeletion        bool
	ClusterAppName               string
	ClusterAppCatalog            string
	ClusterAppVersion            string
	ClusterNamespace             string
	ConfigureContainerRegistries bool
	DefaultAppsName              string
	DefaultAppsCatalog           string
	DefaultAppsVersion           string
	PrivateCA                    bool
	CertManagerDNSChallenge      bool
	MCCustomCoreDNSConfig        string
	MCProxyEnabled               bool
	MCHTTPSProxy                 string
}

type SecretFlags struct {
	SSHDeployKey                      DeployKey
	CustomerDeployKey                 DeployKey
	SharedDeployKey                   DeployKey
	VSphereCredentials                string
	CloudDirectorRefreshToken         string
	Azure                             AzureFlags
	ContainerRegistryConfiguration    string
	ClusterValues                     string
	CertManagerRoute53Region          string
	CertManagerRoute53Role            string
	CertManagerRoute53AccessKeyID     string
	CertManagerRoute53SecretAccessKey string
}

type AzureFlags struct {
	UAClientID   string
	UATenantID   string
	UAResourceID string
	ClientID     string
	TenantID     string
	ClientSecret string
}

type DeployKey struct {
	Passphrase string
	Identity   string
	KnownHosts string
}

func (c *Config) Run(ctx context.Context) (*cmc.CMC, error) {
	err := c.Validate()
	if err != nil {
		return nil, err
	}
	err = c.ReadSecretFlags()
	if err != nil {
		return nil, fmt.Errorf("failed to set secret flags.\n%w", err)
	}
	return c.PushCMC(ctx)
}

func (c *Config) PushCMC(ctx context.Context) (*cmc.CMC, error) {
	if err := c.Branch(ctx); err != nil {
		return nil, err
	}
	// pulling current cmc
	cmc, err := c.Pull(ctx)
	if err != nil {
		// if there is no current cmc, create a new one
		// check if the error is a github.ErrNotFound
		if github.IsNotFound(err) {
			log.Debug().Msg(fmt.Sprintf("no current %s entry for %s found, creating a new one", c.CMCRepository, c.Cluster))
			return c.Create(ctx)
		} else {
			return nil, fmt.Errorf("failed to pull %s entry for %s.\n%w", c.CMCRepository, c.Cluster, err)
		}
	}
	cmc, err = decode(c.Flags.AgePubKey, cmc)
	if err != nil {
		return nil, fmt.Errorf("failed to decode cmc.\n%w", err)
	}
	return c.Update(ctx, cmc)
}

func (c *Config) Branch(ctx context.Context) error {
	log.Debug().Msg(fmt.Sprintf("getting %s branch %s", c.CMCRepository, c.CMCBranch))

	cmcRepository := github.Repository{
		Github:       c.Github,
		Name:         c.CMCRepository,
		Organization: key.OrganizationGiantSwarm,
		Branch:       c.CMCBranch,
	}
	err := cmcRepository.CheckBranch(ctx)
	if err != nil {
		if github.IsNotFound(err) {
			log.Debug().Msg(fmt.Sprintf("%s branch %s not found, creating it", c.CMCRepository, c.CMCBranch))
			err = cmcRepository.CreateBranch(ctx, key.CMCMainBranch)
			if err != nil {
				return fmt.Errorf("failed to create %s branch %s.\n%w", c.CMCRepository, c.CMCBranch, err)
			}
		} else {
			return fmt.Errorf("failed to check %s branch %s.\n%w", c.CMCRepository, c.CMCBranch, err)
		}
	}
	return nil
}

func (c *Config) Pull(ctx context.Context) (map[string]string, error) {
	log.Debug().Msgf("pulling current %s entry for %s", c.CMCRepository, c.Cluster)

	cmcRepository := github.Repository{
		Github:       c.Github,
		Name:         c.CMCRepository,
		Organization: key.OrganizationGiantSwarm,
		Branch:       c.CMCBranch,
	}
	if err := cmcRepository.Check(ctx); err != nil {
		return nil, err
	}
	data, err := cmcRepository.GetDirectory(ctx, key.GetCMCPath(c.Cluster))
	if err != nil {
		return nil, err
	}
	sops, err := cmcRepository.GetFile(ctx, cmc.SopsFile)
	if err != nil {
		return nil, err
	}
	data[cmc.SopsFile] = sops
	return data, nil
}

func (c *Config) Create(ctx context.Context) (*cmc.CMC, error) {
	var err error
	log.Debug().Msg(fmt.Sprintf("creating new %s entry for %s", c.CMCRepository, c.Cluster))
	var desiredCMC *cmc.CMC
	{
		if c.Input == nil {
			desiredCMC, err = getNewCMCFromFlags(*c)
			if err != nil {
				return nil, fmt.Errorf("failed to get new %s object from flags.\n%w", c.CMCRepository, err)
			}
		} else {
			desiredCMC = c.Input
		}
	}
	template, err := c.PullTemplate()
	if err != nil {
		return nil, fmt.Errorf("failed to pull template.\n%w", err)
	}
	create, err := desiredCMC.GetMap(template)
	if err != nil {
		return nil, fmt.Errorf("failed to get cmc map.\n%w", err)
	}
	return c.Push(ctx, create)
}

func (c *Config) Update(ctx context.Context, currentCMCmap map[string]string) (*cmc.CMC, error) {
	var err error
	log.Debug().Msg(fmt.Sprintf("updating %s entry for %s", c.CMCRepository, c.Cluster))
	currentCMC, err := cmc.GetCMCFromMap(currentCMCmap, c.Cluster)
	if err != nil {
		return nil, fmt.Errorf("failed to get cmc from map.\n%w", err)
	}

	var desiredCMC *cmc.CMC
	{
		if c.Input == nil {
			desiredCMC, err = overrideCMCWithFlags(currentCMC, *c)
			if err != nil {
				return nil, fmt.Errorf("failed to override cmc with flags.\n%w", err)
			}
		} else {
			desiredCMC = currentCMC.Override(c.Input)
		}
	}
	if currentCMC.Equals(desiredCMC) {
		log.Debug().Msg(fmt.Sprintf("%s entry for %s is up to date", c.CMCRepository, c.Cluster))
		return desiredCMC, nil
	}
	// if deny all netpol is being enabled, we need to get the file from the template
	if currentCMC.DisableDenyAllNetPol && !desiredCMC.DisableDenyAllNetPol {
		template, err := c.PullTemplateFile(cmc.DenyNetPolFile)
		if err != nil {
			return nil, fmt.Errorf("failed to pull template file.\n%w", err)
		}
		currentCMCmap[fmt.Sprintf("%s/%s", key.GetCMCPath(c.Cluster), cmc.DenyNetPolFile)] = template
	}
	update, err := desiredCMC.GetMap(currentCMCmap)
	if err != nil {
		return nil, fmt.Errorf("failed to get cmc map.\n%w", err)
	}
	update, err = encode(desiredCMC.AgePubKey, update)
	if err != nil {
		return nil, fmt.Errorf("failed to encode cmc map.\n%w", err)
	}
	return c.Push(ctx, update)
}

func (c *Config) Push(ctx context.Context, desiredCMC map[string]string) (*cmc.CMC, error) {
	log.Debug().Msg(fmt.Sprintf("pushing %s entry for %s", c.CMCRepository, c.Cluster))

	cmcRepository := github.Repository{
		Github:       c.Github,
		Name:         c.CMCRepository,
		Organization: key.OrganizationGiantSwarm,
		Branch:       c.CMCBranch,
	}
	err := cmcRepository.Check(ctx)
	if err != nil {
		return nil, err
	}
	// TODO actually push the CMC
	log.Debug().Msg("This is a dry run. No changes will be made.")

	return cmc.GetCMCFromMap(desiredCMC, c.Cluster)
}

func (c *Config) Validate() error {
	if c.CMCBranch == "main" || c.CMCBranch == "master" {
		return fmt.Errorf("cannot push to cmc branch %s\n%w", c.CMCBranch, ErrInvalidFlag)
	}
	if c.Input != nil {
		log.Debug().Msg("using input file. Other cmc flags will be ignored")
		return nil
	}
	if c.CMCRepository == "" {
		return fmt.Errorf("cmc repository is required\n%w", ErrInvalidFlag)
	}
	if c.Provider != "" {
		if !key.IsValidProvider(c.Provider) {
			return fmt.Errorf("invalid provider %s. Valid values: %s:\n%w", c.Provider, key.GetValidProviders(), ErrInvalidFlag)
		}
	}
	return nil
}

func getNewCMCFromFlags(c Config) (*cmc.CMC, error) {
	// Ensure that all the needed flags are set
	if c.Flags.ClusterAppName == "" ||
		c.Flags.ClusterAppCatalog == "" ||
		c.Flags.ClusterAppVersion == "" ||
		c.Flags.DefaultAppsName == "" ||
		c.Flags.DefaultAppsCatalog == "" ||
		c.Flags.DefaultAppsVersion == "" ||
		c.Provider == "" ||
		c.Cluster == "" ||
		c.Flags.AgePubKey == "" ||
		c.Flags.AgeKey == "" ||
		c.Flags.TaylorBotToken == "" {
		return nil, fmt.Errorf("not all required flags are set\n%w", ErrInvalidFlag)
	}

	// Ensure that all the needed secret values are set
	if c.Flags.Secrets.SSHDeployKey.Identity == "" ||
		c.Flags.Secrets.SSHDeployKey.KnownHosts == "" ||
		c.Flags.Secrets.SSHDeployKey.Passphrase == "" ||
		c.Flags.Secrets.CustomerDeployKey.Identity == "" ||
		c.Flags.Secrets.CustomerDeployKey.KnownHosts == "" ||
		c.Flags.Secrets.CustomerDeployKey.Passphrase == "" ||
		c.Flags.Secrets.SharedDeployKey.Identity == "" ||
		c.Flags.Secrets.SharedDeployKey.KnownHosts == "" ||
		c.Flags.Secrets.SharedDeployKey.Passphrase == "" ||
		c.Flags.Secrets.ClusterValues == "" {
		return nil, fmt.Errorf("not all required values are set\n%w", ErrInvalidFlag)
	}
	// Ensure that needed values for enabled features are set
	if c.Flags.ConfigureContainerRegistries {
		if c.Flags.Secrets.ContainerRegistryConfiguration == "" {
			return nil, fmt.Errorf("container registry configuration is required\n%w", ErrInvalidFlag)
		}
	}
	if c.Flags.CertManagerDNSChallenge {
		if c.Flags.Secrets.CertManagerRoute53Region == "" ||
			c.Flags.Secrets.CertManagerRoute53Role == "" ||
			c.Flags.Secrets.CertManagerRoute53AccessKeyID == "" ||
			c.Flags.Secrets.CertManagerRoute53SecretAccessKey == "" {
			return nil, fmt.Errorf("cert manager dns challenge configuration is required\n%w", ErrInvalidFlag)
		}
	}
	if c.Flags.MCProxyEnabled {
		if c.Flags.MCHTTPSProxy == "" {
			return nil, fmt.Errorf("mc proxy is enabled but no https proxy is set\n%w", ErrInvalidFlag)
		}
	}
	// Ensure that needed values for the provider are set
	if c.Provider == key.ProviderVsphere {
		if c.Flags.Secrets.VSphereCredentials == "" {
			return nil, fmt.Errorf("vsphere credentials are required\n%w", ErrInvalidFlag)
		}
	} else if c.Provider == key.ProviderVCD {
		if c.Flags.Secrets.CloudDirectorRefreshToken == "" {
			return nil, fmt.Errorf("cloud director refresh token is required\n%w", ErrInvalidFlag)
		}
	} else if c.Provider == key.ProviderAzure {
		if c.Flags.Secrets.Azure.ClientID == "" ||
			c.Flags.Secrets.Azure.TenantID == "" ||
			c.Flags.Secrets.Azure.ClientSecret == "" ||
			c.Flags.Secrets.Azure.UAClientID == "" ||
			c.Flags.Secrets.Azure.UATenantID == "" ||
			c.Flags.Secrets.Azure.UAResourceID == "" {
			return nil, fmt.Errorf("azure credentials are required\n%w", ErrInvalidFlag)
		}
	}
	return getCMC(c)
}

func getCMC(c Config) (*cmc.CMC, error) {
	newCMC := &cmc.CMC{
		AgePubKey:        c.Flags.AgePubKey,
		AgeKey:           c.Flags.AgeKey,
		Cluster:          c.Cluster,
		ClusterNamespace: c.Flags.ClusterNamespace,
		ClusterApp: cmc.App{
			Name:    c.Flags.ClusterAppName,
			AppName: c.Cluster,
			Catalog: c.Flags.ClusterAppCatalog,
			Version: c.Flags.ClusterAppVersion,
			Values:  c.Flags.Secrets.ClusterValues,
		},
		DefaultApps: cmc.App{
			Name:    c.Flags.DefaultAppsName,
			AppName: fmt.Sprintf("%s-default-apps", c.Cluster),
			Catalog: c.Flags.DefaultAppsCatalog,
			Version: c.Flags.DefaultAppsVersion,
		},
		MCAppsPreventDeletion: c.Flags.MCAppsPreventDeletion,
		PrivateCA:             c.Flags.PrivateCA,
		Provider: cmc.Provider{
			Name: c.Provider,
		},
		TaylorBotToken: c.Flags.TaylorBotToken,
		SSHdeployKey: cmc.DeployKey{
			Passphrase: c.Flags.Secrets.SSHDeployKey.Passphrase,
			Identity:   c.Flags.Secrets.SSHDeployKey.Identity,
			KnownHosts: c.Flags.Secrets.SSHDeployKey.KnownHosts,
		},
		CustomerDeployKey: cmc.DeployKey{
			Passphrase: c.Flags.Secrets.CustomerDeployKey.Passphrase,
			Identity:   c.Flags.Secrets.CustomerDeployKey.Identity,
			KnownHosts: c.Flags.Secrets.CustomerDeployKey.KnownHosts,
		},
		SharedDeployKey: cmc.DeployKey{
			Passphrase: c.Flags.Secrets.SharedDeployKey.Passphrase,
			Identity:   c.Flags.Secrets.SharedDeployKey.Identity,
			KnownHosts: c.Flags.Secrets.SharedDeployKey.KnownHosts,
		},
		DisableDenyAllNetPol: disableDenyAllNetPol(c.Provider),
	}
	if c.Flags.ConfigureContainerRegistries {
		newCMC.ConfigureContainerRegistries = cmc.ConfigureContainerRegistries{
			Enabled: true,
			Values:  c.Flags.Secrets.ContainerRegistryConfiguration,
		}
	}
	if c.Flags.CertManagerDNSChallenge {
		newCMC.CertManagerDNSChallenge = cmc.CertManagerDNSChallenge{
			Enabled:         true,
			Region:          c.Flags.Secrets.CertManagerRoute53Region,
			Role:            c.Flags.Secrets.CertManagerRoute53Role,
			AccessKeyID:     c.Flags.Secrets.CertManagerRoute53AccessKeyID,
			SecretAccessKey: c.Flags.Secrets.CertManagerRoute53SecretAccessKey,
		}
	}
	if c.Flags.MCProxyEnabled {
		newCMC.MCProxy = cmc.MCProxy{
			Enabled:  true,
			Hostname: strings.Split(strings.Split(c.Flags.MCHTTPSProxy, "/")[2], ":")[0],
			Port:     strings.Split(strings.Split(c.Flags.MCHTTPSProxy, "/")[2], ":")[1],
		}
	}
	if c.Flags.MCCustomCoreDNSConfig != "" {
		newCMC.CustomCoreDNS = cmc.CustomCoreDNS{
			Enabled: true,
			Values:  c.Flags.MCCustomCoreDNSConfig,
		}
	}

	if c.Provider == key.ProviderVsphere {
		newCMC.Provider.CAPV = cmc.CAPV{
			CloudConfig: c.Flags.Secrets.VSphereCredentials,
		}
	} else if c.Provider == key.ProviderVCD {
		newCMC.Provider.CAPVCD = cmc.CAPVCD{
			RefreshToken: c.Flags.Secrets.CloudDirectorRefreshToken,
		}
	} else if c.Provider == key.ProviderAzure {
		newCMC.Provider.CAPZ = cmc.CAPZ{
			ClientID:     c.Flags.Secrets.Azure.ClientID,
			TenantID:     c.Flags.Secrets.Azure.TenantID,
			ClientSecret: c.Flags.Secrets.Azure.ClientSecret,
			UAClientID:   c.Flags.Secrets.Azure.UAClientID,
			UATenantID:   c.Flags.Secrets.Azure.UATenantID,
			UAResourceID: c.Flags.Secrets.Azure.UAResourceID,
		}
	}
	if err := newCMC.SetDefaultAppValues(); err != nil {
		return nil, fmt.Errorf("failed to set default app values.\n%w", err)
	}
	return newCMC, nil
}

func overrideCMCWithFlags(currentCMC *cmc.CMC, c Config) (*cmc.CMC, error) {
	newCMC, err := getCMC(c)
	if err != nil {
		return nil, fmt.Errorf("failed to get cmc from flags.\n%w", err)
	}
	return currentCMC.Override(newCMC), nil
}

func disableDenyAllNetPol(provider string) bool {
	return provider != key.ProviderAWS &&
		provider != key.ProviderVsphere &&
		provider != key.ProviderVCD
}

func (c *Config) PullTemplate() (map[string]string, error) {
	githubRepository := github.Repository{
		Github:       c.Github,
		Name:         key.RepositoryMCBootstrap,
		Organization: key.OrganizationGiantSwarm,
		Branch:       key.CMCMainBranch,
	}
	log.Debug().Msg(fmt.Sprintf("pulling cmc entry template from %s repository", key.RepositoryMCBootstrap))
	err := githubRepository.CheckBranch(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to check repository %s branch %s.\n%w", key.RepositoryMCBootstrap, key.CMCMainBranch, err)
	}
	template, err := githubRepository.GetDirectory(context.Background(), key.CMCEntryTemplatePath)
	if err != nil {
		return nil, fmt.Errorf("unable to get template directory from repository %s.\n%w", key.RepositoryMCBootstrap, err)
	}
	//copy each value from the template into the cmc map, Replace the cmc entry template path with the cmc path
	cmcMap := make(map[string]string)
	for k, v := range template {
		cmcMap[key.GetCMCPath(c.Cluster)+k[len(key.CMCEntryTemplatePath):]] = v
	}
	return template, nil
}

func (c *Config) PullTemplateFile(path string) (string, error) {
	githubRepository := github.Repository{
		Github:       c.Github,
		Name:         key.RepositoryMCBootstrap,
		Organization: key.OrganizationGiantSwarm,
		Branch:       key.CMCMainBranch,
	}
	log.Debug().Msg(fmt.Sprintf("pulling cmc entry template file %s from %s repository", path, key.RepositoryMCBootstrap))
	err := githubRepository.CheckBranch(context.Background())
	if err != nil {
		return "", fmt.Errorf("unable to check repository %s branch %s.\n%w", key.RepositoryMCBootstrap, key.CMCMainBranch, err)
	}
	template, err := githubRepository.GetFile(context.Background(), fmt.Sprintf("%s/%s", key.CMCEntryTemplatePath, path))
	if err != nil {
		return "", fmt.Errorf("unable to get file %s from repository %s.\n%w", path, key.RepositoryMCBootstrap, err)
	}
	return template, nil
}
