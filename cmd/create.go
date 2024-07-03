package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	createcmc "github.com/giantswarm/mcli/cmd/create/cmc"
	"github.com/giantswarm/mcli/pkg/github"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates repository to hold configuration of the Management Cluster",
	Long:  `Creates repository to hold configuration of the Management Cluster, such as CMC repository.`,
}

var createCMCCmd = &cobra.Command{
	Use:   "cmc",
	Short: "Creates CMC repository",
	Long: `Creates CMC repository. For example:

mcli create cmc --customer=gigamac`,
	RunE: func(cmd *cobra.Command, args []string) error {
		defaultCreate()
		err := validateRoot(cmd, args)
		if err != nil {
			return err
		}
		ctx := context.Background()
		client := github.New(github.Config{
			Token: githubToken,
		})
		c := createcmc.Config{
			Github:        client,
			CMCRepository: cmcRepository,
			CMCBranch:     cmcBranch,
			Customer:      customer,
		}
		_, err = c.Run(ctx)
		if err != nil {
			return fmt.Errorf("failed to create CMC repository.\n%w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.AddCommand(createCMCCmd)
}
