package base

import (
	"fmt"
	"regexp"

	"github.com/rs/zerolog/log"

	"github.com/giantswarm/mcli/pkg/key"
	"github.com/giantswarm/mcli/pkg/template"
)

const (
	baseDomainKey = "baseDomain"
	branchKey     = "branch"
	domainKey     = "domain"
)

const (
	customBranchCollectionFile = "custom-branch-collection.yaml"
	customBranchCMCFile        = "custom-branch-management-clusters-fleet.yaml"
	customBranchConfigFile     = "custom-branch-config.yaml"
	catalogPatchFile           = "catalogs/patches/appcatalog-default-patch.yaml"
	catalogKustomizationFile   = "catalogs/kustomization.yaml"
)

type Config struct {
	Provider              string
	Cluster               string
	CMCRepository         string
	CMCBranch             string
	MCBBranchSource       string
	ConfigBranch          string
	MCAppCollectionBranch string
	BaseDomain            string
	RegistryDomain        string
}

type Template struct {
	PROVIDER                 string
	INSTALLATION             string
	CMC_REPOSITORY           string
	CMC_BRANCH               string
	MCB_BRANCH_SOURCE        string
	CONFIG_BRANCH            string
	MC_APP_COLLECTION_BRANCH string
	BASE_DOMAIN              string
	CATALOG_REGISTRY_VALUES  string
}

func GetBaseFiles(c Config, templates map[string]string) (map[string]string, error) {
	log.Debug().Msg("Getting base files")
	manifests := make(map[string]string)

	tmp := Template{
		PROVIDER:                 c.Provider,
		INSTALLATION:             c.Cluster,
		CMC_REPOSITORY:           c.CMCRepository,
		CMC_BRANCH:               c.CMCBranch,
		MCB_BRANCH_SOURCE:        c.MCBBranchSource,
		CONFIG_BRANCH:            c.ConfigBranch,
		MC_APP_COLLECTION_BRANCH: c.MCAppCollectionBranch,
		BASE_DOMAIN:              c.BaseDomain,
		CATALOG_REGISTRY_VALUES:  getCustomCatalogRegistryValues(c.RegistryDomain),
	}

	for k, t := range templates {

		containsVar, manifest, err := formatTemplate(t)
		if err != nil {
			return nil, fmt.Errorf("failed to format template %s.\n%w", k, err)
		}
		if !containsVar {
			continue
		}
		manifest, err = template.Execute(manifest, tmp)
		if err != nil {
			return nil, fmt.Errorf("failed to execute template %s.\n%w", k, err)
		}
		manifests[k] = manifest
	}
	return manifests, nil
}

func GetBaseConfig(templates map[string]string, path string) (Config, error) {
	log.Debug().Msg("Getting base config")

	baseDomain, err := getBaseDomain(templates[fmt.Sprintf("%s/%s", path, catalogPatchFile)])
	if err != nil {
		return Config{}, fmt.Errorf("failed to get base domain.\n%w", err)
	}
	registryDomain := getRegistryDomain(templates[fmt.Sprintf("%s/%s", path, catalogPatchFile)])

	cmcBranch := getBranch(templates[fmt.Sprintf("%s/%s", path, customBranchCMCFile)])
	mcAppCollectionBranch := getBranch(templates[fmt.Sprintf("%s/%s", path, customBranchCollectionFile)])
	configBranch := getBranch(templates[fmt.Sprintf("%s/%s", path, customBranchConfigFile)])

	mcbBranchSource := getMCBBranchSource(templates[fmt.Sprintf("%s/%s", path, catalogKustomizationFile)])

	return Config{
		CMCBranch:             cmcBranch,
		MCBBranchSource:       mcbBranchSource,
		ConfigBranch:          configBranch,
		BaseDomain:            baseDomain,
		RegistryDomain:        registryDomain,
		MCAppCollectionBranch: mcAppCollectionBranch,
	}, nil
}

func getBaseDomain(file string) (string, error) {
	return key.GetValue(baseDomainKey, file)
}

func getBranch(file string) string {
	b, err := key.GetValue(branchKey, file)
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("Failed to get custom branch from file.\n%s", err))
		return key.CMCMainBranch
	}
	return b
}

func getRegistryDomain(file string) string {
	d, err := key.GetValue(domainKey, file)
	if err != nil {
		log.Debug().Msg(fmt.Sprintf("Failed to get registry domain from file.%s. Assuming default value is used.", err))
		return ""
	}
	return d
}

func getMCBBranchSource(file string) string {

	k := `https://github.com/giantswarm/management-cluster-bases//bases/catalogs\?ref`

	/*we will use regex here to get the branch name from the reference k
		for example in the following file we would want to get the branch name "hello"
		apiVersion: kustomize.config.k8s.io/v1beta1
	kind: Kustomization
	patches:
	  - path: patches/appcatalog-default-patch.yaml
	resources:
	  - https://github.com/giantswarm/management-cluster-bases//bases/catalogs?ref=hello
	*/

	re := regexp.MustCompile(fmt.Sprintf(`%s=(.*)`, k))
	match := re.FindStringSubmatch(file)
	if len(match) < 2 {
		log.Debug().Msg(fmt.Sprintf("Failed to get MCB branch source from file.\n%s", match))
		return key.CMCMainBranch
	}
	return match[1]
}

func getCustomCatalogRegistryValues(domain string) string {
	if domain == "" {
		return ""
	}
	return fmt.Sprintf(`            image:
              registry: %s
            registry:
              domain: %s`, domain, domain)
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
