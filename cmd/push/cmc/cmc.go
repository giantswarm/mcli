package pushcmc

import (
	"context"
	"fmt"

	"github.com/giantswarm/mcli/pkg/github"
	"github.com/giantswarm/mcli/pkg/key"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc/certmanager"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc/deploykey"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc/mcproxy"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc/provider/capz"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc/taylorbot"
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
	SecretFolder                 string
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
	MCCustomCoreDNSConfig        bool
	MCProxyEnabled               bool
	MCHTTPSProxy                 string
}

type Values struct {
	AgePubKey         string
	TaylorBot         taylorbot.Config
	Deploykey         deploykey.Config
	CertManager       certmanager.Config
	MCProxy           mcproxy.Config
	ClusterAppValues  string
	DefaultAppsValues string
	RegistryValues    string
	CAPV              string
	CAPVCD            string
	CAPZ              capz.Config
}

func (c *Config) Run(ctx context.Context) (*cmc.CMC, error) {
	err := c.Validate()
	if err != nil {
		return nil, err
	}
	return c.PushCMC(ctx)
}

func (c *Config) PushCMC(ctx context.Context) (*cmc.CMC, error) {
	if err := c.Branch(ctx); err != nil {
		return nil, err
	}
	// pulling current installations
	installations, err := c.Pull(ctx)
	if err != nil {
		// if there is no current installations, create a new one
		// check if the error is a github.ErrNotFound
		if github.IsNotFound(err) {
			log.Debug().Msg(fmt.Sprintf("no current installations %s found, creating a new one", c.Cluster))
			return c.Create(ctx)
		} else {
			return nil, fmt.Errorf("failed to pull installations.\n%w", err)
		}
	}
	return c.Update(ctx, installations)
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
			err = cmcRepository.CreateBranch(ctx, key.InstallationsMainBranch)
			if err != nil {
				return fmt.Errorf("failed to create %s branch %s.\n%w", c.CMCRepository, c.CMCBranch, err)
			}
		} else {
			return fmt.Errorf("failed to check %s branch %s.\n%w", c.CMCRepository, c.CMCBranch, err)
		}
	}
	return nil
}

func (c *Config) Pull(ctx context.Context) (*cmc.CMC, error) {
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
	return cmc.GetCMCFromMap(data), nil
}

func (c *Config) Create(ctx context.Context) (*cmc.CMC, error) {

	log.Debug().Msg(fmt.Sprintf("creating new %s entry for %s", c.CMCRepository, c.Cluster))
	var desiredCMC *cmc.CMC
	{
		if c.Input == nil {
			v, err := c.GetValues()
			if err != nil {
				return nil, fmt.Errorf("failed to get values.\n%w", err)
			}
			desiredCMC, err = getNewCMCFromFlags(*c, v)
			if err != nil {
				return nil, fmt.Errorf("failed to get new %s object from flags.\n%w", c.CMCRepository, err)
			}
		} else {
			desiredCMC = c.Input
		}
	}
	return c.Push(ctx, desiredCMC)
}

func (c *Config) Update(ctx context.Context, currentCMC *cmc.CMC) (*cmc.CMC, error) {
	log.Debug().Msg(fmt.Sprintf("updating %s entry for %s", c.CMCRepository, c.Cluster))
	var desiredCMC *cmc.CMC
	{
		if c.Input == nil {
			v, err := c.GetValues()
			if err != nil {
				return nil, fmt.Errorf("failed to get values.\n%w", err)
			}
			desiredCMC = overrideCMCWithFlags(currentCMC, *c, v)
		} else {
			desiredCMC = currentCMC.Override(c.Input)
		}
	}
	if currentCMC.Equals(desiredCMC) {
		log.Debug().Msg(fmt.Sprintf("%s entry for %s is up to date", c.CMCRepository, c.Cluster))
		return desiredCMC, nil
	}
	return c.Push(ctx, desiredCMC)
}

