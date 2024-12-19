/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/giantswarm/mcli/pkg/key"
)

func addFlagsPull() {
	pullCmd.Flags().StringArrayVarP(&skip, flagSkip, "s", []string{}, fmt.Sprintf("List of repositories to skip. (default: none) Valid values: %s", key.GetValidRepositories()))
}

func defaultPull() {
	if installationsBranch == "" {
		installationsBranch = key.InstallationsMainBranch
	}
	if cmcBranch == "" {
		cmcBranch = key.CMCMainBranch
	}
	if customer == "" {
		customer = key.OrganizationGiantSwarm
	}
	if cmcRepository == "" {
		cmcRepository = key.GetCMCName(customer)
	}
}
