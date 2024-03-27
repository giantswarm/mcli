package registry

import (
	"fmt"

	"github.com/giantswarm/mcli/pkg/key"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ValuesKey = "values.yaml"
)

func GetRegistryConfig(file string) (string, error) {
	log.Debug().Msg("Getting registry config")

	return key.GetSecretValue(ValuesKey, file)
}

func GetRegistryFile(values string) (string, error) {
	log.Debug().Msg("Creating container-registries-configuration Secret")
	secret := v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "container-registries-configuration",
			Namespace: "default",
		},
		Data: map[string][]byte{
			ValuesKey: []byte(values),
		},
	}
	data, err := yaml.Marshal(secret)
	if err != nil {
		return "", fmt.Errorf("failed to marshal container-registries-configuration object.\n%w", err)
	}
	return string(data), nil
}
