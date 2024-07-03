package certmanager

import (
	"encoding/base64"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/giantswarm/mcli/pkg/key"
	"github.com/giantswarm/mcli/pkg/template"
)

const (
	ValuesKey          = "values"
	RegionKey          = "region"
	RoleKey            = "role"
	AccessKeyIDKey     = "accessKeyID"
	SecretAccessKeyKey = "secretAccessKey"
)

const CertManagerTemplate = `apiVersion: v1
kind: Secret
metadata:
  name: {{ .Cluster }}-cert-manager-user-secrets
  namespace: {{ .ClusterNamespace }}
data:
  values: {{ .Values }}
`

type Config struct {
	Cluster          string
	ClusterNamespace string
	Region           string
	Role             string
	AccessKeyID      string
	SecretAccessKey  string
	MCProxyEnabled   bool
}

type CMSecret struct {
	Global                        map[string]string       `yaml:"global"`
	GiantSwarmClusterIssuer       GiantSwarmClusterIssuer `yaml:"giantSwarmClusterIssuer"`
	DNS01RecursiveNameserversOnly bool                    `yaml:"dns01RecursiveNameserversOnly,omitempty"`
	DNS01RecursiveNameservers     []string                `yaml:"dns01RecursiveNameservers,omitempty"`
}
type GiantSwarmClusterIssuer struct {
	Acme Acme `yaml:"acme"`
}
type Acme struct {
	Dns01 Dns01 `yaml:"dns01"`
}
type Dns01 struct {
	Route53 map[string]string `yaml:"route53"`
}

func GetCertManagerConfig(file string) (Config, error) {
	log.Debug().Msg("Getting Route53 configuration")

	values, err := key.GetSecretValue(ValuesKey, file)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get values.\n%w", err)
	}

	region, err := key.GetValue(RegionKey, values)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get region.\n%w", err)
	}

	role, err := key.GetValue(RoleKey, values)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get role.\n%w", err)
	}

	accessKeyID, err := key.GetValue(AccessKeyIDKey, values)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get accessKeyID.\n%w", err)
	}

	secretAccessKey, err := key.GetValue(SecretAccessKeyKey, values)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get secretAccessKey.\n%w", err)
	}

	return Config{
		Region:          region,
		Role:            role,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
	}, nil
}

func GetCertManagerFile(c Config) (string, error) {
	log.Debug().Msg("Creating Route53 configuration file for cert-manager")

	cmSecret := CMSecret{
		GiantSwarmClusterIssuer: GiantSwarmClusterIssuer{
			Acme: Acme{
				Dns01: Dns01{
					Route53: map[string]string{
						RegionKey:          c.Region,
						RoleKey:            c.Role,
						AccessKeyIDKey:     c.AccessKeyID,
						SecretAccessKeyKey: c.SecretAccessKey,
					},
				},
			},
		},
	}
	if c.MCProxyEnabled {
		cmSecret.DNS01RecursiveNameserversOnly = true
		cmSecret.DNS01RecursiveNameservers = []string{"$(COREDNS_SERVICE_HOST):53"}
	}
	data, err := key.GetData(cmSecret)
	if err != nil {
		return "", fmt.Errorf("failed to marshal cm secret object.\n%w", err)
	}

	type config struct {
		Cluster          string
		ClusterNamespace string
		Values           string
	}

	return template.Execute(CertManagerTemplate, config{
		Cluster:          c.Cluster,
		ClusterNamespace: c.ClusterNamespace,
		Values:           base64.StdEncoding.EncodeToString(data),
	})
}

func GetCertManagerDefaultAppConfigMap() string {
	log.Debug().Msg("Creating cert-manager user values config map")
	return `kind: ConfigMap
apiVersion: v1
metadata:
  annotations:
    meta.helm.sh/release-name: cert-manager
    meta.helm.sh/release-namespace: kube-system
  labels:
    app.kubernetes.io/managed-by: Helm
    k8s-app: cert-manager
  name: cert-manager-user-values
  namespace: kube-system
data:
  values: |
    controller:
      defaultIssuer:
        name: private-giantswarm`
}
