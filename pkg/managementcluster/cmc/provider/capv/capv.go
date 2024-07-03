package capv

import (
	"github.com/rs/zerolog/log"

	"github.com/giantswarm/mcli/pkg/key"
	"github.com/giantswarm/mcli/pkg/template"
)

const (
	CloudConfigKey = "values.yaml"
)

const VsphereTemplate = `apiVersion: v1
kind: Secret
metadata:
  name: vsphere-credentials
  namespace: {{ .Namespace }}
  labels:
    clusterctl.cluster.x-k8s.io/move: true
data:
  values.yaml: {{ .CloudConfig }}
`

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

	return template.Execute(VsphereTemplate, c)
}
