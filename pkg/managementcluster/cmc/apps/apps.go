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
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
)

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
	appCROutput, err := GetAppsFile(c)
	if err != nil {
		return "", fmt.Errorf("failed to get app CRs for %s.\n%w", c.Name, err)
	}
	/* get the App object - it'd be better to return it initially, but kgs doesn't export a function for this
	// - let's change that
	app := &applicationv1alpha1.App{}
	if err = yaml.Unmarshal([]byte(appCROutput.AppCR), app); err != nil {
		return "", fmt.Errorf("failed to unmarshal app CR for %s.\n%w", c.Name, err)
	}

	if c.ConfigureContainerRegistries {
		log.Debug().Msg(fmt.Sprintf("Configuring container registries for %s", c.Name))
		app.Spec.ExtraConfigs = append(app.Spec.ExtraConfigs, applicationv1alpha1.AppExtraConfig{
			Kind:      "secret",
			Name:      ContainerRegistrySecretName,
			Namespace: "default",
		})
	}
	if key.IsProviderVsphere(c.Provider) {
		log.Debug().Msg(fmt.Sprintf("Configuring vsphere credentials for %s", c.Name))
		app.Spec.UserConfig.Secret.Name = "vsphere-credentials"
		app.Spec.UserConfig.Secret.Namespace = "org-giantswarm"
	}

	// update the app CR - see reasoning in above comment
	data, err := yaml.Marshal(app)
	if err != nil {
		return "", fmt.Errorf("failed to marshal app CR for %s.\n%w", c.Name, err)
	}
	appCROutput.AppCR = string(data)
	*/

	return TemplateApp(appCROutput)
}

func GetDefaultAppsFile(c Config) (string, error) {
	log.Debug().Msg(fmt.Sprintf("Getting the default apps config for %s", c.Cluster))
	appCROutput, err := GetAppsFile(c)
	if err != nil {
		return "", fmt.Errorf("failed to get app CRs for %s.\n%w", c.Name, err)
	}

	/* TODO: change this to not use the App object directly
	app := &applicationv1alpha1.App{}
	if err = yaml.Unmarshal([]byte(appCROutput.AppCR), app); err != nil {
		return "", fmt.Errorf("failed to unmarshal app CR for %s.\n%w", c.Name, err)
	}

	// Pass cluster-values configmap to default apps app.
	app.Spec.Config.ConfigMap.Name = fmt.Sprintf("%s-cluster-values", c.Cluster)
	app.Spec.Config.ConfigMap.Namespace = c.Namespace

	// The app-admission-controller will prevent the creation of the default apps `App` CR because the cluster-values configmap does not exist yet.
	if app.Labels == nil {
		app.Labels = map[string]string{}
	}
	app.Labels[label.Cluster] = c.Cluster
	app.Labels[label.ManagedBy] = "cluster"

	data, err := yaml.Marshal(app)
	if err != nil {
		return "", fmt.Errorf("failed to marshal app CR for %s.\n%w", c.Name, err)
	}
	appCROutput.AppCR = string(data)
	*/

	return TemplateApp(appCROutput)
}

func GetAppsFile(c Config) (*templateapp.AppCROutput, error) {
	var userConfigConfigMapYaml []byte
	var userConfigSecretYaml []byte
	var err error

	log.Debug().Msg(fmt.Sprintf("Creating the App for %s/%s-cluster", c.Name, c.Cluster))

	userConfigMap, err := GetUserConfigMap(c)
	if err != nil {
		return nil, err
	}
	appCRYaml, err := GetAppCRYaml(c, userConfigMap.GetName())
	if err != nil {
		return nil, err
	}
	userConfigConfigMapYaml, err = yaml.Marshal(userConfigMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal user config map for %s.\n%w", c.Name, err)
	}

	return &templateapp.AppCROutput{
		AppCR:               string(appCRYaml),
		UserConfigConfigMap: string(userConfigConfigMapYaml),
		UserConfigSecret:    string(userConfigSecretYaml),
	}, nil
}

func GetUserConfigMap(c Config) (*v1.ConfigMap, error) {
	log.Debug().Msg(fmt.Sprintf("Creating the user config map for %s/%s-cluster", c.Name, c.Cluster))
	configMapConfig := templateapp.UserConfig{
		Data:      c.Values,
		Name:      c.Cluster + "-user-values",
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
		Name:                  name,
		AppName:               app.Spec.Name,
		Catalog:               app.Spec.Catalog,
		Version:               app.Spec.Version,
		Namespace:             namespace,
		Values:                userConfigConfigMap.Data[ValuesKey],
		MCAppsPreventDeletion: userConfigConfigMap.Labels[label.PreventDeletion] == "true",
		Provider:              getProvider(app.Spec.Name),
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

func GetAppCRYaml(c Config, configmapName string) ([]byte, error) {
	log.Debug().Msg(fmt.Sprintf("Creating the app CR for %s/%s-cluster", c.Name, c.Cluster))
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
	appCRYaml, err := templateapp.NewAppCR(appConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create app CR for %s.\n%w", c.Name, err)
	}
	return appCRYaml, nil
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
