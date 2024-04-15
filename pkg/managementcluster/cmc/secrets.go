package cmc

import (
	"encoding/base64"
	"fmt"
	"reflect"

	"github.com/rs/zerolog/log"

	"github.com/giantswarm/mcli/pkg/github"
	"github.com/giantswarm/mcli/pkg/key"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc/kustomization"
)

const (
	Redacted = "REDACTED"
)

func (c *CMC) RedactSecrets() {
	log.Debug().Msg("Redacting secret values")
	c.TaylorBotToken = Redacted
	c.SharedDeployKey.Identity = Redacted
	c.SharedDeployKey.Passphrase = Redacted
	c.CustomerDeployKey.Identity = Redacted
	c.CustomerDeployKey.Passphrase = Redacted
	c.SSHdeployKey.Identity = Redacted
	c.SSHdeployKey.Passphrase = Redacted

	if c.CertManagerDNSChallenge.Enabled {
		c.CertManagerDNSChallenge.SecretAccessKey = Redacted
	}
	if c.ConfigureContainerRegistries.Enabled {
		c.ConfigureContainerRegistries.Values = Redacted
	}

	if key.IsProviderAzure(c.Provider.Name) {
		c.Provider.CAPZ.ClientSecret = Redacted
	} else if key.IsProviderVCD(c.Provider.Name) {
		c.Provider.CAPVCD.RefreshToken = Redacted
	} else if key.IsProviderVsphere(c.Provider.Name) {
		c.Provider.CAPV.CloudConfig = Redacted
	}
}

func (c *CMC) EncodeSecrets() {
	log.Debug().Msg("Base64 encoding secret values")
	c.TaylorBotToken = encodeSecret(c.TaylorBotToken)
	c.SharedDeployKey.Identity = encodeSecret(c.SharedDeployKey.Identity)
	c.SharedDeployKey.Passphrase = encodeSecret(c.SharedDeployKey.Passphrase)
	c.SharedDeployKey.KnownHosts = encodeSecret(c.SharedDeployKey.KnownHosts)
	c.CustomerDeployKey.Identity = encodeSecret(c.CustomerDeployKey.Identity)
	c.CustomerDeployKey.Passphrase = encodeSecret(c.CustomerDeployKey.Passphrase)
	c.CustomerDeployKey.KnownHosts = encodeSecret(c.CustomerDeployKey.KnownHosts)
	c.SSHdeployKey.Identity = encodeSecret(c.SSHdeployKey.Identity)
	c.SSHdeployKey.Passphrase = encodeSecret(c.SSHdeployKey.Passphrase)
	c.SSHdeployKey.KnownHosts = encodeSecret(c.SSHdeployKey.KnownHosts)
	c.ConfigureContainerRegistries.Values = encodeSecret(c.ConfigureContainerRegistries.Values)
	c.Provider.CAPZ.ClientSecret = encodeSecret(c.Provider.CAPZ.ClientSecret)
	c.Provider.CAPVCD.RefreshToken = encodeSecret(c.Provider.CAPVCD.RefreshToken)
	c.Provider.CAPV.CloudConfig = encodeSecret(c.Provider.CAPV.CloudConfig)
}

func (c *CMC) DecodeSecrets() error {
	log.Debug().Msg("Base64 decoding secret values")
	var err error
	c.TaylorBotToken, err = decodeSecret(c.TaylorBotToken)
	if err != nil {
		return fmt.Errorf("failed to decode TaylorBotToken: %w", err)
	}
	c.SharedDeployKey.Identity, err = decodeSecret(c.SharedDeployKey.Identity)
	if err != nil {
		return fmt.Errorf("failed to decode SharedDeployKey Identity: %w", err)
	}
	c.SharedDeployKey.Passphrase, err = decodeSecret(c.SharedDeployKey.Passphrase)
	if err != nil {
		return fmt.Errorf("failed to decode SharedDeployKey Passphrase: %w", err)
	}
	c.SharedDeployKey.KnownHosts, err = decodeSecret(c.SharedDeployKey.KnownHosts)
	if err != nil {
		return fmt.Errorf("failed to decode SharedDeployKey KnownHosts: %w", err)
	}
	c.CustomerDeployKey.Identity, err = decodeSecret(c.CustomerDeployKey.Identity)
	if err != nil {
		return fmt.Errorf("failed to decode CustomerDeployKey Identity: %w", err)
	}
	c.CustomerDeployKey.Passphrase, err = decodeSecret(c.CustomerDeployKey.Passphrase)
	if err != nil {
		return fmt.Errorf("failed to decode CustomerDeployKey Passphrase: %w", err)
	}
	c.CustomerDeployKey.KnownHosts, err = decodeSecret(c.CustomerDeployKey.KnownHosts)
	if err != nil {
		return fmt.Errorf("failed to decode CustomerDeployKey KnownHosts: %w", err)
	}
	c.SSHdeployKey.Identity, err = decodeSecret(c.SSHdeployKey.Identity)
	if err != nil {
		return fmt.Errorf("failed to decode SSHdeployKey Identity: %w", err)
	}
	c.SSHdeployKey.Passphrase, err = decodeSecret(c.SSHdeployKey.Passphrase)
	if err != nil {
		return fmt.Errorf("failed to decode SSHdeployKey Passphrase: %w", err)
	}
	c.SSHdeployKey.KnownHosts, err = decodeSecret(c.SSHdeployKey.KnownHosts)
	if err != nil {
		return fmt.Errorf("failed to decode SSHdeployKey KnownHosts: %w", err)
	}
	c.ConfigureContainerRegistries.Values, err = decodeSecret(c.ConfigureContainerRegistries.Values)
	if err != nil {
		return fmt.Errorf("failed to decode ConfigureContainerRegistries Values: %w", err)
	}
	c.Provider.CAPZ.ClientSecret, err = decodeSecret(c.Provider.CAPZ.ClientSecret)
	if err != nil {
		return fmt.Errorf("failed to decode CAPZ ClientSecret: %w", err)
	}
	c.Provider.CAPVCD.RefreshToken, err = decodeSecret(c.Provider.CAPVCD.RefreshToken)
	if err != nil {
		return fmt.Errorf("failed to decode CAPVCD RefreshToken: %w", err)
	}
	c.Provider.CAPV.CloudConfig, err = decodeSecret(c.Provider.CAPV.CloudConfig)
	if err != nil {
		return fmt.Errorf("failed to decode CAPV CloudConfig: %w", err)
	}
	return nil
}

