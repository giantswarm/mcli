package capz

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func GetCAPZConfig(sp string, ua string, staticsp string) Config {
	return Config{}
}

func GetCAPZSecret(c Config) (string, error) {
	log.Debug().Msg("Creating CAPZ static SP file")
	secret := v1.Secret{
		Type: v1.SecretTypeOpaque,
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cluster-identity-secret-static-sp",
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

func GetCAPZSPFile(c Config) string {
	log.Debug().Msg("Creating CAPZ SP file")
	return getAzureClusterIdentityManualSP(c.Namespace, c.ClientID, c.TenantID)
}

func GetCAPZUAFile(c Config) string {
	log.Debug().Msg("Creating CAPZ UA file")
	return getAzureClusterIdentityUA(c.Namespace, c.UAClientID, c.UATenantID, c.UAResourceID)
}

func getAzureClusterIdentityUA(namespace string, clientID string, tenantID string, resourceID string) string {
	return fmt.Sprintf(`apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: AzureClusterIdentity
metadata:
  name: cluster-identity
  namespace: %s
spec:
  allowedNamespaces: {}
  clientID: %s
  tenantID: %s
  resourceID: %s
  type: UserAssignedMSI`, namespace, clientID, tenantID, resourceID)
}

func getAzureClusterIdentityManualSP(namespace string, clientID string, tenantID string) string {
	return fmt.Sprintf(`apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: AzureClusterIdentity
metadata:
  name: cluster-identity-static-sp
  namespace: %s
  labels:
	clusterctl.cluster.x-k8s.io/move: "true"
spec:
  allowedNamespaces: {}
  clientID: %s
  clientSecret:
	name: cluster-identity-secret-static-sp
	namespace: org-giantswarm
  tenantID: %s
  type: ManualServicePrincipal`, namespace, clientID, tenantID)
}
