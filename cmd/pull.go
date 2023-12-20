/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/giantswarm/mcli/cmd/pull"
	pullinstallations "github.com/giantswarm/mcli/cmd/pull/installations"
	"github.com/giantswarm/mcli/pkg/github"
	"github.com/giantswarm/mcli/pkg/key"
	"github.com/spf13/cobra"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pulls the current configuration of the Management Cluster",
	Long: `Pulls the current configuration of a Management Cluster from all
relevant git repositories. For example:

mcli pull --cluster=gigmac`,
	PreRun: toggleVerbose,
	RunE: func(cmd *cobra.Command, args []string) error {
		defaultPull()
		err := validate(cmd, args)
		if err != nil {
			return err
		}
		ctx := context.Background()
		c := pull.Config{
			Cluster:             cluster,
			GithubToken:         githubToken,
			InstallationsBranch: installationsBranch,
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
	PreRun: toggleVerbose,
	RunE: func(cmd *cobra.Command, args []string) error {
		defaultPull()
		err := validate(cmd, args)
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

func init() {
	rootCmd.AddCommand(pullCmd)
	pullCmd.AddCommand(pullInstallationsCmd)
	pullCmd.Flags().StringArrayVarP(&skip, flagSkip, "s", []string{}, fmt.Sprintf("List of repositories to skip. (default: none) Valid values: %s", key.GetValidRepositories()))
}

func defaultPull() {
	if installationsBranch == "" {
		installationsBranch = key.InstallationsMainBranch
	}
}
