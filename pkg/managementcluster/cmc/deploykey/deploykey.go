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
	Passphrasekey       = "password"
	Identitykey         = "identity"
	knownhostskey       = "known_hosts"
)

type Config struct {
	Name       string
	Passphrase string
	Identity   string
	KnownHosts string
}

func GetDeployKeyConfig(file string) (Config, error) {
	log.Debug().Msg("Getting DeployKey config")

	passphrase, err := key.GetSecretValue(Passphrasekey, file)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get DeployKey passphrase.\n%w", err)
	}

	identity, err := key.GetSecretValue(Identitykey, file)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get DeployKey identity.\n%w", err)
	}

	knownhosts, err := key.GetSecretValue(knownhostskey, file)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get DeployKey knownhosts.\n%w", err)
	}

	return Config{
		Passphrase: passphrase,
		Identity:   identity,
		KnownHosts: knownhosts,
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
			Passphrasekey: []byte(c.Passphrase),
			Identitykey:   []byte(c.Identity),
			knownhostskey: []byte(c.KnownHosts),
		},
	}
	data, err := yaml.Marshal(secret)
	if err != nil {
		return "", fmt.Errorf("failed to marshal DeployKey object %s.\n%w", c.Name, err)
	}
	return string(data), nil
}
