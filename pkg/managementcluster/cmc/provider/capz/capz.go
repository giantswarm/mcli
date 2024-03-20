package capz

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
)

const (
	secretName = "cluster-identity-secret-static"
)

type Config struct {
	Namespace    string
	UAClientID   string
	UATenantID   string
	UAResourceID string
	ClientID     string
	ClientSecret string
	TenantID     string
}

func GetCAPZConfig(sp string, ua string, staticsp string) (Config, error) {
	secret := v1.Secret{}
	err := yaml.Unmarshal([]byte(sp), &secret)
	if err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal static-sp object.\n%w", err)
	}
	uaCR := v1beta1.AzureClusterIdentity{}
	err = yaml.Unmarshal([]byte(ua), &uaCR)
	if err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal UA object.\n%w", err)
	}
	spCR := v1beta1.AzureClusterIdentity{}
	err = yaml.Unmarshal([]byte(staticsp), &spCR)
	if err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal static-sp object.\n%w", err)
	}
	return Config{
		Namespace:    uaCR.Namespace,
		UAClientID:   uaCR.Spec.ClientID,
		UATenantID:   uaCR.Spec.TenantID,
		UAResourceID: uaCR.Spec.ResourceID,
		ClientID:     spCR.Spec.ClientID,
		ClientSecret: string(secret.Data["clientSecret"]),
		TenantID:     spCR.Spec.TenantID,
	}, nil
}

func GetCAPZSecret(c Config) (string, error) {
	log.Debug().Msg("Creating CAPZ static SP file")
	secret := v1.Secret{
		Type: v1.SecretTypeOpaque,
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: c.Namespace,
			Labels: map[string]string{
				"clusterctl.cluster.x-k8s.io/move": "true",
			},
		},
		Data: map[string][]byte{
			"clientSecret": []byte(c.ClientSecret),
		},
	}
	data, err := yaml.Marshal(secret)
	if err != nil {
		return "", fmt.Errorf("failed to marshal static-sp object.\n%w", err)
	}
	return string(data), nil
}

func GetCAPZSPFile(c Config) (string, error) {
	log.Debug().Msg("Creating CAPZ SP file")
	sp := v1beta1.AzureClusterIdentity{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "infrastructure.cluster.x-k8s.io/v1beta1",
			Kind:       "AzureClusterIdentity",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cluster-identity-static-sp",
			Namespace: c.Namespace,
			Labels: map[string]string{
				"clusterctl.cluster.x-k8s.io/move": "true",
			},
		},
		Spec: v1beta1.AzureClusterIdentitySpec{
			AllowedNamespaces: &v1beta1.AllowedNamespaces{},
			ClientID:          c.ClientID,
			ClientSecret: v1.SecretReference{
				Name:      secretName,
				Namespace: "org-giantswarm",
			},
			TenantID: c.TenantID,
			Type:     "ManualServicePrincipal",
		},
	}
	data, err := yaml.Marshal(sp)
	if err != nil {
		return "", fmt.Errorf("failed to marshal static-sp object.\n%w", err)
	}
	return string(data), nil
}

func GetCAPZUAFile(c Config) (string, error) {
	log.Debug().Msg("Creating CAPZ UA file")
	ua := v1beta1.AzureClusterIdentity{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "infrastructure.cluster.x-k8s.io/v1beta1",
			Kind:       "AzureClusterIdentity",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cluster-identity",
			Namespace: c.Namespace,
		},
		Spec: v1beta1.AzureClusterIdentitySpec{
			AllowedNamespaces: &v1beta1.AllowedNamespaces{},
			ClientID:          c.UAClientID,
			TenantID:          c.UATenantID,
			ResourceID:        c.UAResourceID,
			Type:              "UserAssignedMSI",
		},
	}
	data, err := yaml.Marshal(ua)
	if err != nil {
		return "", fmt.Errorf("failed to marshal UA object.\n%w", err)
	}
	return string(data), nil
}
