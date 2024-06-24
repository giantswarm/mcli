package pushcmc

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/giantswarm/mcli/pkg/github"
	"github.com/giantswarm/mcli/pkg/key"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc/kustomization"
	"github.com/giantswarm/mcli/pkg/sops"
)

type Config struct {
	Cluster        string
	Github         *github.Github
	BaseDomain     string
	CMCRepository  string
	CMCBranch      string
	Provider       string
	Input          *cmc.CMC
	Flags          CMCFlags
	DisplaySecrets bool
}

type CMCFlags struct {
	Secrets                      SecretFlags
	SecretFolder                 string
	AgePubKey                    string
	TaylorBotToken               string
	MCAppsPreventDeletion        bool
	ClusterAppName               string
	ClusterAppCatalog            string
	ClusterAppVersion            string
	ClusterNamespace             string
	ClusterIntegratesDefaultApps bool
	ConfigureContainerRegistries bool
	DefaultAppsName              string
	DefaultAppsCatalog           string
	DefaultAppsVersion           string
	PrivateCA                    bool
	PrivateMC                    bool
	CertManagerDNSChallenge      bool
	MCCustomCoreDNSConfig        string
	MCProxyEnabled               bool
	MCHTTPSProxy                 string
	CatalogRegistryValues        string
	MCBBranchSource              string
	ConfigBranch                 string
	MCAppCollectionBranch        string
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
	UAClientID     string
	UATenantID     string
	UAResourceID   string
	ClientID       string
	TenantID       string
	ClientSecret   string
	SubscriptionID string
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
	cmcRepository := github.Repository{
		Github:       c.Github,
		Name:         c.CMCRepository,
		Organization: key.OrganizationGiantSwarm,
		Branch:       c.CMCBranch,
	}

	if err := c.Branch(ctx, cmcRepository); err != nil {
		return nil, err
	}
	// pulling current cmc
	cmc, err := c.Pull(ctx, cmcRepository)
	if err != nil {
		// if there is no current cmc, create a new one
		// check if the error is a github.ErrNotFound
		if github.IsNotFound(err) {
			sopsFile, err := c.pullSopsFile(ctx, cmcRepository)
			if err != nil {
				return nil, err
			}
			log.Debug().Msg(fmt.Sprintf("no current %s entry for %s found, creating a new one", c.CMCRepository, c.Cluster))
			return c.Create(ctx, sopsFile)
		} else {
			return nil, fmt.Errorf("failed to pull %s entry for %s.\n%w", c.CMCRepository, c.Cluster, err)
		}
	}
	return c.Update(ctx, cmc)
}

func (c *Config) Branch(ctx context.Context, cmcRepository github.Repository) error {
	log.Debug().Msg(fmt.Sprintf("getting %s branch %s", c.CMCRepository, c.CMCBranch))

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

func (c *Config) Pull(ctx context.Context, cmcRepository github.Repository) (map[string]string, error) {
	log.Debug().Msgf("pulling current %s entry for %s", c.CMCRepository, c.Cluster)

	sopsfile, err := c.pullSopsFile(ctx, cmcRepository)
	if err != nil {
		return nil, err
	}

	data, err := cmcRepository.GetDirectory(ctx, key.GetCMCPath(c.Cluster))
	if err != nil {
		return nil, err
	}
	data[cmc.SopsFile] = sopsfile

	return data, nil
}

func (c *Config) pullSopsFile(ctx context.Context, cmcRepository github.Repository) (string, error) {
	if err := cmcRepository.Check(ctx); err != nil {
		return "", err
	}
	sopsfile, err := cmcRepository.GetFile(ctx, cmc.SopsFile)
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("no %s file found in %s repository.\n%s", cmc.SopsFile, c.CMCRepository, err))
		return "", nil
	}
	return sopsfile, nil
}

