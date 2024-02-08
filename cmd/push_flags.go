/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/giantswarm/mcli/pkg/key"
)

const (
	flagInput        = "input"
	flagBaseDomain   = "base-domain"
	flagTeam         = "team"
	flagProvider     = "provider"
	flagAWSRegion    = "aws-region"
	flagAWSAccountID = "aws-account-id"
)

const (
	envBaseDomain   = "BASE_DOMAIN"
	envTeam         = "TEAM_NAME"
	envProvider     = "PROVIDER"
	envAWSRegion    = "AWS_REGION"
	envAWSAccountID = "INSTALLATION_AWS_ACCOUNT"
)

var (
	input        string
	baseDomain   string
	team         string
	provider     string
	awsRegion    string
	awsAccountID string
)

func addFlagsPush() {
	pushCmd.Flags().StringArrayVarP(&skip, flagSkip, "s", []string{}, fmt.Sprintf("List of repositories to skip. (default: none) Valid values: %s", key.GetValidRepositories()))
	pushCmd.PersistentFlags().StringVarP(&input, flagInput, "i", "", "Input configuration file to use. If not specified, configuration is read from other flags.")
	pushCmd.PersistentFlags().StringVar(&baseDomain, flagBaseDomain, viper.GetString(envBaseDomain), "Base domain to use for the cluster")
	pushCmd.PersistentFlags().StringVar(&team, flagTeam, viper.GetString(envTeam), "Name of the team that owns the cluster")
	pushCmd.PersistentFlags().StringVar(&provider, flagProvider, viper.GetString(envProvider), "Provider of the cluster")
	pushCmd.PersistentFlags().StringVar(&awsRegion, flagAWSRegion, viper.GetString(envAWSRegion), "AWS region of the cluster")
	pushCmd.PersistentFlags().StringVar(&awsAccountID, flagAWSAccountID, viper.GetString(envAWSAccountID), "AWS account ID of the cluster")
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
	if installationsBranch == "" {
		return invalidFlagError(flagInstallationsBranch)
	}
	return nil
}

func defaultPush() {
	if installationsBranch == "" {
		installationsBranch = key.GetDefaultInstallationsBranch(cluster)
	}
	if customer == "" {
		customer = "giantswarm"
	}
}
