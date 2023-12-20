/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/giantswarm/mcli/pkg/key"
	log "github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagCluster             = "cluster"
	flagVerbose             = "verbose"
	flagGithubToken         = "github-token"
	flagSkip                = "skip"
	flagInstallationsBranch = "installations-branch"
)

const (
	envCluster             = "INSTALLATION"
	envGithubToken         = "OPSCTL_GITHUB_TOKEN"
	envInstallationsBranch = "INSTALLATIONS_BRANCH"
)

var (
	cluster             string
	verbose             bool
	githubToken         string
	installationsBranch string
	skip                []string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mcli",
	Short: "A CLI tool to manage Giant Swarm Management Cluster Configuration",
	Long: `A CLI tool to manage Giant Swarm Management Cluster Configuration.
Configuration is stored across multiple git repositories.
This tool allows you to pull and push configuration for new and existing clusters.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	viper.AutomaticEnv()

	rootCmd.PersistentFlags().StringVar(&githubToken, flagGithubToken, viper.GetString(envGithubToken), "Github token to use for authentication")
	rootCmd.PersistentFlags().StringVarP(&cluster, flagCluster, "c", viper.GetString(envCluster), "Name of the management cluster to operate on")
	rootCmd.PersistentFlags().StringVar(&installationsBranch, flagInstallationsBranch, viper.GetString(envInstallationsBranch), "Branch to use for the installations repository")
	rootCmd.PersistentFlags().BoolVarP(&verbose, flagVerbose, "v", false, "Display more verbose output in console output. (default: false)")
}

func toggleVerbose(cmd *cobra.Command, args []string) {
	if verbose {
		log.SetGlobalLevel(log.DebugLevel)
	} else {
		log.SetGlobalLevel(log.ErrorLevel)
	}
}

func validate(cmd *cobra.Command, args []string) error {
	if cluster == "" {
		return invalidFlagError(flagCluster)
	}
	if githubToken == "" {
		return invalidFlagError(flagGithubToken)
	}
	if installationsBranch == "" {
		return invalidFlagError(flagInstallationsBranch)
	}
	for _, repository := range skip {
		if !key.IsValidRepository(repository) {
			return fmt.Errorf("invalid repository %s. Valid values: %s:\n%w", repository, key.GetValidRepositories(), ErrInvalidFlag)
		}
	}
	return nil
}
