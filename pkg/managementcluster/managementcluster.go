package managementcluster

import (
	"fmt"
	"os"

	"github.com/giantswarm/mcli/pkg/managementcluster/cmc"
	"github.com/giantswarm/mcli/pkg/managementcluster/installations"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type ManagementCluster struct {
	Installations installations.Installations `yaml:"installations,omitempty"`
	CMC           cmc.CMC                     `yaml:"cmc,omitempty"`
}

func (mc *ManagementCluster) Print() error {
	data, err := mc.GetData()
	if err != nil {
		return err
	}
	fmt.Print(string(data))
	return nil
}

func (mc *ManagementCluster) GetData() ([]byte, error) {
	log.Debug().Msg("getting management cluster data")
	data, err := yaml.Marshal(mc)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal management cluster object.\n%w", err)
	}
	return data, nil
}

func GetManagementCluster(data []byte) (*ManagementCluster, error) {
	log.Debug().Msg("getting management cluster object from data")
	managementcluster := ManagementCluster{}
	if err := yaml.Unmarshal(data, &managementcluster); err != nil {
		return nil, fmt.Errorf("failed to unmarshal management cluster object.\n%w", err)
	}
	return &managementcluster, nil
}

func GetManagementClusterFromFile(input string) (*ManagementCluster, error) {
	log.Debug().Msg("getting new management cluster object from input file")

	// read data from input file
	data, err := os.ReadFile(input)
	if err != nil {
		return nil, fmt.Errorf("failed to read input file %s.\n%w", input, err)
	}

	return GetManagementCluster(data)
}