func (c *Config) Create(ctx context.Context, sopsFile string) (*cmc.CMC, error) {
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
	template[cmc.SopsFile] = sopsFile
	create, err := desiredCMC.GetMap(template)
	if err != nil {
		return nil, fmt.Errorf("failed to get cmc map.\n%w", err)
	}

	message := fmt.Sprintf("Create configuration of management cluster %s", c.Cluster)

	return c.Push(ctx, create, message)
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
		if !c.DisplaySecrets {
			desiredCMC.RedactSecrets()
		}
		return desiredCMC, nil
	}
	// if deny all netpol is being enabled, we need to get the file from the template
	if currentCMC.DisableDenyAllNetPol && !desiredCMC.DisableDenyAllNetPol {
		template, err := c.PullTemplateFile(kustomization.DenyNetPolFile)
		if err != nil {
			return nil, fmt.Errorf("failed to pull template file.\n%w", err)
		}
		currentCMCmap[fmt.Sprintf("%s/%s", key.GetCMCPath(c.Cluster), kustomization.DenyNetPolFile)] = template
	}
	update, err := desiredCMC.GetMap(currentCMCmap)
	if err != nil {
		return nil, fmt.Errorf("failed to get cmc map.\n%w", err)
	}

	update, err = cmc.MarkUnchangedSecretsInMap(currentCMC, desiredCMC, update)
	if err != nil {
		return nil, fmt.Errorf("failed to mark unchanged secrets.\n%w", err)
	}
	message := fmt.Sprintf("Update configuration of management cluster %s", c.Cluster)
	return c.Push(ctx, update, message)
}

func (c *Config) Push(ctx context.Context, desiredCMC map[string]string, message string) (*cmc.CMC, error) {
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

	err = cmcRepository.CreateDirectory(ctx, desiredCMC, message)
	if err != nil {
		return nil, fmt.Errorf("failed to create directory %s.\n%w", key.GetCMCPath(c.Cluster), err)
	}

	result, err := cmc.GetCMCFromMap(desiredCMC, c.Cluster)
	if err != nil {
		return nil, fmt.Errorf("failed to get cmc from map.\n%w", err)
	}
	if !c.DisplaySecrets {
		result.RedactSecrets()
	}
	return result, nil
}

func (c *Config) Validate() error {
	// check if environment variable age key is set
	if val, present := os.LookupEnv(sops.EnvAgeKey); !present || val == "" {
		return fmt.Errorf("environment variable %s is not set\n%w", sops.EnvAgeKey, ErrInvalidFlag)
	}
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
	err := c.EnsureFlagsAreSet()
	if err != nil {
		return nil, fmt.Errorf("failed to ensure flags are set.\n%w", err)
	}
	return getCMC(c)
}

