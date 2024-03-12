package deploykey

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/mcli/pkg/key"
)

const (
	DeployKeySecretName = "configs-ssh-credentials"
)

type Config struct {
	Name       string
	Passphrase string
	Identity   string
	KnownHosts string
}

func GetDeployKeyConfig(file string) (Config, error) {
	log.Debug().Msg("Getting DeployKey config")
	var secret v1.Secret
	err := yaml.Unmarshal([]byte(file), &secret)
	if err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal DeployKey object.\n%w", err)
	}
	return Config{
		Name:       secret.ObjectMeta.Name,
		Passphrase: string(secret.Data["passphrase"]),
		Identity:   string(secret.Data["identity"]),
		KnownHosts: string(secret.Data["knownHosts"]),
	}, nil
}

func GetDeployKeyFile(c Config) (string, error) {
	log.Debug().Msg(fmt.Sprintf("Creating DeployKey file for %s", c.Name))

	secret := v1.Secret{
		Type: v1.SecretTypeOpaque,
		ObjectMeta: metav1.ObjectMeta{
			Name:      c.Name,
			Namespace: key.FluxNamespace,
		},
		Data: map[string][]byte{
			"passphrase": []byte(c.Passphrase),
			"identity":   []byte(c.Identity),
			"knownHosts": []byte(c.KnownHosts),
		},
	}
	data, err := yaml.Marshal(secret)
	if err != nil {
		return "", fmt.Errorf("failed to marshal DeployKey object %s.\n%w", c.Name, err)
	}
	return string(data), nil
}
