/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/giantswarm/mcli/pkg/key"
)

const (
	flagInput        = "input"
	flagBaseDomain   = "base-domain"
	flagTeam         = "team"
	flagAWSRegion    = "aws-region"
	flagAWSAccountID = "aws-account-id"
)

const (
	envBaseDomain   = "BASE_DOMAIN"
	envTeam         = "TEAM_NAME"
	envAWSRegion    = "AWS_REGION"
	envAWSAccountID = "INSTALLATION_AWS_ACCOUNT"
)

/*
based on this
type CMCFlags struct {
	SecretFolder                 string
	MCAppsPreventDeletion        bool
	ClusterAppName               string
	ClusterAppCatalog            string
	ClusterAppVersion            string
	ClusterNamespace             string
	ConfigureContainerRegistries bool
	DefaultAppsName              string
	DefaultAppsCatalog           string
	DefaultAppsVersion           string
	PrivateCA                    string
	CertManagerDNSChallenge      bool
	MCCustomCoreDNSConfig        bool
	MCProxyEnabled               bool
	MCHTTPSProxy                 string
}
*/

const (
	flagSecretFolder                 = "secret-folder"
	flagMCAppsPreventDeletion        = "mc-apps-prevent-deletion"
	flagClusterAppName               = "cluster-app-name"
	flagClusterAppCatalog            = "cluster-app-catalog"
	flagClusterAppVersion            = "cluster-app-version"
	flagClusterNamespace             = "cluster-namespace"
	flagConfigureContainerRegistries = "configure-container-registries"
	flagDefaultAppsName              = "default-apps-name"
	flagDefaultAppsCatalog           = "default-apps-catalog"
	flagDefaultAppsVersion           = "default-apps-version"
	flagPrivateCA                    = "private-ca"
	flagCertManagerDNSChallenge      = "cert-manager-dns-challenge"
	flagMCCustomCoreDNSConfig        = "mc-custom-coredns-config"
	flagMCProxyEnabled               = "mc-proxy-enabled"
	flagMCHTTPSProxy                 = "mc-https-proxy"
)

const (
	envSecretFolder                 = "SECRETS_FOLDER"
	envMCAppsPreventDeletion        = "MC_APPS_PREVENT_DELETION"
	envClusterAppName               = "CLUSTER_APP_NAME"
	envClusterAppCatalog            = "CLUSTER_APP_CATALOG"
	envClusterAppVersion            = "CLUSTER_APP_VERSION"
	envClusterNamespace             = "CLUSTER_NAMESPACE"
	envConfigureContainerRegistries = "CONFIGURE_CONTAINER_REGISTRIES"
	envDefaultAppsName              = "DEFAULT_APPS_APP_NAME"
	envDefaultAppsCatalog           = "DEFAULT_APPS_APP_CATALOG"
	envDefaultAppsVersion           = "DEFAULT_APPS_APP_VERSION"
	envPrivateCA                    = "PRIVATE_CA"
	envCertManagerDNSChallenge      = "CERT_MANAGER_DNS01_CHALLENGE"
	envMCCustomCoreDNSConfig        = "MC_CUSTOM_COREDNS_CONFIG"
	envMCProxyEnabled               = "MC_PROXY_ENABLED"
	envMCHTTPSProxy                 = "MC_HTTPS_PROXY"
)

var (
	input                        string
	baseDomain                   string
	team                         string
	provider                     string
	awsRegion                    string
	awsAccountID                 string
	secretFolder                 string
	mcAppsPreventDeletion        bool
	clusterAppName               string
	clusterAppCatalog            string
	clusterAppVersion            string
	clusterNamespace             string
	configureContainerRegistries bool
	defaultAppsName              string
	defaultAppsCatalog           string
	defaultAppsVersion           string
	privateCA                    bool
	certManagerDNSChallenge      bool
	mcCustomCoreDNSConfig        bool
	mcProxyEnabled               bool
	mcHTTPSProxy                 string
)

