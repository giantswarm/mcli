package deploykey

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/giantswarm/mcli/pkg/key"
	"github.com/giantswarm/mcli/pkg/template"
)

const (
	Passphrasekey = "password"
	Identitykey   = "identity"
	knownhostskey = "known_hosts"
)

const DeployKeyTemplate = `apiVersion: v1
kind: Secret
type: Opaque
metadata:
  name: {{ .Name }}
  namespace: flux-giantswarm
data:
  password: {{ .Passphrase }}
  identity: {{ .Identity }}
  known_hosts: {{ .KnownHosts }}
`

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

	return template.Execute(DeployKeyTemplate, c)
}
