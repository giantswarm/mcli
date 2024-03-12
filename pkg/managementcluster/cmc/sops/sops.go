package sops

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Cluster   string
	AgePubKey string
}

type Sops struct {
	CreationRules []CreationRule `yaml:"creation_rules"`
}

type CreationRule struct {
	Age            string `yaml:"age"`
	PathRegex      string `yaml:"path_regex"`
	EncryptedRegex string `yaml:"encrypted_regex"`
}

func GetSopsFile(c Config, file string) (string, error) {
	log.Debug().Msg(fmt.Sprintf("Adding SOPS pubkey for the installation %s", c.Cluster))
	sops, err := getSops(file)
	if err != nil {
		return "", err
	}
	for i, rule := range sops.CreationRules {
		if rule.PathRegex == getRegex(c.Cluster) {
			// delete existing entry
			sops.CreationRules = append(sops.CreationRules[:i], sops.CreationRules[i+1:]...)
		}
	}
	sops.CreationRules = append(sops.CreationRules, CreationRule{
		Age:            c.AgePubKey,
		PathRegex:      getRegex(c.Cluster),
		EncryptedRegex: "^(data|stringData)$",
	})
	data, err := yaml.Marshal(sops)
	if err != nil {
		return "", fmt.Errorf("failed to marshal sops object.\n%w", err)
	}
	return string(data), nil
}

func GetSopsConfig(file string, cluster string) (Config, error) {
	log.Debug().Msg(fmt.Sprintf("Getting SOPS pubkey for the installation %s", cluster))
	sops, err := getSops(file)
	if err != nil {
		return Config{}, err
	}
	for _, rule := range sops.CreationRules {
		if rule.PathRegex == getRegex(cluster) {
			return Config{
				Cluster:   cluster,
				AgePubKey: rule.Age,
			}, nil
		}
	}
	return Config{}, fmt.Errorf("no SOPS pubkey found for the installation %s", cluster)
}

func getSops(file string) (Sops, error) {
	sops := Sops{}
	if err := yaml.Unmarshal([]byte(file), &sops); err != nil {
		return sops, fmt.Errorf("failed to unmarshal sops object.\n%w", err)
	}
	return sops, nil
}

func getRegex(cluster string) string {
	return fmt.Sprintf("management-clusters/%s/.*(secret|credential).*", cluster)
}