func (c *Config) EnsureFlagsAreSet() error {
	// Ensure that all the needed flags are set
	if c.Flags.ClusterAppName == "" {
		return fmt.Errorf("cluster app name is required\n%w", ErrInvalidFlag)
	}
	if c.BaseDomain == "" {
		return fmt.Errorf("base domain is required\n%w", ErrInvalidFlag)
	}
	if c.Flags.CatalogRegistryValues == "" {
		return fmt.Errorf("catalog registry values are required\n%w", ErrInvalidFlag)
	}
	if c.Flags.ConfigBranch == "" {
		return fmt.Errorf("config branch is required\n%w", ErrInvalidFlag)
	}
	if c.Flags.MCBBranchSource == "" {
		return fmt.Errorf("mcb branch source is required\n%w", ErrInvalidFlag)
	}
	if c.Flags.MCAppCollectionBranch == "" {
		return fmt.Errorf("mc app collection branch is required\n%w", ErrInvalidFlag)
	}
	if c.Flags.ClusterAppCatalog == "" {
		return fmt.Errorf("cluster app catalog is required\n%w", ErrInvalidFlag)
	}
	if c.Flags.ClusterAppVersion == "" {
		return fmt.Errorf("cluster app version is required\n%w", ErrInvalidFlag)
	}
	if c.Flags.ClusterNamespace == "" {
		return fmt.Errorf("cluster namespace is required\n%w", ErrInvalidFlag)
	}
	if c.Provider == "" {
		return fmt.Errorf("provider is required\n%w", ErrInvalidFlag)
	}
	if c.Cluster == "" {
		return fmt.Errorf("cluster is required\n%w", ErrInvalidFlag)
	}
	if c.Flags.AgePubKey == "" {
		return fmt.Errorf("age public key is required\n%w", ErrInvalidFlag)
	}
	if c.Flags.TaylorBotToken == "" {
		return fmt.Errorf("taylor bot token is required\n%w", ErrInvalidFlag)
	}

	// Ensure that all the needed secret values are set
	if c.Flags.Secrets.SSHDeployKey.Identity == "" {
		return fmt.Errorf("ssh deploy key identity is required\n%w", ErrInvalidFlag)
	}
	if c.Flags.Secrets.SSHDeployKey.KnownHosts == "" {
		return fmt.Errorf("ssh deploy key known hosts is required\n%w", ErrInvalidFlag)
	}
	if c.Flags.Secrets.SSHDeployKey.Passphrase == "" {
		return fmt.Errorf("ssh deploy key passphrase is required\n%w", ErrInvalidFlag)
	}
	if c.Flags.Secrets.CustomerDeployKey.Identity == "" {
		return fmt.Errorf("customer deploy key identity is required\n%w", ErrInvalidFlag)
	}
	if c.Flags.Secrets.CustomerDeployKey.KnownHosts == "" {
		return fmt.Errorf("customer deploy key known hosts is required\n%w", ErrInvalidFlag)
	}
	if c.Flags.Secrets.CustomerDeployKey.Passphrase == "" {
		return fmt.Errorf("customer deploy key passphrase is required\n%w", ErrInvalidFlag)
	}
	if c.Flags.Secrets.SharedDeployKey.Identity == "" {
		return fmt.Errorf("shared deploy key identity is required\n%w", ErrInvalidFlag)
	}
	if c.Flags.Secrets.SharedDeployKey.KnownHosts == "" {
		return fmt.Errorf("shared deploy key known hosts is required\n%w", ErrInvalidFlag)
	}
	if c.Flags.Secrets.SharedDeployKey.Passphrase == "" {
		return fmt.Errorf("shared deploy key passphrase is required\n%w", ErrInvalidFlag)
	}
	if c.Flags.Secrets.ClusterValues == "" {
		return fmt.Errorf("cluster values are required\n%w", ErrInvalidFlag)
	}
	// Ensure that needed values for enabled features are set
	if !c.Flags.ClusterIntegratesDefaultApps {
		if c.Flags.DefaultAppsName == "" {
			return fmt.Errorf("default apps name is required\n%w", ErrInvalidFlag)
		}
		if c.Flags.DefaultAppsCatalog == "" {
			return fmt.Errorf("default apps catalog is required\n%w", ErrInvalidFlag)
		}
		if c.Flags.DefaultAppsVersion == "" {
			return fmt.Errorf("default apps version is required\n%w", ErrInvalidFlag)
		}
	}

	if c.Flags.ConfigureContainerRegistries {
		if c.Flags.Secrets.ContainerRegistryConfiguration == "" {
			return fmt.Errorf("container registry configuration is required\n%w", ErrInvalidFlag)
		}
	}
	if c.Flags.CertManagerDNSChallenge {
		if c.Flags.Secrets.CertManagerRoute53Region == "" {
			return fmt.Errorf("cert manager route53 region is required\n%w", ErrInvalidFlag)
		}
		if c.Flags.Secrets.CertManagerRoute53Role == "" {
			return fmt.Errorf("cert manager route53 role is required\n%w", ErrInvalidFlag)
		}
		if c.Flags.Secrets.CertManagerRoute53AccessKeyID == "" {
			return fmt.Errorf("cert manager route53 access key id is required\n%w", ErrInvalidFlag)
		}
		if c.Flags.Secrets.CertManagerRoute53SecretAccessKey == "" {
			return fmt.Errorf("cert manager route53 secret access key is required\n%w", ErrInvalidFlag)
		}
	}
	if c.Flags.MCProxyEnabled {
		if c.Flags.MCHTTPSProxy == "" {
			return fmt.Errorf("mc https proxy is required\n%w", ErrInvalidFlag)
		}
	}
	// Ensure that needed values for the provider are set
	if key.IsProviderVsphere(c.Provider) {
		if c.Flags.Secrets.VSphereCredentials == "" {
			return fmt.Errorf("vsphere credentials are required\n%w", ErrInvalidFlag)
		}
	} else if key.IsProviderVCD(c.Provider) {
		if c.Flags.Secrets.CloudDirectorRefreshToken == "" {
			return fmt.Errorf("cloud director refresh token is required\n%w", ErrInvalidFlag)
		}
	} else if key.IsProviderAzure(c.Provider) {
		if c.Flags.Secrets.Azure.ClientID == "" {
			return fmt.Errorf("azure client id is required\n%w", ErrInvalidFlag)
		}
		if c.Flags.Secrets.Azure.TenantID == "" {
			return fmt.Errorf("azure tenant id is required\n%w", ErrInvalidFlag)
		}
		if c.Flags.Secrets.Azure.ClientSecret == "" {
			return fmt.Errorf("azure client secret is required\n%w", ErrInvalidFlag)
		}
		if c.Flags.Secrets.Azure.SubscriptionID == "" {
			return fmt.Errorf("azure subscription id is required\n%w", ErrInvalidFlag)
		}
		if c.Flags.Secrets.Azure.UAClientID == "" {
			return fmt.Errorf("azure user assigned client id is required\n%w", ErrInvalidFlag)
		}
		if c.Flags.Secrets.Azure.UATenantID == "" {
			return fmt.Errorf("azure user assigned tenant id is required\n%w", ErrInvalidFlag)
		}
		if c.Flags.Secrets.Azure.UAResourceID == "" {
			return fmt.Errorf("azure user assigned resource id is required\n%w", ErrInvalidFlag)
		}
	}
	return nil
}

