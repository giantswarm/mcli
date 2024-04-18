package apps

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	applicationv1alpha1 "github.com/giantswarm/apiextensions-application/api/v1alpha1"
	templateapp "github.com/giantswarm/kubectl-gs/v2/pkg/template/app"
	"github.com/giantswarm/mcli/pkg/key"

	"github.com/giantswarm/k8smetadata/pkg/label"
	"github.com/rs/zerolog/log"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

// The app templating in mc-bootstrap is done through kubectl gs
// That is why the code here replicates kubectl gs (not all packages are exported)
// I don't think it should stay this way but for now we want to be as close as possible to the original

const (
	ContainerRegistrySecretName = "container-registries-configuration"
	ValuesKey                   = "values"
)

type Config struct {
	Cluster                      string
	Name                         string
	AppName                      string
	Catalog                      string
	Version                      string
	Namespace                    string
	Values                       string
	ConfigureContainerRegistries bool
	Provider                     string
	MCAppsPreventDeletion        bool
}

func GetClusterAppsFile(c Config) (string, error) {
	log.Debug().Msg(fmt.Sprintf("Getting the cluster apps config for %s", c.Cluster))

	userConfigMap, err := GetUserConfigMap(c)
	if err != nil {
		return "", err
	}

	appsConfig, err := getClusterAppConfig(c, userConfigMap.GetName())
	if err != nil {
		return "", fmt.Errorf("failed to get app CRs for %s.\n%w", c.Name, err)
	}

	appCROutput, err := GetAppsFile(appsConfig, userConfigMap)
	if err != nil {
		return "", fmt.Errorf("failed to get app CRs for %s.\n%w", c.Name, err)
	}

	return TemplateApp(appCROutput)
}

func GetDefaultAppsFile(c Config) (string, error) {
	log.Debug().Msg(fmt.Sprintf("Getting the default apps config for %s", c.Cluster))
	userConfigMap, err := GetUserConfigMap(c)
	if err != nil {
		return "", err
	}

	appsConfig, err := getDefaultAppConfig(c, userConfigMap.GetName())
	if err != nil {
		return "", fmt.Errorf("failed to get app CRs for %s.\n%w", c.Name, err)
	}

	appCROutput, err := GetAppsFile(appsConfig, userConfigMap)
	if err != nil {
		return "", fmt.Errorf("failed to get app CRs for %s.\n%w", c.Name, err)
	}

	appCROutput.AppCR = addClusterConfig(appCROutput.AppCR, c)

	return TemplateApp(appCROutput)
}

func getClusterAppConfig(c Config, configmapName string) (templateapp.Config, error) {

	log.Debug().Msg(fmt.Sprintf("Creating the cluster app CR for %s/%s-cluster", c.Name, c.Cluster))
	appConfig := templateapp.Config{
		AppName:                 c.Name, //todo: this needs to be swapped
		Catalog:                 c.Catalog,
		Cluster:                 c.Cluster,
		DefaultingEnabled:       true,
		ExtraLabels:             map[string]string{},
		InCluster:               true,
		Name:                    c.AppName,
		Namespace:               c.Namespace,
		Version:                 c.Version,
		UserConfigConfigMapName: configmapName,
	}
	if c.MCAppsPreventDeletion {
		appConfig.ExtraLabels[label.PreventDeletion] = "true"
	}
	if c.ConfigureContainerRegistries {
		log.Debug().Msg(fmt.Sprintf("Configuring container registries for %s", c.Name))
		appConfig.ExtraConfigs = []applicationv1alpha1.AppExtraConfig{
			{
				Kind:      "secret",
				Name:      ContainerRegistrySecretName,
				Namespace: "default",
			},
		}
	}
	if key.IsProviderVsphere(c.Provider) {
		log.Debug().Msg(fmt.Sprintf("Configuring vsphere credentials for %s", c.Name))
		appConfig.UserConfigSecretName = "vsphere-credentials"
	}
	return appConfig, nil
}

func getDefaultAppConfig(c Config, configmapName string) (templateapp.Config, error) {

	log.Debug().Msg(fmt.Sprintf("Creating the app CR for %s/%s-cluster", c.Name, c.Cluster))
	appConfig := templateapp.Config{
		AppName:           c.Name, //todo: this needs to be swapped
		Catalog:           c.Catalog,
		Cluster:           c.Cluster,
		DefaultingEnabled: true,
		ExtraLabels: map[string]string{
			label.Cluster:   c.Cluster,
			label.ManagedBy: "cluster",
		},
		InCluster:               true,
		Name:                    c.AppName,
		Namespace:               c.Namespace,
		Version:                 c.Version,
		UserConfigConfigMapName: configmapName,
	}
	if c.MCAppsPreventDeletion {
		appConfig.ExtraLabels[label.PreventDeletion] = "true"
	}
	return appConfig, nil
}

