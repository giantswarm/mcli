package base

import (
	"fmt"
	"regexp"

	"github.com/giantswarm/mcli/pkg/template"
	"github.com/rs/zerolog/log"
)

type Config struct {
	PROVIDER                 string
	CMC_REPOSITORY           string
	CMC_BRANCH               string
	INSTALLATION             string
	MCB_BRANCH_SOURCE        string
	CONFIG_BRANCH            string
	MC_APP_COLLECTION_BRANCH string
	BASE_DOMAIN              string
	CATALOG_REGISTRY_VALUES  string
}

func GetBaseFiles(c Config, templates map[string]string) (map[string]string, error) {
	log.Debug().Msg("Getting base files")
	manifests := make(map[string]string)

	for k, t := range templates {

		containsVar, manifest, err := formatTemplate(t)
		if err != nil {
			return nil, fmt.Errorf("failed to format template %s.\n%w", k, err)
		}
		if !containsVar {
			continue
		}
		manifest, err = template.Execute(manifest, c)
		if err != nil {
			return nil, fmt.Errorf("failed to execute template %s.\n%w", k, err)
		}
		manifests[k] = manifest
	}
	return manifests, nil
}

func formatTemplate(t string) (bool, string, error) {
	// if the template contains variables of the format ${VARIABLE_NAME}, format them as {{ .VARIABLE_NAME }}
	regexp := regexp.MustCompile(`\${([A-Z0-9_]+)}`)
	containsVar := regexp.MatchString(t)
	if containsVar {
		t = regexp.ReplaceAllString(t, `{{ .${1} }}`)
	}
	return containsVar, t, nil
}
