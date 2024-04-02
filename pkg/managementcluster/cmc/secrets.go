package cmc

import (
	"encoding/base64"

	"github.com/giantswarm/mcli/pkg/key"
	"github.com/rs/zerolog/log"
)

const (
	Redacted = "REDACTED"
)

func (c *CMC) RedactSecrets() {
	log.Debug().Msg("Redacting secret values")
	c.AgeKey = Redacted
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
	c.AgeKey = base64.StdEncoding.EncodeToString([]byte(c.AgeKey))
	c.TaylorBotToken = base64.StdEncoding.EncodeToString([]byte(c.TaylorBotToken))
	c.SharedDeployKey.Identity = base64.StdEncoding.EncodeToString([]byte(c.SharedDeployKey.Identity))
	c.SharedDeployKey.Passphrase = base64.StdEncoding.EncodeToString([]byte(c.SharedDeployKey.Passphrase))
	c.SharedDeployKey.KnownHosts = base64.StdEncoding.EncodeToString([]byte(c.SharedDeployKey.KnownHosts))
	c.CustomerDeployKey.Identity = base64.StdEncoding.EncodeToString([]byte(c.CustomerDeployKey.Identity))
	c.CustomerDeployKey.Passphrase = base64.StdEncoding.EncodeToString([]byte(c.CustomerDeployKey.Passphrase))
	c.CustomerDeployKey.KnownHosts = base64.StdEncoding.EncodeToString([]byte(c.CustomerDeployKey.KnownHosts))
	c.SSHdeployKey.Identity = base64.StdEncoding.EncodeToString([]byte(c.SSHdeployKey.Identity))
	c.SSHdeployKey.Passphrase = base64.StdEncoding.EncodeToString([]byte(c.SSHdeployKey.Passphrase))
	c.SSHdeployKey.KnownHosts = base64.StdEncoding.EncodeToString([]byte(c.SSHdeployKey.KnownHosts))

	if c.ConfigureContainerRegistries.Enabled {
		c.ConfigureContainerRegistries.Values = base64.StdEncoding.EncodeToString([]byte(c.ConfigureContainerRegistries.Values))
	}

	if key.IsProviderAzure(c.Provider.Name) {
		c.Provider.CAPZ.ClientSecret = base64.StdEncoding.EncodeToString([]byte(c.Provider.CAPZ.ClientSecret))
	} else if key.IsProviderVCD(c.Provider.Name) {
		c.Provider.CAPVCD.RefreshToken = base64.StdEncoding.EncodeToString([]byte(c.Provider.CAPVCD.RefreshToken))
	} else if key.IsProviderVsphere(c.Provider.Name) {
		c.Provider.CAPV.CloudConfig = base64.StdEncoding.EncodeToString([]byte(c.Provider.CAPV.CloudConfig))
	}
}
