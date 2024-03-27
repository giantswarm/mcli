package capv

import (
	"fmt"

	"github.com/giantswarm/mcli/pkg/key"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	CloudConfigKey = "values.yaml"
)

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

	secret := v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "vsphere-credentials",
			Namespace: c.Namespace,
			Labels: map[string]string{
				"clusterctl.cluster.x-k8s.io/move": "true",
			},
		},
		Data: map[string][]byte{
			CloudConfigKey: []byte(c.CloudConfig),
		},
	}
	data, err := yaml.Marshal(secret)
	if err != nil {
		return "", fmt.Errorf("failed to marshal vsphere-credentials object.\n%w", err)
	}
	return string(data), nil
}
