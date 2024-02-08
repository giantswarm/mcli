package cmc

import (
	"fmt"

	flux "github.com/fluxcd/kustomize-controller/api/v1beta2"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type CMC struct {
	//TODO: add fields, the map is a placeholder for now
	Contents map[string]string
}

func GetKustomization(data []byte) (*flux.Kustomization, error) {
	log.Debug().Msg("getting kustomization object from data")
	kustomization := flux.Kustomization{}
	if err := yaml.Unmarshal(data, &kustomization); err != nil {
		return nil, fmt.Errorf("failed to unmarshal kustomization object.\n%w", err)
	}
	return &kustomization, nil
}

func (c *CMC) Print() error {
	return nil
}
