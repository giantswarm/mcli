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

// installations flags
const (
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

var (
	baseDomain   string
	team         string
	awsRegion    string
	awsAccountID string
)

// cmc flags
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
	flagAgePubKey                    = "age-pub-key"
	flagAgeKey                       = "age-key"
	flagTaylorBotToken               = "taylor-bot-token"
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
	envAgePubKey                    = "AGE_PUBKEY"
	envAgeKey                       = "AGE_KEY"
)

var (
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
	mcCustomCoreDNSConfig        string
	mcProxyEnabled               bool
	mcHTTPSProxy                 string
	agePubKey                    string
	ageKey                       string
)

// extra cmc flags that are read from the secrets folder and not exposed
const (
	flagDeployKey                          = "deploy-key-passphrase"
	flagDeployKeyIdentity                  = "deploy-key-identity"
	flagDeployKeyKnownHosts                = "deploy-key-known-hosts"
	flagCustomerDeployKey                  = "customer-deploy-key-passphrase"
	flagCustomerDeployKeyIdentity          = "customer-deploy-key-identity"
	flagCustomerDeployKeyKnownHosts        = "customer-deploy-key-known-hosts"
	flagSharedDeployKey                    = "shared-deploy-key-passphrase"
	flagSharedDeployKeyIdentity            = "shared-deploy-key-identity"
	flagSharedDeployKeyKnownHosts          = "shared-deploy-key-known-hosts"
	flagVSphereCredentials                 = "vsphere-credentials"
	flagCloudDirectorCredentials           = "cloud-director-credentials"
	flagAzureClusterIdentityUA             = "azure-cluster-identity-ua"
	flagAzureClusterIdentitySP             = "azure-cluster-identity-sp"
	flagAzureSecretClusterIdentityStaticSP = "azure-secret-cluster-identity-static-sp"
	flagContainerRegistryConfiguration     = "container-registry-configuration"
	flagClusterValues                      = "cluster-values"
	flagCertManagerRoute53Region           = "cert-manager-route53-region"
	flagCertManagerRoute53Role             = "cert-manager-route53-role"
	flagCertManagerRoute53AccessKeyID      = "cert-manager-route53-access-key-id"
	flagCertManagerRoute53SecretAccessKey  = "cert-manager-route53-secret-access-key"
)

var (
	taylorBotToken                     string
	deployKeyPassphrase                string
	deployKeyIdentity                  string
	deployKeyKnownHosts                string
	customerDeployKeyPassphrase        string
	customerDeployKeyIdentity          string
	customerDeployKeyKnownHosts        string
	sharedDeployKeyPassphrase          string
	sharedDeployKeyIdentity            string
	sharedDeployKeyKnownHosts          string
	vSphereCredentials                 string
	cloudDirectorCredentials           string
	azureClusterIdentityUA             string
	azureClusterIdentitySP             string
	azureSecretClusterIdentityStaticSP string
	containerRegistryConfiguration     string
	clusterValues                      string
	certManagerRoute53Region           string
	certManagerRoute53Role             string
	certManagerRoute53AccessKeyID      string
	certManagerRoute53SecretAccessKey  string
)

