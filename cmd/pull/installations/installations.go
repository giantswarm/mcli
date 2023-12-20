package pullinstallations

import (
	"context"
	"fmt"

	"github.com/giantswarm/mcli/pkg/github"
	"github.com/giantswarm/mcli/pkg/key"
	"github.com/giantswarm/mcli/pkg/managementcluster/installations"
	"github.com/rs/zerolog/log"
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
		Path:         key.GetInstallationsPath(c.Cluster),
	}
	data, err := installationsRepository.GetFile(ctx)
	if err != nil {
		return nil, err
	}
	installations, err := installations.GetInstallations([]byte(data))
	if err != nil {
		return nil, err
	}
	return installations, nil
}
