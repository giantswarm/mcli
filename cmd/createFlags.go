package cmd

import (
	"github.com/giantswarm/mcli/pkg/key"
)

func defaultCreate() {
	if cluster == "" {
		cluster = "default"
	}
	if cmcBranch == "" {
		cmcBranch = key.CMCMainBranch
	}
	if cmcRepository == "" {
		if customer != "" {
			cmcRepository = key.GetCMCName(customer)
		}
	}
}
