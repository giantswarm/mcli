/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/giantswarm/mcli/cmd/pull"
	pullcmc "github.com/giantswarm/mcli/cmd/pull/cmc"
	pullinstallations "github.com/giantswarm/mcli/cmd/pull/installations"
	"github.com/giantswarm/mcli/pkg/github"
	"github.com/spf13/cobra"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pulls the current configuration of the Management Cluster",
	Long: `Pulls the current configuration of a Management Cluster from all
relevant git repositories. For example:

mcli pull --cluster=gigmac`,
	RunE: func(cmd *cobra.Command, args []string) error {
		defaultPull()
		err := validateRoot(cmd, args)
		if err != nil {
			return err
		}
		ctx := context.Background()
		c := pull.Config{
			Cluster:             cluster,
			GithubToken:         githubToken,
			InstallationsBranch: installationsBranch,
			CMCBranch:           cmcBranch,
			CMCRepository:       cmcRepository,
			Skip:                skip,
		}
		err = pull.Run(c, ctx)
		if err != nil {
			return err
		}
		return nil
	},
}

// pullInstallationsCmd represents the pull installations command
var pullInstallationsCmd = &cobra.Command{
	Use:   "installations",
	Short: "Pulls the current configuration of the Management Clusters installations repository entry",
	Long: `Pulls the current configuration of a Management Clusters installations repository entry from
installations repository. For example:

mcli pull installations --cluster=gigmac`,
	RunE: func(cmd *cobra.Command, args []string) error {
		defaultPull()
		err := validateRoot(cmd, args)
		if err != nil {
			return err
		}
		ctx := context.Background()
		client := github.New(github.Config{
			Token: githubToken,
		})
		i := pullinstallations.Config{
			Cluster:             cluster,
			Github:              client,
			InstallationsBranch: installationsBranch,
		}
		installations, err := i.Run(ctx)
		if err != nil {
			return fmt.Errorf("failed to pull installations.\n%w", err)
		}
		return installations.Print()
	},
}

var pullCMCCmd = &cobra.Command{
	Use:   "cmc",
	Short: "Pulls the current configuration of the Management Clusters CMC repository entry",
	Long: `Pulls the current configuration of a Management Clusters CMC repository entry from
installations repository. For example:

mcli pull cmc --cluster=gigmac`,
	RunE: func(cmd *cobra.Command, args []string) error {
		defaultPull()
		err := validateRoot(cmd, args)
		if err != nil {
			return err
		}
		ctx := context.Background()
		client := github.New(github.Config{
			Token: githubToken,
		})
		c := pullcmc.Config{
			Cluster:       cluster,
			Github:        client,
			CMCRepository: cmcRepository,
			CMCBranch:     cmcBranch,
		}
		cmc, err := c.Run(ctx)
		if err != nil {
			return fmt.Errorf("failed to pull CMC.\n%w", err)
		}
		return cmc.Print()
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
	pullCmd.AddCommand(pullInstallationsCmd)
	pullCmd.AddCommand(pullCMCCmd)
	addFlagsPull()
}
