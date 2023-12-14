package pull

import (
	"context"
	"fmt"

	pullinstallations "github.com/giantswarm/mcli/cmd/pull/installations"
	"github.com/giantswarm/mcli/pkg/github"
	"github.com/giantswarm/mcli/pkg/key"
	"github.com/giantswarm/mcli/pkg/managementcluster"

	"github.com/rs/zerolog/log"
)

type Config struct {
	Cluster             string
	GithubToken         string
	InstallationsBranch string
	Skip                []string
}

func Run(c Config, ctx context.Context) error {
	mc, err := c.Pull(ctx)
	if err != nil {
		return fmt.Errorf("failed to pull management cluster configuration.\n%w", err)
	}
	return mc.Print()
}

func (c *Config) Pull(ctx context.Context) (*managementcluster.ManagementCluster, error) {
	var mc managementcluster.ManagementCluster
	log.Debug().Msg(fmt.Sprintf("pulling management cluster %s", c.Cluster))

	client := github.New(github.Config{
		Token: c.GithubToken,
	})

	if !key.Skip(key.RepositoryInstallations, c.Skip) {
		i := pullinstallations.Config{
			Cluster:             c.Cluster,
			Github:              client,
			InstallationsBranch: c.InstallationsBranch,
		}
		installations, err := i.Run(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to pull installations.\n%w", err)
		}
		mc.Installations = *installations
	}
	return &mc, nil
}
