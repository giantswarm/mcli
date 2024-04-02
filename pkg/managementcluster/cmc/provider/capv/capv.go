package capv

import (
	"github.com/rs/zerolog/log"

	"github.com/giantswarm/mcli/pkg/key"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc/kustomization"
	"github.com/giantswarm/mcli/pkg/template"
)

const (
	CloudConfigKey = "values.yaml"
)

type Config struct {
	CloudConfig string
	Namespace   string
}

func GetCAPVConfig(file string) (string, error) {
	log.Debug().Msg("Getting CAPV config")

	return key.GetSecretValue(CloudConfigKey, file)
}

func GetCAPVFile(c Config) (string, error) {
	log.Debug().Msg("Creating CAPV file")

	return template.Execute(key.GetTMPLFile(kustomization.VsphereCredentialsFile), c)
}
