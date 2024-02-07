package pullcmc

import (
	"context"

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
	return &cmc.CMC{
		Contents: data,
	}, nil
}
