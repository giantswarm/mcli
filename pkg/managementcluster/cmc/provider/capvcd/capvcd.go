package capvcd

import (
	"github.com/rs/zerolog/log"

	"github.com/giantswarm/mcli/pkg/key"
	"github.com/giantswarm/mcli/pkg/template"
)

const (
	RefreshTokenKey = "refreshToken"
)

const CloudDirectorTemplate = `apiVersion: v1
kind: Secret
metadata:
  name: vcd-credentials
  namespace: {{ .Namespace }}
  labels:
    clusterctl.cluster.x-k8s.io/move: true
data:
  refreshToken: {{ .RefreshToken }}
`

type Config struct {
	Namespace    string
	RefreshToken string
}

func GetCAPVCDConfig(file string) (string, error) {
	log.Debug().Msg("Getting CAPVCD config")

	return key.GetSecretValue(RefreshTokenKey, file)
}

func GetCAPVCDFile(c Config) (string, error) {
	log.Debug().Msg("Creating CAPV file")

	return template.Execute(CloudDirectorTemplate, c)
}
