package pullinstallations

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
	Github              *github.Github
	InstallationsBranch string
}

func (c *Config) Run(ctx context.Context) (*installations.Installations, error) {
	log.Debug().Msg(fmt.Sprintf("pulling installations %s", c.Cluster))

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
