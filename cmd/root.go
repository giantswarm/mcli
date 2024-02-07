/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	log "github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mcli",
	Short: "A CLI tool to manage Giant Swarm Management Cluster Configuration",
	Long: `A CLI tool to manage Giant Swarm Management Cluster Configuration.
Configuration is stored across multiple git repositories.
This tool allows you to pull and push configuration for new and existing clusters.`,
	PersistentPreRun: toggleVerbose,
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
	addFlagsRoot()
}

func toggleVerbose(cmd *cobra.Command, args []string) {
	if verbose {
		log.SetGlobalLevel(log.DebugLevel)
	} else {
		log.SetGlobalLevel(log.ErrorLevel)
	}
}
