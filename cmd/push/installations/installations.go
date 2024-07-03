package pushinstallations

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/giantswarm/mcli/pkg/github"
	"github.com/giantswarm/mcli/pkg/key"
	"github.com/giantswarm/mcli/pkg/managementcluster/installations"
)

type Config struct {
	Cluster             string
	BaseDomain          string
	Provider            string
	CMCRepository       string
	Github              *github.Github
	InstallationsBranch string
	Input               *installations.Installations
	Flags               InstallationsFlags
}

type InstallationsFlags struct {
	CCRRepository string
	Team          string
	AWS           AWSFlags
	Customer      string
}

type AWSFlags struct {
	Region                 string
	InstallationAWSAccount string
}

func (c *Config) Run(ctx context.Context) (*installations.Installations, error) {
	err := c.Validate()
	if err != nil {
		return nil, err
	}

	return c.PushInstallations(ctx)
}

func (c *Config) PushInstallations(ctx context.Context) (*installations.Installations, error) {
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

func (c *Config) Create(ctx context.Context) (*installations.Installations, error) {
	var err error

	log.Debug().Msg(fmt.Sprintf("creating new installations %s", c.Cluster))
	var desiredInstallations *installations.Installations
	{
		if c.Input == nil {
			desiredInstallations, err = getNewInstallationsFromFlags(*c)
			if err != nil {
				return nil, fmt.Errorf("failed to get new installations object from flags.\n%w", err)
			}
		} else {
			desiredInstallations = c.Input
		}
	}
	return c.Push(ctx, desiredInstallations)
}

func (c *Config) Update(ctx context.Context, currentInstallations *installations.Installations) (*installations.Installations, error) {
	log.Debug().Msg(fmt.Sprintf("updating installations %s", c.Cluster))
	var desiredInstallations *installations.Installations
	{
		if c.Input == nil {
			desiredInstallations = overrideInstallationsWithFlags(currentInstallations, *c)
		} else {
			desiredInstallations = currentInstallations.Override(c.Input)
		}
	}
	if currentInstallations.Equals(desiredInstallations) {
		log.Debug().Msg("installations are up to date")
		return desiredInstallations, nil
	}
	return c.Push(ctx, desiredInstallations)
}

func (c *Config) Pull(ctx context.Context) (*installations.Installations, error) {
	log.Debug().Msg(fmt.Sprintf("pulling current installations %s", c.Cluster))

	installationsRepository := github.Repository{
		Github:       c.Github,
		Name:         key.RepositoryInstallations,
		Organization: key.OrganizationGiantSwarm,
		Branch:       c.InstallationsBranch,
	}
	if err := installationsRepository.Check(ctx); err != nil {
		return nil, err
	}
	data, err := installationsRepository.GetFile(ctx, key.GetInstallationsPath(c.Cluster))
	if err != nil {
		return nil, err
	}
	installations, err := installations.GetInstallations([]byte(data))
	if err != nil {
		return nil, err
	}
	return installations, nil
}

func (c *Config) Push(ctx context.Context, i *installations.Installations) (*installations.Installations, error) {
	// check if i is valid
	err := i.Validate()
	if err != nil {
		return nil, err
	}
	if i.Codename != c.Cluster {
		return nil, fmt.Errorf("cluster name %s does not match installations codename %s", c.Cluster, i.Codename)
	}

	log.Debug().Msg(fmt.Sprintf("pushing installations %s", c.Cluster))
	data, err := installations.GetData(i)
	if err != nil {
		return nil, err
	}

	installationsRepository := github.Repository{
		Github:       c.Github,
		Name:         key.RepositoryInstallations,
		Organization: key.OrganizationGiantSwarm,
		Branch:       c.InstallationsBranch,
	}
	err = installationsRepository.CreateFile(ctx, data, key.GetInstallationsPath(c.Cluster))
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (c *Config) Branch(ctx context.Context) error {
	log.Debug().Msg(fmt.Sprintf("getting installations branch %s", c.InstallationsBranch))

	installationsRepository := github.Repository{
		Github:       c.Github,
		Name:         key.RepositoryInstallations,
		Organization: key.OrganizationGiantSwarm,
		Branch:       c.InstallationsBranch,
	}
	err := installationsRepository.CheckBranch(ctx)
	if err != nil {
		// if the branch doesn't exist, create it
		// check if the error is a github.ErrNotFound
		if github.IsNotFound(err) {
			log.Debug().Msg(fmt.Sprintf("installations branch %s not found, creating it", c.InstallationsBranch))
			err = installationsRepository.CreateBranch(ctx, key.InstallationsMainBranch)
			if err != nil {
				return fmt.Errorf("failed to create installations branch %s.\n%w", c.InstallationsBranch, err)
			}
		} else {
			return fmt.Errorf("failed to check installations branch %s.\n%w", c.InstallationsBranch, err)
		}
	}
	return nil
}

func getNewInstallationsFromFlags(flags Config) (*installations.Installations, error) {
	//Ensure that all the needed flags are set
	if flags.BaseDomain == "" {
		return nil, fmt.Errorf("base domain is not set.\n%w", ErrInvalidFlag)
	}
	if flags.CMCRepository == "" {
		return nil, fmt.Errorf("CMC repository is not set.\n%w", ErrInvalidFlag)
	}
	if flags.Flags.CCRRepository == "" {
		return nil, fmt.Errorf("CCR repository is not set.\n%w", ErrInvalidFlag)
	}
	if flags.Flags.Team == "" {
		return nil, fmt.Errorf("team is not set.\n%w", ErrInvalidFlag)
	}
	if flags.Flags.Customer == "" {
		return nil, fmt.Errorf("customer is not set.\n%w", ErrInvalidFlag)
	}
	if flags.Provider == "" {
		return nil, fmt.Errorf("provider is not set.\n%w", ErrInvalidFlag)
	}
	if flags.Cluster == "" {
		return nil, fmt.Errorf("cluster is not set.\n%w", ErrInvalidFlag)
	}
	if key.IsProviderAWS(flags.Provider) {
		if flags.Flags.AWS.Region == "" {
			return nil, fmt.Errorf("AWS region is not set.\n%w", ErrInvalidFlag)
		}
		if flags.Flags.AWS.InstallationAWSAccount == "" {
			return nil, fmt.Errorf("AWS installation account is not set.\n%w", ErrInvalidFlag)
		}
	}

	log.Debug().Msg("getting new installations object from flags")

	c := installations.InstallationsConfig{
		Base:            flags.BaseDomain,
		Codename:        flags.Cluster,
		Customer:        flags.Flags.Customer,
		CmcRepository:   flags.CMCRepository,
		CcrRepository:   flags.Flags.CCRRepository,
		AccountEngineer: flags.Flags.Team,
		Pipeline:        "testing",
		Provider:        fmt.Sprintf("%s-test", flags.Provider),
	}
	if key.IsProviderAWS(flags.Provider) {
		c.AwsRegion = flags.Flags.AWS.Region
		c.AwsHostClusterAccount = flags.Flags.AWS.InstallationAWSAccount
		c.AwsHostClusterAdminRoleArn = fmt.Sprintf("arn:aws:iam::%s:role/GiantSwarmAdmin", flags.Flags.AWS.InstallationAWSAccount)
		c.AwsHostClusterCloudtrailBucket = ""
		c.AwsHostClusterGuardDuty = false
		c.AwsGuestClusterAccount = flags.Flags.AWS.InstallationAWSAccount
		c.AwsGuestClusterCloudtrailBucket = ""
		c.AwsGuestClusterGuardDuty = false
	}

	return installations.NewInstallations(c), nil
}

func overrideInstallationsWithFlags(current *installations.Installations, flags Config) *installations.Installations {
	log.Debug().Msg("overriding installations object with flags")

	c := installations.InstallationsConfig{
		Codename:        flags.Cluster,
		Base:            flags.BaseDomain,
		CmcRepository:   flags.CMCRepository,
		CcrRepository:   flags.Flags.CCRRepository,
		AccountEngineer: flags.Flags.Team,
		Provider:        flags.Provider,
		Customer:        flags.Flags.Customer,
	}
	if key.IsProviderAWS(flags.Provider) {
		c.AwsRegion = flags.Flags.AWS.Region
		if flags.Flags.AWS.InstallationAWSAccount != "" {
			c.AwsHostClusterAccount = flags.Flags.AWS.InstallationAWSAccount
			c.AwsGuestClusterAccount = flags.Flags.AWS.InstallationAWSAccount
			c.AwsHostClusterAdminRoleArn = fmt.Sprintf("arn:aws:iam::%s:role/GiantSwarmAdmin", flags.Flags.AWS.InstallationAWSAccount)
		}
	}
	return current.Override(installations.NewInstallations(c))
}

func (c *Config) Validate() error {
	if c.InstallationsBranch == "main" || c.InstallationsBranch == "master" {
		return fmt.Errorf("cannot push to installations branch %s.\n%w", c.InstallationsBranch, ErrInvalidFlag)
	}
	if c.Input != nil {
		log.Debug().Msg("using input file. Other installations flags will be ignored")
		return nil
	}
	if c.BaseDomain == "" &&
		c.CMCRepository == "" &&
		c.Flags.CCRRepository == "" &&
		c.Flags.Team == "" &&
		c.Provider == "" &&
		c.Flags.AWS.Region == "" &&
		c.Flags.Customer == "" &&
		c.Flags.AWS.InstallationAWSAccount == "" {
		return fmt.Errorf("no input file or flags specified.\n%w", ErrInvalidFlag)
	}

	if c.Provider != "" {
		if !key.IsValidProvider(c.Provider) {
			return fmt.Errorf("invalid provider %s. Valid values: %s:\n%w", c.Provider, key.GetValidProviders(), ErrInvalidFlag)
		}
	}
	// todo: check format of other flags if they are set
	return nil
}