func addFlagsPush() {
	// add general flags
	pushCmd.Flags().StringArrayVarP(&skip, flagSkip, "s", []string{}, fmt.Sprintf("List of repositories to skip. (default: none) Valid values: %s", key.GetValidRepositories()))
	pushCmd.PersistentFlags().StringVarP(&input, flagInput, "i", "", "Input configuration file to use. If not specified, configuration is read from other flags.")

	// add installations flags
	pushCmd.PersistentFlags().StringVar(&baseDomain, flagBaseDomain, viper.GetString(envBaseDomain), "Base domain to use for the cluster")
	pushCmd.PersistentFlags().StringVar(&team, flagTeam, viper.GetString(envTeam), "Name of the team that owns the cluster")
	pushCmd.PersistentFlags().StringVar(&awsRegion, flagAWSRegion, viper.GetString(envAWSRegion), "AWS region of the cluster")
	pushCmd.PersistentFlags().StringVar(&awsAccountID, flagAWSAccountID, viper.GetString(envAWSAccountID), "AWS account ID of the cluster")

	// add cmc flags
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
	pushCmd.PersistentFlags().StringVar(&mcCustomCoreDNSConfig, flagMCCustomCoreDNSConfig, viper.GetString(envMCCustomCoreDNSConfig), "Custom CoreDNS configuration")
	pushCmd.PersistentFlags().BoolVar(&mcProxyEnabled, flagMCProxyEnabled, viper.GetBool(envMCProxyEnabled), "Use proxy")
	pushCmd.PersistentFlags().StringVar(&mcHTTPSProxy, flagMCHTTPSProxy, viper.GetString(envMCHTTPSProxy), "HTTPS proxy to use")
	pushCmd.PersistentFlags().StringVar(&agePubKey, flagAgePubKey, viper.GetString(envAgePubKey), "Age public key for the cluster")
	pushCmd.PersistentFlags().StringVar(&ageKey, flagAgeKey, viper.GetString(envAgeKey), "Age key for the cluster")
	err := pushCmd.PersistentFlags().MarkHidden(flagAgePubKey)
	if err != nil {
		panic(err)
	}
	err = pushCmd.PersistentFlags().MarkHidden(flagAgeKey)
	if err != nil {
		panic(err)
	}
	pushCmd.PersistentFlags().StringVar(&taylorBotToken, flagTaylorBotToken, "", "Taylor bot token")
	err = pushCmd.PersistentFlags().MarkHidden(flagTaylorBotToken)
	if err != nil {
		panic(err)
	}

	// add extra cmc flags
	pushCmd.PersistentFlags().StringVar(&deployKeyPassphrase, flagDeployKey, "", "Deploy key")
	err = pushCmd.PersistentFlags().MarkHidden(flagDeployKey)
	if err != nil {
		panic(err)
	}
	pushCmd.PersistentFlags().StringVar(&deployKeyIdentity, flagDeployKeyIdentity, "", "Deploy key identity")
	err = pushCmd.PersistentFlags().MarkHidden(flagDeployKeyIdentity)
	if err != nil {
		panic(err)
	}
	pushCmd.PersistentFlags().StringVar(&deployKeyKnownHosts, flagDeployKeyKnownHosts, "", "Deploy key known hosts")
	err = pushCmd.PersistentFlags().MarkHidden(flagDeployKeyKnownHosts)
	if err != nil {
		panic(err)
	}
	pushCmd.PersistentFlags().StringVar(&customerDeployKeyPassphrase, flagCustomerDeployKey, "", "Customer deploy key")
	err = pushCmd.PersistentFlags().MarkHidden(flagCustomerDeployKey)
	if err != nil {
		panic(err)
	}
	pushCmd.PersistentFlags().StringVar(&customerDeployKeyIdentity, flagCustomerDeployKeyIdentity, "", "Customer deploy key identity")
	err = pushCmd.PersistentFlags().MarkHidden(flagCustomerDeployKeyIdentity)
	if err != nil {
		panic(err)
	}
	pushCmd.PersistentFlags().StringVar(&customerDeployKeyKnownHosts, flagCustomerDeployKeyKnownHosts, "", "Customer deploy key known hosts")
	err = pushCmd.PersistentFlags().MarkHidden(flagCustomerDeployKeyKnownHosts)
	if err != nil {
		panic(err)
	}
	pushCmd.PersistentFlags().StringVar(&sharedDeployKeyPassphrase, flagSharedDeployKey, "", "Shared deploy key")
	err = pushCmd.PersistentFlags().MarkHidden(flagSharedDeployKey)
	if err != nil {
		panic(err)
	}
	pushCmd.PersistentFlags().StringVar(&sharedDeployKeyIdentity, flagSharedDeployKeyIdentity, "", "Shared deploy key identity")
	err = pushCmd.PersistentFlags().MarkHidden(flagSharedDeployKeyIdentity)
	if err != nil {
		panic(err)
	}
	pushCmd.PersistentFlags().StringVar(&sharedDeployKeyKnownHosts, flagSharedDeployKeyKnownHosts, "", "Shared deploy key known hosts")
	err = pushCmd.PersistentFlags().MarkHidden(flagSharedDeployKeyKnownHosts)
	if err != nil {
		panic(err)
	}
	pushCmd.PersistentFlags().StringVar(&vSphereCredentials, flagVSphereCredentials, "", "vSphere credentials")
	err = pushCmd.PersistentFlags().MarkHidden(flagVSphereCredentials)
	if err != nil {
		panic(err)
	}
	pushCmd.PersistentFlags().StringVar(&cloudDirectorCredentials, flagCloudDirectorCredentials, "", "Cloud Director credentials")
	err = pushCmd.PersistentFlags().MarkHidden(flagCloudDirectorCredentials)
	if err != nil {
		panic(err)
	}
	pushCmd.PersistentFlags().StringVar(&azureClusterIdentityUA, flagAzureClusterIdentityUA, "", "Azure cluster identity user-assigned")
	err = pushCmd.PersistentFlags().MarkHidden(flagAzureClusterIdentityUA)
	if err != nil {
		panic(err)
	}
	pushCmd.PersistentFlags().StringVar(&azureClusterIdentitySP, flagAzureClusterIdentitySP, "", "Azure cluster identity service principal")
	err = pushCmd.PersistentFlags().MarkHidden(flagAzureClusterIdentitySP)
	if err != nil {
		panic(err)
	}
	pushCmd.PersistentFlags().StringVar(&azureSecretClusterIdentityStaticSP, flagAzureSecretClusterIdentityStaticSP, "", "Azure secret cluster identity static service principal")
	err = pushCmd.PersistentFlags().MarkHidden(flagAzureSecretClusterIdentityStaticSP)
	if err != nil {
		panic(err)
	}
	pushCmd.PersistentFlags().StringVar(&containerRegistryConfiguration, flagContainerRegistryConfiguration, "", "Container registry configuration")
	err = pushCmd.PersistentFlags().MarkHidden(flagContainerRegistryConfiguration)
	if err != nil {
		panic(err)
	}
	pushCmd.PersistentFlags().StringVar(&clusterValues, flagClusterValues, "", "Cluster values")
	err = pushCmd.PersistentFlags().MarkHidden(flagClusterValues)
	if err != nil {
		panic(err)
	}
	pushCmd.PersistentFlags().StringVar(&certManagerRoute53Region, flagCertManagerRoute53Region, "", "Cert Manager Route53 region")
	err = pushCmd.PersistentFlags().MarkHidden(flagCertManagerRoute53Region)
	if err != nil {
		panic(err)
	}
	pushCmd.PersistentFlags().StringVar(&certManagerRoute53Role, flagCertManagerRoute53Role, "", "Cert Manager Route53 role")
	err = pushCmd.PersistentFlags().MarkHidden(flagCertManagerRoute53Role)
	if err != nil {
		panic(err)
	}
	pushCmd.PersistentFlags().StringVar(&certManagerRoute53AccessKeyID, flagCertManagerRoute53AccessKeyID, "", "Cert Manager Route53 access key ID")
	err = pushCmd.PersistentFlags().MarkHidden(flagCertManagerRoute53AccessKeyID)
	if err != nil {
		panic(err)
	}
	pushCmd.PersistentFlags().StringVar(&certManagerRoute53SecretAccessKey, flagCertManagerRoute53SecretAccessKey, "", "Cert Manager Route53 secret access key")
	err = pushCmd.PersistentFlags().MarkHidden(flagCertManagerRoute53SecretAccessKey)
	if err != nil {
		panic(err)
	}
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
	if cmcBranch == "" {
		return invalidFlagError(flagCMCBranch)
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
