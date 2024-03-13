package certmanager

import (
	"fmt"

	"github.com/giantswarm/mcli/pkg/key"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	var secret v1.Secret
	err := yaml.Unmarshal([]byte(file), &secret)
	if err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal Route53 configuration secret.\n%w", err)
	}
	values := secret.Data["values"]
	var cmSecret CMSecret
	err = yaml.Unmarshal(values, &cmSecret)
	if err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal cm secret object.\n%w", err)
	}

	return Config{
		Region:          cmSecret.GiantSwarmClusterIssuer.Acme.Dns01.Route53["region"],
		Role:            cmSecret.GiantSwarmClusterIssuer.Acme.Dns01.Route53["role"],
		AccessKeyID:     cmSecret.GiantSwarmClusterIssuer.Acme.Dns01.Route53["accessKeyID"],
		SecretAccessKey: cmSecret.GiantSwarmClusterIssuer.Acme.Dns01.Route53["secretAccessKey"],
	}, nil
}

func GetCertManagerFile(c Config) (string, error) {
	log.Debug().Msg("Creating Route53 configuration file for cert-manager")

	cmSecret := CMSecret{
		GiantSwarmClusterIssuer: GiantSwarmClusterIssuer{
			Acme: Acme{
				Dns01: Dns01{
					Route53: map[string]string{
						"region":          c.Region,
						"role":            c.Role,
						"accessKeyID":     c.AccessKeyID,
						"secretAccessKey": c.SecretAccessKey,
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