func addFlagsPush() {
	pushCmd.Flags().StringArrayVarP(&skip, flagSkip, "s", []string{}, fmt.Sprintf("List of repositories to skip. (default: none) Valid values: %s", key.GetValidRepositories()))
	pushCmd.PersistentFlags().StringVarP(&input, flagInput, "i", "", "Input configuration file to use. If not specified, configuration is read from other flags.")
	pushCmd.PersistentFlags().StringVar(&baseDomain, flagBaseDomain, viper.GetString(envBaseDomain), "Base domain to use for the cluster")
	pushCmd.PersistentFlags().StringVar(&team, flagTeam, viper.GetString(envTeam), "Name of the team that owns the cluster")
	pushCmd.PersistentFlags().StringVar(&awsRegion, flagAWSRegion, viper.GetString(envAWSRegion), "AWS region of the cluster")
	pushCmd.PersistentFlags().StringVar(&awsAccountID, flagAWSAccountID, viper.GetString(envAWSAccountID), "AWS account ID of the cluster")
	pushCmd.PersistentFlags().StringVar(&secretFolder, flagSecretFolder, viper.GetString(envSecretFolder), "Secrets folder to use for the cluster")
	pushCmd.PersistentFlags().BoolVar(&mcAppsPreventDeletion, flagMCAppsPreventDeletion, viper.GetBool(envMCAppsPreventDeletion), "Prevent deletion of management cluster apps")
	pushCmd.PersistentFlags().StringVar(&clusterAppName, flagClusterAppName, viper.GetString(envClusterAppName), "Name of the management cluster app")
	pushCmd.PersistentFlags().StringVar(&clusterAppCatalog, flagClusterAppCatalog, viper.GetString(envClusterAppCatalog), "Catalog of the management cluster app")
	pushCmd.PersistentFlags().StringVar(&clusterAppVersion, flagClusterAppVersion, viper.GetString(envClusterAppVersion), "Version of the management cluster app")
	pushCmd.PersistentFlags().StringVar(&clusterNamespace, flagClusterNamespace, viper.GetString(envClusterNamespace), "Namespace of the management cluster")
	pushCmd.PersistentFlags().BoolVar(&configureContainerRegistries, flagConfigureContainerRegistries, viper.GetBool(envConfigureContainerRegistries), "Configure container registries")
	pushCmd.PersistentFlags().StringVar(&defaultAppsName, flagDefaultAppsName, viper.GetString(envDefaultAppsName), "Name of the default apps")
	pushCmd.PersistentFlags().StringVar(&defaultAppsCatalog, flagDefaultAppsCatalog, viper.GetString(envDefaultAppsCatalog), "Catalog of the default apps")
	pushCmd.PersistentFlags().StringVar(&defaultAppsVersion, flagDefaultAppsVersion, viper.GetString(envDefaultAppsVersion), "Version of the default apps")
	pushCmd.PersistentFlags().BoolVar(&privateCA, flagPrivateCA, viper.GetBool(envPrivateCA), "Use private CA")
	pushCmd.PersistentFlags().BoolVar(&certManagerDNSChallenge, flagCertManagerDNSChallenge, viper.GetBool(envCertManagerDNSChallenge), "Use cert-manager DNS01 challenge")
	pushCmd.PersistentFlags().BoolVar(&mcCustomCoreDNSConfig, flagMCCustomCoreDNSConfig, viper.GetBool(envMCCustomCoreDNSConfig), "Use custom CoreDNS config")
	pushCmd.PersistentFlags().BoolVar(&mcProxyEnabled, flagMCProxyEnabled, viper.GetBool(envMCProxyEnabled), "Use proxy")
	pushCmd.PersistentFlags().StringVar(&mcHTTPSProxy, flagMCHTTPSProxy, viper.GetString(envMCHTTPSProxy), "HTTPS proxy to use")
}

func validatePush(cmd *cobra.Command, args []string) error {
	if input != "" {
		_, err := os.Stat(input)
		if err != nil {
			return fmt.Errorf("input file %s can not be accessed.\n%w", input, err)
		}
		log.Debug().Msg(fmt.Sprintf("using input file %s", input))
		return nil
	}
	if installationsBranch == "" {
		return invalidFlagError(flagInstallationsBranch)
	}
	if provider == "" {
		return invalidFlagError(flagProvider)
	}
	return nil
}

func defaultPush() {
	if installationsBranch == "" {
		installationsBranch = key.GetDefaultPRBranch(cluster)
	}
	if cmcBranch == "" {
		cmcBranch = key.GetDefaultPRBranch(cluster)
	}
	if customer == "" {
		customer = "giantswarm"
	}
	if cmcRepository == "" {
		cmcRepository = key.GetCMCName(customer)
	}
}
