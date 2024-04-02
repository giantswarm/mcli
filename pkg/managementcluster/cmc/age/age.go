package age

import (
	"github.com/rs/zerolog/log"

	"github.com/giantswarm/mcli/pkg/key"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc/kustomization"
	"github.com/giantswarm/mcli/pkg/template"
)

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

	return template.Execute(key.GetTMPLFile(kustomization.AgeKeyFile), c)
}
