package registry

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetRegistryConfig(file string) (string, error) {
	log.Debug().Msg("Getting registry config")
	var secret v1.Secret
	err := yaml.Unmarshal([]byte(file), &secret)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal registry object.\n%w", err)
	}
	return string(secret.Data["values.yaml"]), nil
}

func GetRegistryFile(values string) (string, error) {
	log.Debug().Msg("Creating container-registries-configuration Secret")
	secret := v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "container-registries-configuration",
			Namespace: "default",
		},
		Data: map[string][]byte{
			"values.yaml": []byte(values),
		},
	}
	data, err := yaml.Marshal(secret)
	if err != nil {
		return "", fmt.Errorf("failed to marshal container-registries-configuration object.\n%w", err)
	}
	return string(data), nil
}