func encodeSecret(secret string) string {
	if secret == "" {
		return ""
	}
	return base64.StdEncoding.EncodeToString([]byte(secret))
}

func decodeSecret(secret string) (string, error) {
	if secret == "" {
		return "", nil
	}
	decoded, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to decode secret: %w", err)
	}
	return string(decoded), nil
}

// the reason for this is that when secrets are encrypted, they will always appear to be different
// so we need to remove secrets that have not changed from the update map
func MarkUnchangedSecretsInMap(currentCMC *CMC, desiredCMC *CMC, update map[string]string) (map[string]string, error) {
	if reflect.DeepEqual(currentCMC.SSHdeployKey, desiredCMC.SSHdeployKey) {
		markUnchanged(update, fmt.Sprintf("%s/%s", key.GetCMCPath(desiredCMC.Cluster), kustomization.SSHdeployKeyFile))
	}
	if reflect.DeepEqual(currentCMC.CustomerDeployKey, desiredCMC.CustomerDeployKey) {
		markUnchanged(update, fmt.Sprintf("%s/%s", key.GetCMCPath(desiredCMC.Cluster), kustomization.CustomerDeployKeyFile))
	}
	if reflect.DeepEqual(currentCMC.SharedDeployKey, desiredCMC.SharedDeployKey) {
		markUnchanged(update, fmt.Sprintf("%s/%s", key.GetCMCPath(desiredCMC.Cluster), kustomization.SharedDeployKeyFile))
	}
	if currentCMC.TaylorBotToken == desiredCMC.TaylorBotToken {
		markUnchanged(update, fmt.Sprintf("%s/%s", key.GetCMCPath(desiredCMC.Cluster), kustomization.TaylorBotFile))
	}
	if reflect.DeepEqual(currentCMC.CertManagerDNSChallenge, desiredCMC.CertManagerDNSChallenge) {
		markUnchanged(update, fmt.Sprintf("%s/%s", key.GetCMCPath(desiredCMC.Cluster), kustomization.CertManagerFile))
	}
	if reflect.DeepEqual(currentCMC.ConfigureContainerRegistries, desiredCMC.ConfigureContainerRegistries) {
		markUnchanged(update, fmt.Sprintf("%s/%s", key.GetCMCPath(desiredCMC.Cluster), kustomization.RegistryFile))
	}
	if reflect.DeepEqual(currentCMC.Provider.CAPV, desiredCMC.Provider.CAPV) {
		markUnchanged(update, fmt.Sprintf("%s/%s", key.GetCMCPath(desiredCMC.Cluster), kustomization.VsphereCredentialsFile))
	}
	if reflect.DeepEqual(currentCMC.Provider.CAPVCD, desiredCMC.Provider.CAPVCD) {
		markUnchanged(update, fmt.Sprintf("%s/%s", key.GetCMCPath(desiredCMC.Cluster), kustomization.CloudDirectorCredentialsFile))
	}
	if reflect.DeepEqual(currentCMC.Provider.CAPZ, desiredCMC.Provider.CAPZ) {
		markUnchanged(update, fmt.Sprintf("%s/%s", key.GetCMCPath(desiredCMC.Cluster), kustomization.AzureSecretClusterIdentityStaticSP))
	}
	return update, nil
}

func markUnchanged(update map[string]string, key string) {
	if _, ok := update[key]; ok {
		update[key] = fmt.Sprintf("%s\n%s", update[key], github.ActionNoChangesMarker)
	}
}
