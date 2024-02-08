package repositories

import (
	"bytes"
	"fmt"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type Repo struct {
	Name          string          `yaml:"name"`
	ComponentType string          `yaml:"componentType"`
	Gen           Gen             `yaml:"gen,omitempty"`
	Replace       map[string]bool `yaml:"replace,omitempty"`
	Lifecycle     string          `yaml:"lifecycle,omitempty"`
	System        string          `yaml:"system,omitempty"`
}

type Gen struct {
	Flavours                []string `yaml:"flavours,omitempty"`
	Language                string   `yaml:"language,omitempty"`
	EnableFloatingMajorTags bool     `yaml:"enableFloatingMajorTags,omitempty"`
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
	w := new(bytes.Buffer)
	encoder := yaml.NewEncoder(w)
	encoder.SetIndent(2)
	err := encoder.Encode(repo)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal repositories.\n%w", err)
	}
	return w.Bytes(), nil
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
