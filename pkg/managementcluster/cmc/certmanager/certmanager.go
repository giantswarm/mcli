package certmanager

import (
	"fmt"

	"github.com/giantswarm/mcli/pkg/key"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ValuesKey          = "values"
	RegionKey          = "region"
	RoleKey            = "role"
	AccessKeyIDKey     = "accessKeyID"
	SecretAccessKeyKey = "secretAccessKey"
)

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
	data, err := yaml.Marshal(cmSecret)
	if err != nil {
		return "", fmt.Errorf("failed to marshal cm secret object.\n%w", err)
	}

	secret := v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      key.GetCertManagerSecretName(c.Cluster),
			Namespace: c.ClusterNamespace,
		},
		Data: map[string][]byte{
			"values": []byte(data),
		},
	}
	data, err = yaml.Marshal(secret)
	if err != nil {
		return "", fmt.Errorf("failed to marshal cert manager configuration secret.\n%w", err)
	}
	return string(data), nil
}
