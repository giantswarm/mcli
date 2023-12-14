/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/giantswarm/mcli/cmd/push"
	pushinstallations "github.com/giantswarm/mcli/cmd/push/installations"
	"github.com/giantswarm/mcli/pkg/key"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagInput         = "input"
	flagBaseDomain    = "base-domain"
	flagCMCRepository = "cmc-repository"
	flagTeam          = "team"
	flagProvider      = "provider"
	flagAWSRegion     = "aws-region"
	flagAWSAccountID  = "aws-account-id"
)

const (
	envBaseDomain    = "BASE_DOMAIN"
	envCMCRepository = "CMC_REPOSITORY"
	envTeam          = "TEAM_NAME"
	envProvider      = "PROVIDER"
	envAWSRegion     = "AWS_REGION"
	envAWSAccountID  = "INSTALLATION_AWS_ACCOUNT"
)

var (
	input         string
	baseDomain    string
	cmcRepository string
	team          string
	provider      string
	awsRegion     string
	awsAccountID  string
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Pushes configuration of a Management Cluster",
	Long: `Pushes configuration of a Management Cluster to all
relevant git repositories. For example:

mcli push --cluster=gigmac --input=cluster.yaml`,
	PreRun: toggleVerbose,
	RunE: func(cmd *cobra.Command, args []string) error {
		defaultPush()
		err := validate(cmd, args)
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

func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.Flags().StringVarP(&input, flagInput, "i", "", "Input configuration file to use. If not specified, configuration is read from other flags.")

	pushCmd.Flags().StringVar(&baseDomain, flagBaseDomain, "", "Base domain to use for the cluster")
	viper.BindEnv(flagBaseDomain, envBaseDomain)
	viper.BindPFlag(flagBaseDomain, pushCmd.Flags().Lookup(flagBaseDomain))

	pushCmd.Flags().StringVar(&cmcRepository, flagCMCRepository, "", "Name of CMC repository")
	viper.BindEnv(flagCMCRepository, envCMCRepository)
	viper.BindPFlag(flagCMCRepository, pushCmd.Flags().Lookup(flagCMCRepository))

	pushCmd.Flags().StringVar(&team, flagTeam, "", "Name of the team that owns the cluster")
	viper.BindEnv(flagTeam, envTeam)
	viper.BindPFlag(flagTeam, pushCmd.Flags().Lookup(flagTeam))

	pushCmd.Flags().StringVar(&provider, flagProvider, "", "Provider of the cluster")
	viper.BindEnv(flagProvider, envProvider)
	viper.BindPFlag(flagProvider, pushCmd.Flags().Lookup(flagProvider))

	pushCmd.Flags().StringVar(&awsRegion, flagAWSRegion, "", "AWS region of the cluster")
	viper.BindEnv(flagAWSRegion, envAWSRegion)
	viper.BindPFlag(flagAWSRegion, pushCmd.Flags().Lookup(flagAWSRegion))

	pushCmd.Flags().StringVar(&awsAccountID, flagAWSAccountID, "", "AWS account ID of the cluster")
	viper.BindEnv(flagAWSAccountID, envAWSAccountID)
	viper.BindPFlag(flagAWSAccountID, pushCmd.Flags().Lookup(flagAWSAccountID))
}

func validatePush(cmd *cobra.Command, args []string) error {
	if input != "" {
		_, err := os.Stat(input)
		if err != nil {
			return fmt.Errorf("input file %s can not be accessed.\n%w", input, err)
		}
		log.Debug().Msg(fmt.Sprintf("using input file %s", input))
		return nil
	}
	return nil
}

func defaultPush() {
	if installationsBranch == "" {
		installationsBranch = key.GetDefaultInstallationsBranch(cluster)
	}
}
