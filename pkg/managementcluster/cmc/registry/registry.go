package registry

import (
	"github.com/rs/zerolog/log"

	"github.com/giantswarm/mcli/pkg/key"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc/kustomization"
	"github.com/giantswarm/mcli/pkg/template"
)

const (
	ValuesKey = "values.yaml"
)

type Config struct {
	Values string
}

func GetRegistryConfig(file string) (string, error) {
	log.Debug().Msg("Getting registry config")

	return key.GetSecretValue(ValuesKey, file)
}

func GetRegistryFile(values string) (string, error) {
	log.Debug().Msg("Creating container-registries-configuration Secret")
	return template.Execute(key.GetTMPLFile(kustomization.RegistryFile), Config{Values: values})
}
