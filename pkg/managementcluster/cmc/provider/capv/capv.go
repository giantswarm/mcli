package capv

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Config struct {
	Namespace   string
	CloudConfig string
}

func GetCAPVConfig(file string) Config {
	log.Debug().Msg("Getting CAPV config")
	var secret v1.Secret
	err := yaml.Unmarshal([]byte(file), &secret)
	if err != nil {
		log.Error().Msgf("failed to unmarshal CAPV object.\n%v", err)
	}
	return Config{
		Namespace:   secret.ObjectMeta.Namespace,
		CloudConfig: string(secret.Data["values.yaml"]),
	}
}

func GetCAPVFile(c Config) (string, error) {
	log.Debug().Msg("Creating CAPV file")

	secret := v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "vsphere-credentials",
			Namespace: c.Namespace,
			Labels: map[string]string{
				"clusterctl.cluster.x-k8s.io/move": "true",
			},
		},
		Data: map[string][]byte{
			"values.yaml": []byte(c.CloudConfig),
		},
	}
	data, err := yaml.Marshal(secret)
	if err != nil {
		return "", fmt.Errorf("failed to marshal vsphere-credentials object.\n%w", err)
	}
	return string(data), nil
}
