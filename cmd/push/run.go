package push

import (
	"context"
	"fmt"

	pushinstallations "github.com/giantswarm/mcli/cmd/push/installations"
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
	Input               string
	InstallationsFlags  pushinstallations.InstallationsFlags
}

func Run(c Config, ctx context.Context) error {
	mc, err := c.Push(ctx)
	if err != nil {
		return fmt.Errorf("failed to push management cluster configuration.\n%w", err)
	}
	return mc.Print()
}

func (c *Config) Push(ctx context.Context) (*managementcluster.ManagementCluster, error) {
	var err error

	log.Debug().Msg(fmt.Sprintf("pushing management cluster %s", c.Cluster))
	mc := &managementcluster.ManagementCluster{}

	if c.Input != "" {
		mc, err = managementcluster.GetManagementClusterFromFile(c.Input)
		if err != nil {
			return nil, fmt.Errorf("failed to get new management cluster object from input file.\n%w", err)
		}
	}

	client := github.New(github.Config{
		Token: c.GithubToken,
	})

	if !key.Skip(key.RepositoryInstallations, c.Skip) {

		i := pushinstallations.Config{
			Cluster:             c.Cluster,
			Github:              client,
			InstallationsBranch: c.InstallationsBranch,
			Flags:               c.InstallationsFlags,
		}
		if c.Input != "" {
			i.Input = &mc.Installations
		}
		installations, err := i.Run(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to push installations.\n%w", err)
		}
		mc.Installations = *installations
	}
	return mc, nil
}
