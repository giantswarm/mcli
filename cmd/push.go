/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/giantswarm/mcli/cmd/push"
	pushinstallations "github.com/giantswarm/mcli/cmd/push/installations"
	"github.com/giantswarm/mcli/pkg/github"
	"github.com/giantswarm/mcli/pkg/managementcluster/installations"
	"github.com/spf13/cobra"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Pushes configuration of a Management Cluster",
	Long: `Pushes configuration of a Management Cluster to all
relevant git repositories. For example:

mcli push --cluster=gigmac --input=cluster.yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		defaultPush()
		err := validateRoot(cmd, args)
		if err != nil {
			return err
		}
		err = validatePush(cmd, args)
		if err != nil {
			return err
		}
		ctx := context.Background()
		c := push.Config{
			Cluster:             cluster,
			GithubToken:         githubToken,
			InstallationsBranch: installationsBranch,
			Skip:                skip,
			Input:               input,
			InstallationsFlags: pushinstallations.InstallationsFlags{
				BaseDomain:    baseDomain,
				CMCRepository: cmcRepository,
				Team:          team,
				Provider:      provider,
				Customer:      customer,
				AWS: pushinstallations.AWSFlags{
					Region:                 awsRegion,
					InstallationAWSAccount: awsAccountID,
				},
			},
		}
		err = push.Run(c, ctx)
		if err != nil {
			return err
		}
		return nil
	},
}

// pushInstallationsCmd represents the push installations command
var pushInstallationsCmd = &cobra.Command{
	Use:   "installations",
	Short: "Pushes configuration of a Management Cluster installations repository entry",
	Long: `Pushes configuration of a Management Cluster installations repository entry to
installations repository. For example:

mcli push installations --cluster=gigmac --input=cluster.yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		defaultPush()
		err := validateRoot(cmd, args)
		if err != nil {
			return err
		}
		err = validatePush(cmd, args)
		if err != nil {
			return err
		}
		ctx := context.Background()
		client := github.New(github.Config{
			Token: githubToken,
		})
		i := pushinstallations.Config{
			Cluster:             cluster,
			Github:              client,
			InstallationsBranch: installationsBranch,
			Flags: pushinstallations.InstallationsFlags{
				BaseDomain:    baseDomain,
				CMCRepository: cmcRepository,
				Team:          team,
				Provider:      provider,
				Customer:      customer,
				AWS: pushinstallations.AWSFlags{
					Region:                 awsRegion,
					InstallationAWSAccount: awsAccountID,
				},
			},
		}
		if input != "" {
			i.Input, err = installations.GetInstallationsFromFile(input)
			if err != nil {
				return fmt.Errorf("failed to get new installations object from input file.\n%w", err)
			}
		}
		installations, err := i.Run(ctx)
		if err != nil {
			return fmt.Errorf("failed to push installations.\n%w", err)
		}
		return installations.Print()
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.AddCommand(pushInstallationsCmd)
	addFlagsPush()
}
