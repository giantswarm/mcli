package repositories

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"

	"github.com/giantswarm/mcli/pkg/key"
)

type Repo struct {
	Name            string          `yaml:"name"`
	ComponentType   string          `yaml:"componentType"`
	DeploymentNames []string        `yaml:"deploymentNames,omitempty"`
	Gen             Gen             `yaml:"gen,omitempty"`
	Replace         map[string]bool `yaml:"replace,omitempty"`
	Lifecycle       string          `yaml:"lifecycle,omitempty"`
	System          string          `yaml:"system,omitempty"`
}

type Gen struct {
	Flavours                []string `yaml:"flavours,omitempty"`
	Language                string   `yaml:"language,omitempty"`
	EnableFloatingMajorTags bool     `yaml:"enableFloatingMajorTags,omitempty"`
	InstallUpdateChart      bool     `yaml:"installUpdateChart,omitempty"`
	RunSecurityScoreCard    bool     `yaml:"runSecurityScoreCard,omitempty"`
}

func GetRepos(data []byte) ([]Repo, error) {
	repositories := []Repo{}
	if err := yaml.Unmarshal(data, &repositories); err != nil {
		return nil, fmt.Errorf("failed to unmarshal repositories.\n%w", err)
	}
	return repositories, nil
}

func GetData(repo []Repo) ([]byte, error) {
	log.Debug().Msg("getting data from repositories object")

	return key.GetData(repo)
}

func SortReposAlphabetically(repos []Repo) []Repo {
	for i := 0; i < len(repos); i++ {
		for j := i + 1; j < len(repos); j++ {
			if repos[i].Name > repos[j].Name {
				repos[i], repos[j] = repos[j], repos[i]
			}
		}
	}
	return repos
}
