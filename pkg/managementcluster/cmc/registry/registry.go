package registry

import (
	"github.com/rs/zerolog/log"

	"github.com/giantswarm/mcli/pkg/key"
	"github.com/giantswarm/mcli/pkg/template"
)

const (
	ValuesKey = "values.yaml"
)

const RegistryTemplate = `apiVersion: v1
kind: Secret
metadata:
  name: container-registries-configuration
  namespace: default
data:
  values.yaml: {{ .Values }}
`

type Config struct {
	Values string
}

func GetRegistryConfig(file string) (string, error) {
	log.Debug().Msg("Getting registry config")

	return key.GetSecretValue(ValuesKey, file)
}

func GetRegistryFile(values string) (string, error) {
	log.Debug().Msg("Creating container-registries-configuration Secret")
	return template.Execute(RegistryTemplate, Config{Values: values})
}
