package cmd

import (
	"fmt"

	"github.com/giantswarm/mcli/pkg/key"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagCluster             = "cluster"
	flagVerbose             = "verbose"
	flagGithubToken         = "github-token"
	flagSkip                = "skip"
	flagInstallationsBranch = "installations-branch"
	flagCMCRepository       = "cmc-repository"
	flagCMCBranch           = "cmc-branch"
	flagCustomer            = "customer"
)

const (
	envCluster             = "INSTALLATION"
	envGithubToken         = "OPSCTL_GITHUB_TOKEN"
	envInstallationsBranch = "INSTALLATIONS_BRANCH"
	envCMCRepository       = "CMC_REPOSITORY"
	envCMCBranch           = "CMC_BRANCH"
	envCustomer            = "CUSTOMER"
)

var (
	cluster             string
	verbose             bool
	githubToken         string
	installationsBranch string
	skip                []string
	cmcRepository       string
	cmcBranch           string
	customer            string
)

func addFlagsRoot() {
	viper.AutomaticEnv()

	rootCmd.PersistentFlags().StringVar(&githubToken, flagGithubToken, viper.GetString(envGithubToken), "Github token to use for authentication")
	rootCmd.PersistentFlags().MarkHidden(flagGithubToken)
	rootCmd.PersistentFlags().StringVarP(&cluster, flagCluster, "c", viper.GetString(envCluster), "Name of the management cluster to operate on")
	rootCmd.PersistentFlags().StringVar(&installationsBranch, flagInstallationsBranch, viper.GetString(envInstallationsBranch), "Branch to use for the installations repository")
	rootCmd.PersistentFlags().BoolVarP(&verbose, flagVerbose, "v", false, "Display more verbose output in console output. (default: false)")
	rootCmd.PersistentFlags().StringVar(&cmcRepository, flagCMCRepository, viper.GetString(envCMCRepository), "Name of CMC repository to use")
	rootCmd.PersistentFlags().StringVar(&cmcBranch, flagCMCBranch, viper.GetString(envCMCBranch), "Branch to use for the CMC repository")
	rootCmd.PersistentFlags().StringVar(&customer, flagCustomer, viper.GetString(envCustomer), "Name of the customer who owns the management cluster")
}

func validateRoot(cmd *cobra.Command, args []string) error {
	if cluster == "" {
		return invalidFlagError(flagCluster)
	}
	if githubToken == "" {
		return invalidFlagError(flagGithubToken)
	}
	if cmcRepository == "" {
		return invalidFlagError(flagCMCRepository)
	}
	if cmcBranch == "" {
		return invalidFlagError(flagCMCBranch)
	}
	for _, repository := range skip {
		if !key.IsValidRepository(repository) {
			return fmt.Errorf("invalid repository %s. Valid values: %s:\n%w", repository, key.GetValidRepositories(), ErrInvalidFlag)
		}
	}
	return nil
}
