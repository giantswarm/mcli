package capz

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/giantswarm/mcli/pkg/key"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc/kustomization"
	"github.com/giantswarm/mcli/pkg/template"
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

	return template.Execute(key.GetTMPLFile(kustomization.AzureSecretClusterIdentityStaticSP), c)
}

func GetCAPZSPFile(c Config) (string, error) {
	log.Debug().Msg("Creating CAPZ SP file")

	return template.Execute(key.GetTMPLFile(kustomization.AzureClusterIdentitySPFile), c)
}

func GetCAPZUAFile(c Config) (string, error) {
	log.Debug().Msg("Creating CAPZ UA file")

	return template.Execute(key.GetTMPLFile(kustomization.AzureClusterIdentityUAFile), c)
}
