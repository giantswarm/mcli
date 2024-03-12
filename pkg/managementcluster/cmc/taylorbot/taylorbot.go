package taylorbot

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/giantswarm/mcli/pkg/key"
)

const (
	TaylorBotSecretName = "github-giantswarm-https-credentials"
	TaylorBotSecretURL  = "https://github.com/giantswarm"
	TaylorBotUsername   = "taylorbotgit"
)

func GetTaylorBotToken(file string) (string, error) {
	log.Debug().Msg("Getting TaylorBot token")
	var secret v1.Secret
	err := yaml.Unmarshal([]byte(file), &secret)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal TaylorBot object.\n%w", err)
	}
	return string(secret.Data["password"]), nil
}

func GetTaylorBotFile(token string) (string, error) {
	log.Debug().Msg("Creating TaylorBot file")

	secret := v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      TaylorBotSecretName,
			Namespace: key.FluxNamespace,
		},
		Data: map[string][]byte{
			"url":      []byte(TaylorBotSecretURL),
			"username": []byte(TaylorBotUsername),
			"password": []byte(token),
		},
	}
	data, err := yaml.Marshal(secret)
	if err != nil {
		return "", fmt.Errorf("failed to marshal TaylorBot object.\n%w", err)
	}
	return string(data), nil
}
