package installations

import (
	"bytes"
	"fmt"
	"os"
	"reflect"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"

	"github.com/giantswarm/mcli/pkg/key"
)

type Installations struct {
	Base            string    `yaml:"base"`
	Codename        string    `yaml:"codename"`
	Customer        string    `yaml:"customer"`
	CmcRepository   string    `yaml:"cmc_repository"`
	AccountEngineer string    `yaml:"accountEngineer"`
	Pipeline        string    `yaml:"pipeline"`
	Provider        string    `yaml:"provider"`
	Aws             AwsConfig `yaml:"aws,omitempty"`
}

type AwsConfig struct {
	Region       string       `yaml:"region"`
	HostCluster  HostCluster  `yaml:"hostCluster"`
	GuestCluster GuestCluster `yaml:"guestCluster"`
}

type HostCluster struct {
	Account          string `yaml:"account"`
	CloudtrailBucket string `yaml:"cloudtrailBucket"`
	AdminRoleArn     string `yaml:"adminRoleARN"`
	GuardDuty        bool   `yaml:"guardDuty"`
}

type GuestCluster struct {
	Account          string `yaml:"account"`
	CloudtrailBucket string `yaml:"cloudtrailBucket"`
	GuardDuty        bool   `yaml:"guardDuty"`
}

type InstallationsConfig struct {
	Base                            string
	Codename                        string
	Customer                        string
	CmcRepository                   string
	AccountEngineer                 string
	Pipeline                        string
	Provider                        string
	AwsRegion                       string
	AwsHostClusterAccount           string
	AwsHostClusterAdminRoleArn      string
	AwsHostClusterGuardDuty         bool
	AwsHostClusterCloudtrailBucket  string
	AwsGuestClusterAccount          string
	AwsGuestClusterCloudtrailBucket string
	AwsGuestClusterGuardDuty        bool
}

func NewInstallations(installationsConfig InstallationsConfig) *Installations {
	return &Installations{
		Base:            installationsConfig.Base,
		Codename:        installationsConfig.Codename,
		Customer:        installationsConfig.Customer,
		CmcRepository:   installationsConfig.CmcRepository,
		AccountEngineer: installationsConfig.AccountEngineer,
		Pipeline:        installationsConfig.Pipeline,
		Provider:        installationsConfig.Provider,
		Aws: AwsConfig{
			Region: installationsConfig.AwsRegion,
			HostCluster: HostCluster{
				Account:          installationsConfig.AwsHostClusterAccount,
				AdminRoleArn:     installationsConfig.AwsHostClusterAdminRoleArn,
				CloudtrailBucket: installationsConfig.AwsHostClusterCloudtrailBucket,
				GuardDuty:        installationsConfig.AwsHostClusterGuardDuty,
			},
			GuestCluster: GuestCluster{
				Account:          installationsConfig.AwsGuestClusterAccount,
				CloudtrailBucket: installationsConfig.AwsGuestClusterCloudtrailBucket,
				GuardDuty:        installationsConfig.AwsGuestClusterGuardDuty,
			},
		},
	}
}

func GetInstallations(data []byte) (*Installations, error) {
	log.Debug().Msg("getting installations object from data")
	installations := Installations{}
	if err := yaml.Unmarshal(data, &installations); err != nil {
		return nil, fmt.Errorf("failed to unmarshal installations object.\n%w", err)
	}
	return &installations, nil
}

func GetData(i *Installations) ([]byte, error) {
	log.Debug().Msg("getting data from installations object")
	w := new(bytes.Buffer)
	encoder := yaml.NewEncoder(w)
	encoder.SetIndent(2)
	err := encoder.Encode(i)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal installations object.\n%w", err)
	}
	return w.Bytes(), nil
}

func (i *Installations) Print() error {
	data, err := GetData(i)
	if err != nil {
		return err
	}
	log.Debug().Msg("printing installations object")
	fmt.Print(string(data))
	return nil
}

func (i *Installations) Override(override *Installations) *Installations {
	installation := *i
	if override.Base != "" {
		installation.Base = override.Base
	}
	if override.Codename != "" {
		installation.Codename = override.Codename
	}
	if override.Customer != "" {
		installation.Customer = override.Customer
	}
	if override.CmcRepository != "" {
		installation.CmcRepository = override.CmcRepository
	}
	if override.AccountEngineer != "" {
		installation.AccountEngineer = override.AccountEngineer
	}
	if override.Pipeline != "" {
		installation.Pipeline = override.Pipeline
	}
	if override.Provider != "" {
		installation.Provider = override.Provider
	}
	if override.Aws.Region != "" {
		installation.Aws.Region = override.Aws.Region
	}
	if override.Aws.HostCluster.Account != "" {
		installation.Aws.HostCluster.Account = override.Aws.HostCluster.Account
	}
	if override.Aws.HostCluster.AdminRoleArn != "" {
		installation.Aws.HostCluster.AdminRoleArn = override.Aws.HostCluster.AdminRoleArn
	}
	if override.Aws.HostCluster.GuardDuty {
		installation.Aws.HostCluster.GuardDuty = override.Aws.HostCluster.GuardDuty
	}
	if override.Aws.GuestCluster.Account != "" {
		installation.Aws.GuestCluster.Account = override.Aws.GuestCluster.Account
	}
	if override.Aws.GuestCluster.CloudtrailBucket != "" {
		installation.Aws.GuestCluster.CloudtrailBucket = override.Aws.GuestCluster.CloudtrailBucket
	}
	if override.Aws.GuestCluster.GuardDuty {
		installation.Aws.GuestCluster.GuardDuty = override.Aws.GuestCluster.GuardDuty
	}
	return &installation
}

func (i *Installations) Validate() error {
	if i.Base == "" {
		return fmt.Errorf("base domain is empty")
	}
	if i.Codename == "" {
		return fmt.Errorf("codename is empty")
	}
	if i.Customer == "" {
		return fmt.Errorf("customer is empty")
	}
	if i.CmcRepository == "" {
		return fmt.Errorf("cmc repository is empty")
	}
	if i.AccountEngineer == "" {
		return fmt.Errorf("account engineer is empty")
	}
	if i.Pipeline == "" {
		return fmt.Errorf("pipeline is empty")
	}
	if i.Provider == "" {
		return fmt.Errorf("provider is empty")
	}
	if key.IsProviderAWS(i.Provider) {
		if i.Aws.Region == "" {
			return fmt.Errorf("aws region is empty")
		}
		if i.Aws.HostCluster.Account == "" {
			return fmt.Errorf("aws host cluster account is empty")
		}
		if i.Aws.HostCluster.AdminRoleArn == "" {
			return fmt.Errorf("aws host cluster admin role arn is empty")
		}
		if i.Aws.GuestCluster.Account == "" {
			return fmt.Errorf("aws guest cluster account is empty")
		}
	}
	return nil
}

func (i *Installations) Equals(other *Installations) bool {
	return reflect.DeepEqual(i, other)
}

func GetInstallationsFromFile(path string) (*Installations, error) {
	log.Debug().Msg(fmt.Sprintf("getting installations object from file %s", path))
	// read data from input file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read input file %s.\n%w", path, err)
	}

	return GetInstallations(data)
}
