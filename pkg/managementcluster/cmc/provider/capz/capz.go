package capz

import (
	"fmt"

	"github.com/giantswarm/mcli/pkg/key"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
)

const (
	SecretName      = "cluster-identity-secret-static"
	ClientSecretKey = "clientSecret"
	ClientIDKey     = "clientID"
	TenantIDKey     = "tenantID"
	UAClientIDKey   = "clientID"
	UATenantIDKey   = "tenantID"
	UAResourceIDKey = "resourceID"
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

func GetCAPZConfig(sp string, ua string, secret string) (Config, error) {
	log.Debug().Msg("Getting CAPZ config")
	clientSecret, err := key.GetSecretValue(ClientSecretKey, secret)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get CAPZ client secret.\n%w", err)
	}
	clientID, err := key.GetValue(ClientIDKey, sp)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get CAPZ client ID.\n%w", err)
	}
	tenantID, err := key.GetValue(TenantIDKey, sp)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get CAPZ tenant ID.\n%w", err)
	}
	uaClientID, err := key.GetValue(UAClientIDKey, ua)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get CAPZ UA client ID.\n%w", err)
	}
	uaTenantID, err := key.GetValue(UATenantIDKey, ua)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get CAPZ UA tenant ID.\n%w", err)
	}
	uaResourceID, err := key.GetValue(UAResourceIDKey, ua)
	if err != nil {
		return Config{}, fmt.Errorf("failed to get CAPZ UA resource ID.\n%w", err)
	}

	return Config{
		UAClientID:   uaClientID,
		UATenantID:   uaTenantID,
		UAResourceID: uaResourceID,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TenantID:     tenantID,
	}, nil
}

func GetCAPZSecret(c Config) (string, error) {
	log.Debug().Msg("Creating CAPZ static SP file")
	secret := v1.Secret{
		Type: v1.SecretTypeOpaque,
		ObjectMeta: metav1.ObjectMeta{
			Name:      SecretName,
			Namespace: c.Namespace,
			Labels: map[string]string{
				"clusterctl.cluster.x-k8s.io/move": "true",
			},
		},
		Data: map[string][]byte{
			ClientSecretKey: []byte(c.ClientSecret),
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
				Name:      SecretName,
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
