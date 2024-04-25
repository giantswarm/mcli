package pullcmc

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/giantswarm/mcli/pkg/github"
	"github.com/giantswarm/mcli/pkg/key"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc"
	"github.com/giantswarm/mcli/pkg/sops"
)

type Config struct {
	Cluster        string
	Github         *github.Github
	CMCRepository  string
	CMCBranch      string
	DisplaySecrets bool
}

func (c *Config) Run(ctx context.Context) (*cmc.CMC, error) {
	log.Debug().Msgf("pulling CMC %s", c.Cluster)

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
	sopsfile, err := cmcRepository.GetFile(ctx, cmc.SopsFile)
	if err != nil {
		return nil, err
	}
	data[cmc.SopsFile] = sopsfile

	result, err := cmc.GetCMCFromMap(data, c.Cluster)
	if err != nil {
		return nil, err
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

	if c.CMCRepository == "" {
		return fmt.Errorf("cmc repository is required\n%w", ErrInvalidFlag)
	}
	return nil
}
