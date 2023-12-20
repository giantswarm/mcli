/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"

	"github.com/giantswarm/mcli/cmd/pull"
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

func init() {
	rootCmd.AddCommand(pullCmd)
}

func defaultPull() {
	if installationsBranch == "" {
		installationsBranch = key.InstallationsMainBranch
	}
}
