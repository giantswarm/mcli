package taylorbot

import (
	"github.com/rs/zerolog/log"

	"github.com/giantswarm/mcli/pkg/key"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc/kustomization"
	"github.com/giantswarm/mcli/pkg/template"
)

const (
	PasswordKey = "password"
)

type Config struct {
	Token string
}

func GetTaylorBotToken(file string) (string, error) {
	log.Debug().Msg("Getting TaylorBot token")

	return key.GetSecretValue(PasswordKey, file)
}

func GetTaylorBotFile(token string) (string, error) {
	log.Debug().Msg("Creating TaylorBot file")

	return template.Execute(key.GetTMPLFile(kustomization.TaylorBotFile), Config{
		Token: token,
	})
}