func (c *Config) Push(ctx context.Context, desiredCMC *cmc.CMC) (*cmc.CMC, error) {
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

	return desiredCMC, nil
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

func (c *Config) GetValues() (Values, error) {
	v := Values{}
	// TODO get values
	return v, nil
}

func getNewCMCFromFlags(c Config, v Values) (*cmc.CMC, error) {
	// Ensure that all the needed flags are set
	if c.Flags.ClusterAppName == "" ||
		c.Flags.ClusterAppCatalog == "" ||
		c.Flags.ClusterAppVersion == "" ||
		c.Flags.DefaultAppsName == "" ||
		c.Flags.DefaultAppsCatalog == "" ||
		c.Flags.DefaultAppsVersion == "" ||
		c.Provider == "" ||
		c.Cluster == "" {
		return nil, fmt.Errorf("not all required flags are set\n%w", ErrInvalidFlag)
	}

	// Ensure that all the needed values are set
	if v.AgePubKey == "" ||
		v.TaylorBot.User == "" ||
		v.TaylorBot.Token == "" ||
		v.Deploykey.Key == "" ||
		v.Deploykey.Identity == "" ||
		v.Deploykey.KnownHosts == "" ||
		v.ClusterAppValues == "" ||
		v.DefaultAppsValues == "" {
		return nil, fmt.Errorf("not all required values are set\n%w", ErrInvalidFlag)
	}
	return getCMC(c, v), nil
}

func getCMC(c Config, v Values) *cmc.CMC {
	newCMC := &cmc.CMC{
		AgePubKey:        v.AgePubKey,
		Cluster:          c.Cluster,
		ClusterNamespace: c.Flags.ClusterNamespace,
		ClusterApp: cmc.App{
			Name:    c.Flags.ClusterAppName,
			Catalog: c.Flags.ClusterAppCatalog,
			Version: c.Flags.ClusterAppVersion,
			Values:  v.ClusterAppValues,
		},
		DefaultApps: cmc.App{
			Name:    c.Flags.DefaultAppsName,
			Catalog: c.Flags.DefaultAppsCatalog,
			Version: c.Flags.DefaultAppsVersion,
			Values:  v.DefaultAppsValues,
		},
		MCAppsPreventDeletion: c.Flags.MCAppsPreventDeletion,
		PrivateCA:             c.Flags.PrivateCA,
		Provider: cmc.Provider{
			Name: c.Provider,
		},
		TaylorBot: cmc.TaylorBot{
			User:  v.TaylorBot.User,
			Token: v.TaylorBot.Token,
		},
		DeployKey: cmc.DeployKey{
			Key:        v.Deploykey.Key,
			Identity:   v.Deploykey.Identity,
			KnownHosts: v.Deploykey.KnownHosts,
		},
		CustomCoreDNS:        c.Flags.MCCustomCoreDNSConfig,
		DisableDenyAllNetPol: disableDenyAllNetPol(c.Provider),
	}
	if c.Flags.ConfigureContainerRegistries {
		newCMC.ConfigureContainerRegistries = cmc.ConfigureContainerRegistries{
			Enabled: true,
			Values:  v.RegistryValues,
		}
	}
	if c.Flags.CertManagerDNSChallenge {
		newCMC.CertManagerDNSChallenge = cmc.CertManagerDNSChallenge{
			Enabled:         true,
			Region:          v.CertManager.Region,
			Role:            v.CertManager.Role,
			AccessKeyID:     v.CertManager.AccessKeyID,
			SecretAccessKey: v.CertManager.SecretAccessKey,
		}
	}
	if c.Flags.MCProxyEnabled {
		newCMC.MCProxy = cmc.MCProxy{
			Enabled:  true,
			HostName: v.MCProxy.HostName,
			Port:     v.MCProxy.Port,
		}
	}
	if c.Provider == key.ProviderVsphere {
		newCMC.Provider.CAPV = cmc.CAPV{
			CloudConfig: v.CAPV,
		}
	} else if c.Provider == key.ProviderVCD {
		newCMC.Provider.CAPVCD = cmc.CAPVCD{
			CloudConfig: v.CAPVCD,
		}
	} else if c.Provider == key.ProviderAzure {
		newCMC.Provider.CAPZ = cmc.CAPZ{
			IdentityUA:       v.CAPZ.IdentityUA,
			IdentitySP:       v.CAPZ.IdentitySP,
			IdentityStaticSP: v.CAPZ.IdentityStaticSP,
		}
	}
	return newCMC
}

func overrideCMCWithFlags(currentCMC *cmc.CMC, c Config, v Values) *cmc.CMC {
	newCMC := getCMC(c, v)
	return currentCMC.Override(newCMC)
}

func disableDenyAllNetPol(provider string) bool {
	return provider != key.ProviderAWS &&
		provider != key.ProviderVsphere &&
		provider != key.ProviderVCD
}
