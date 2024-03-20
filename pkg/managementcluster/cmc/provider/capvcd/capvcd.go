package capvcd

import (
	"encoding/base64"
	"fmt"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Config struct {
	Namespace    string
	RefreshToken string
}

func GetCAPVCDConfig(file string) Config {
	log.Debug().Msg("Getting CAPVCD config")
	var secret v1.Secret
	err := yaml.Unmarshal([]byte(file), &secret)
	if err != nil {
		log.Error().Msgf("failed to unmarshal CAPVCD object.\n%v", err)
	}
	return Config{
		Namespace:    secret.ObjectMeta.Namespace,
		RefreshToken: string(secret.Data["refreshToken"]),
	}
}

func GetCAPVCDFile(c Config) (string, error) {
	log.Debug().Msg("Creating CAPV file")

	//b64 encode the refresh token
	refreshToken := base64.StdEncoding.EncodeToString([]byte(c.RefreshToken))

	secret := v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "vcd-credentials",
			Namespace: c.Namespace,
			Labels: map[string]string{
				"clusterctl.cluster.x-k8s.io/move": "true",
			},
		},
		Data: map[string][]byte{
			"refreshToken": []byte(refreshToken),
		},
	}
	data, err := yaml.Marshal(secret)
	if err != nil {
		return "", fmt.Errorf("failed to marshal cloud-director-credentials object.\n%w", err)
	}
	return string(data), nil
}