func getCMC(c Config) (*cmc.CMC, error) {
	newCMC := &cmc.CMC{
		AgePubKey:             c.Flags.AgePubKey,
		Cluster:               c.Cluster,
		BaseDomain:            c.BaseDomain,
		CatalogRegistryValues: c.Flags.CatalogRegistryValues,
		GitOps: cmc.GitOps{
			CMCRepository:         c.CMCRepository,
			CMCBranch:             c.CMCBranch,
			MCBBranchSource:       c.Flags.MCBBranchSource,
			ConfigBranch:          c.Flags.ConfigBranch,
			MCAppCollectionBranch: c.Flags.MCAppCollectionBranch,
		},
		ClusterNamespace: c.Flags.ClusterNamespace,
		ClusterApp: cmc.App{
			Name:    c.Flags.ClusterAppName,
			AppName: c.Cluster,
			Catalog: c.Flags.ClusterAppCatalog,
			Version: c.Flags.ClusterAppVersion,
			Values:  c.Flags.Secrets.ClusterValues,
		},
		ClusterIntegratesDefaultApps: c.Flags.ClusterIntegratesDefaultApps,
		MCAppsPreventDeletion:        c.Flags.MCAppsPreventDeletion,
		PrivateCA:                    c.Flags.PrivateCA,
		PrivateMC:                    c.Flags.PrivateMC,
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
	if !c.Flags.ClusterIntegratesDefaultApps {
		newCMC.DefaultApps = cmc.App{
			Name:    c.Flags.DefaultAppsName,
			AppName: fmt.Sprintf("%s-default-apps", c.Cluster),
			Catalog: c.Flags.DefaultAppsCatalog,
			Version: c.Flags.DefaultAppsVersion,
		}
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
		proxy := strings.Split(c.Flags.MCHTTPSProxy, "/")
		if len(proxy) != 3 {
			return nil, fmt.Errorf("invalid mc https proxy format %s. Expected format: http://<hostname>:<port>.\n%w", c.Flags.MCHTTPSProxy, ErrInvalidFlag)
		}
		host := strings.Split(proxy[2], ":")
		if len(host) != 2 {
			return nil, fmt.Errorf("invalid mc https proxy format %s. Expected format: http://<hostname>:<port>.\n%w", c.Flags.MCHTTPSProxy, ErrInvalidFlag)
		}
		newCMC.MCProxy = cmc.MCProxy{
			Enabled:  true,
			Hostname: host[0],
			Port:     host[1],
		}
	}
	if c.Flags.MCCustomCoreDNSConfig != "" {
		newCMC.CustomCoreDNS = cmc.CustomCoreDNS{
			Enabled: true,
			Values:  c.Flags.MCCustomCoreDNSConfig,
		}
	}

	if key.IsProviderVsphere(c.Provider) {
		newCMC.Provider.CAPV = cmc.CAPV{
			CloudConfig: c.Flags.Secrets.VSphereCredentials,
		}
	} else if key.IsProviderVCD(c.Provider) {
		newCMC.Provider.CAPVCD = cmc.CAPVCD{
			RefreshToken: c.Flags.Secrets.CloudDirectorRefreshToken,
		}
	} else if key.IsProviderAzure(c.Provider) {
		newCMC.Provider.CAPZ = cmc.CAPZ{
			ClientID:       c.Flags.Secrets.Azure.ClientID,
			TenantID:       c.Flags.Secrets.Azure.TenantID,
			ClientSecret:   c.Flags.Secrets.Azure.ClientSecret,
			SubscriptionID: c.Flags.Secrets.Azure.SubscriptionID,
			UAClientID:     c.Flags.Secrets.Azure.UAClientID,
			UATenantID:     c.Flags.Secrets.Azure.UATenantID,
			UAResourceID:   c.Flags.Secrets.Azure.UAResourceID,
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
	return !key.IsProviderAWS(provider) &&
		!key.IsProviderVsphere(provider) &&
		!key.IsProviderVCD(provider)
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
	return cmcMap, nil
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
