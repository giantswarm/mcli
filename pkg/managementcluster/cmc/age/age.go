package age

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/mcli/pkg/key"
)

const (
	AgeSecretName = "sops-keys"
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

	secret := v1.Secret{
		Type: v1.SecretTypeOpaque,
		ObjectMeta: metav1.ObjectMeta{
			Name:      AgeSecretName,
			Namespace: key.FluxNamespace,
		},
		Data: map[string][]byte{
			key.GetAgeKey(c.Cluster): []byte(c.AgeKey),
		},
	}
	data, err := yaml.Marshal(secret)
	if err != nil {
		return "", fmt.Errorf("failed to marshal Age object.\n%w", err)
	}
	return string(data), nil
}
