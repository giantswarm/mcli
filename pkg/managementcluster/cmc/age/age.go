package age

import (
	"github.com/rs/zerolog/log"

	"github.com/giantswarm/mcli/pkg/key"
	"github.com/giantswarm/mcli/pkg/template"
)

const Agetemplate = `apiVersion: v1
kind: Secret
type: Opaque
metadata:
    name: sops-keys
    namespace: flux-giantswarm
data:
    {{ .Cluster }}.agekey: {{ .AgeKey }}
`

type Config struct {
	Cluster string
	AgeKey  string
}

func GetAgeKey(file string, cluster string) (string, error) {
	log.Debug().Msg("Getting Age key")

	return key.GetSecretValue(key.GetAgeKey(cluster), file)
}

func GetAgeFile(c Config) (string, error) {
	log.Debug().Msg("Creating Age file")

	return template.Execute(Agetemplate, c)
}
