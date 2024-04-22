package taylorbot

import (
	"github.com/rs/zerolog/log"

	"github.com/giantswarm/mcli/pkg/key"
	"github.com/giantswarm/mcli/pkg/template"
)

const (
	PasswordKey = "password"
)

const TaylorBotTemplate = `apiVersion: v1
kind: Secret
metadata:
  name: github-giantswarm-https-credentials
  namespace: flux-giantswarm
data:
  password: {{ .Token }}
  url: https://github.com/giantswarm
  username: taylorbotgit
`

type Config struct {
	Token string
}

func GetTaylorBotToken(file string) (string, error) {
	log.Debug().Msg("Getting TaylorBot token")

	return key.GetSecretValue(PasswordKey, file)
}

func GetTaylorBotFile(token string) (string, error) {
	log.Debug().Msg("Creating TaylorBot file")

	return template.Execute(TaylorBotTemplate, Config{
		Token: token,
	})
}