func GetUserConfigMap(c Config) (*v1.ConfigMap, error) {
	log.Debug().Msg(fmt.Sprintf("Creating the user config map for %s/%s-cluster", c.Name, c.Cluster))
	configMapConfig := templateapp.UserConfig{
		Data:      c.Values,
		Name:      c.Name + "-user-values",
		Namespace: c.Namespace,
	}
	userConfigMap, err := templateapp.NewConfigMap(configMapConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating user config map for %s.\n%w", c.Name, err)
	}
	if c.MCAppsPreventDeletion {
		if userConfigMap.Labels == nil {
			userConfigMap.Labels = map[string]string{}
		}
		userConfigMap.Labels[label.PreventDeletion] = "true"
	}
	return userConfigMap, nil
}

func GetAppsConfig(file string) (Config, error) {
	var app applicationv1alpha1.App
	var userConfigConfigMap v1.ConfigMap
	var name, namespace string

	log.Debug().Msg("Getting the apps config")
	files := strings.Split(file, "---")
	for _, f := range files {
		if strings.Contains(f, "kind: ConfigMap") {
			if err := yaml.Unmarshal([]byte(f), &userConfigConfigMap); err != nil {
				return Config{}, fmt.Errorf("failed to unmarshal user config map for %s.\n%w", file, err)
			}
		} else if strings.Contains(f, "kind: App") {
			if err := yaml.Unmarshal([]byte(f), &app); err != nil {
				return Config{}, fmt.Errorf("failed to unmarshal app CR for %s.\n%w", file, err)
			}
			// get the name and namespace from the app CR - metadata is removed during unmarshalling
			// todo: make this a bit more elegant
			name, namespace = key.GetNamespacedName(f)
		}
	}
	return Config{
		Name:                         name,
		AppName:                      app.Spec.Name,
		Catalog:                      app.Spec.Catalog,
		Version:                      app.Spec.Version,
		Namespace:                    namespace,
		Values:                       userConfigConfigMap.Data[ValuesKey],
		MCAppsPreventDeletion:        userConfigConfigMap.Labels[label.PreventDeletion] == "true",
		ConfigureContainerRegistries: len(app.Spec.ExtraConfigs) > 0,
		Provider:                     getProvider(app.Spec.Name),
	}, nil
}

func getProvider(appName string) string {
	if strings.Contains(appName, "aws") {
		return key.ProviderAWS
	}
	if strings.Contains(appName, "azure") {
		return key.ProviderAzure
	}
	if strings.Contains(appName, key.ProviderVsphere) {
		return key.ProviderVsphere
	}
	if strings.Contains(appName, key.ProviderVCD) {
		return key.ProviderVCD
	}
	return ""
}

func GetAppsFile(config templateapp.Config, userConfigMap *v1.ConfigMap) (*templateapp.AppCROutput, error) {
	var userConfigConfigMapYaml []byte
	var userConfigSecretYaml []byte
	var err error

	log.Debug().Msg(fmt.Sprintf("Creating the App %s for %s", config.Name, config.Cluster))

	appCRYaml, err := templateapp.NewAppCR(config)
	if err != nil {
		return nil, err
	}
	userConfigConfigMapYaml, err = yaml.Marshal(userConfigMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user config map for %s.\n%w", config.Name, err)
	}

	return &templateapp.AppCROutput{
		AppCR:               string(appCRYaml),
		UserConfigConfigMap: string(userConfigConfigMapYaml),
		UserConfigSecret:    string(userConfigSecretYaml),
	}, nil
}

func TemplateApp(appCROutput *templateapp.AppCROutput) (string, error) {
	log.Debug().Msg("Templating the app CR")
	var b bytes.Buffer

	t := template.Must(template.New("appCR").Parse(AppCRTemplate))

	if err := t.Execute(&b, appCROutput); err != nil {
		return "", fmt.Errorf("failed to execute app CR template.\n%w", err)
	}

	return b.String(), nil
}

// oddly enough this is not taken care of by kgs
func addClusterConfig(appCR string, c Config) string {
	clusterConfig := fmt.Sprintf(`  config:
    configMap:
      name: %s-cluster-values
      namespace: %s
`, c.Cluster, c.Namespace)
	return appCR + clusterConfig
}

const AppCRTemplate = `
{{- if .UserConfigConfigMap -}}
---
{{ .UserConfigConfigMap -}}
{{- end -}}
{{- if .UserConfigSecret -}}
---
{{ .UserConfigSecret -}}
{{- end -}}
---
{{ .AppCR -}}
`
