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
	flagCCRRepository = "ccr-repository"
	flagPipeline      = "pipeline"
	flagTeam          = "team"
	flagAWSRegion     = "aws-region"
	flagAWSAccountID  = "aws-account-id"
)

const (
	envCCRRepository = "CCR_REPOSITORY"
	envPipeline      = "MC_PIPELINE"
	envTeam          = "TEAM_NAME"
	envAWSRegion     = "AWS_REGION"
	envAWSAccountID  = "INSTALLATION_AWS_ACCOUNT"
)

var (
	ccrRepository string
	pipeline      string
	team          string
	awsRegion     string
	awsAccountID  string
)

// cmc flags
const (
	flagSecretFolder                 = "secret-folder"
	flagMCAppsPreventDeletion        = "mc-apps-prevent-deletion"
	flagClusterAppName               = "cluster-app-name"
	flagClusterAppCatalog            = "cluster-app-catalog"
	flagClusterAppVersion            = "cluster-app-version"
	flagClusterNamespace             = "cluster-namespace"
	flagClusterIntegratesDefaultApps = "cluster-integrates-default-apps"
	flagConfigureContainerRegistries = "configure-container-registries"
	flagDefaultAppsName              = "default-apps-name"
	flagDefaultAppsCatalog           = "default-apps-catalog"
	flagDefaultAppsVersion           = "default-apps-version"
	flagPrivateCA                    = "private-ca"
	flagPrivateMC                    = "private-mc"
	flagCertManagerDNSChallenge      = "cert-manager-dns-challenge"
	flagMCCustomCoreDNSConfig        = "mc-custom-coredns-config"
	flagMCProxyEnabled               = "mc-proxy-enabled"
	flagMCHTTPSProxy                 = "mc-https-proxy"
	flagAgePubKey                    = "age-pub-key"
	flagTaylorBotToken               = "taylor-bot-token"
	flagMCBBranchSource              = "mcb-branch-source"
	flagConfigBranch                 = "config-branch"
	flagMCAppCollectionBranch        = "mc-app-collection-branch"
	flagRegistryDomain               = "registry-domain"
)

const (
	envSecretFolder                 = "SECRETS_FOLDER"
	envMCAppsPreventDeletion        = "MC_APPS_PREVENT_DELETION"
	envClusterAppName               = "CLUSTER_APP_NAME"
	envClusterAppCatalog            = "CLUSTER_APP_CATALOG"
	envClusterAppVersion            = "CLUSTER_APP_VERSION"
	envClusterNamespace             = "CLUSTER_NAMESPACE"
	envClusterIntegratesDefaultApps = "CLUSTER_INTEGRATES_DEFAULT_APPS"
	envConfigureContainerRegistries = "CONFIGURE_CONTAINER_REGISTRIES"
	envDefaultAppsName              = "DEFAULT_APPS_APP_NAME"
	envDefaultAppsCatalog           = "DEFAULT_APPS_APP_CATALOG"
	envDefaultAppsVersion           = "DEFAULT_APPS_APP_VERSION"
	envPrivateCA                    = "PRIVATE_CA"
	envPrivateMC                    = "MC_PRIVATE"
	envCertManagerDNSChallenge      = "CERT_MANAGER_DNS01_CHALLENGE"
	envMCCustomCoreDNSConfig        = "MC_CUSTOM_COREDNS_CONFIG"
	envMCProxyEnabled               = "MC_PROXY_ENABLED"
	envMCHTTPSProxy                 = "MC_HTTPS_PROXY"
	envAgePubKey                    = "AGE_PUBKEY"
	envMCBBranchSource              = "MCB_BRANCH_SOURCE"
	envConfigBranch                 = "CONFIG_BRANCH"
	envMCAppCollectionBranch        = "MC_APP_COLLECTION_BRANCH"
	envRegistryDomain               = "REGISTRY_DOMAIN"
)

var (
	secretFolder                 string
	mcAppsPreventDeletion        bool
	clusterAppName               string
	clusterAppCatalog            string
	clusterAppVersion            string
	clusterNamespace             string
	clusterIntegratesDefaultApps bool
	configureContainerRegistries bool
	defaultAppsName              string
	defaultAppsCatalog           string
	defaultAppsVersion           string
	privateCA                    bool
	privateMC                    bool
	certManagerDNSChallenge      bool
	mcCustomCoreDNSConfig        string
	mcProxyEnabled               bool
	mcHTTPSProxy                 string
	agePubKey                    string
	mcbBranchSource              string
	configBranch                 string
	mcAppCollectionBranch        string
	registryDomain               string
)

// extra cmc flags that are read from the secrets folder and not exposed
const (
	flagDeployKey                         = "deploy-key-passphrase"
	flagDeployKeyIdentity                 = "deploy-key-identity"
	flagDeployKeyKnownHosts               = "deploy-key-known-hosts"
	flagCustomerDeployKey                 = "customer-deploy-key-passphrase"
	flagCustomerDeployKeyIdentity         = "customer-deploy-key-identity"
	flagCustomerDeployKeyKnownHosts       = "customer-deploy-key-known-hosts"
	flagSharedDeployKey                   = "shared-deploy-key-passphrase"
	flagSharedDeployKeyIdentity           = "shared-deploy-key-identity"
	flagSharedDeployKeyKnownHosts         = "shared-deploy-key-known-hosts"
	flagVSphereCredentials                = "vsphere-credentials" // #nosec G101
	flagCloudDirectorRefreshToken         = "cloud-director-refresh-token"
	flagAzureUAClientID                   = "azure-ua-client-id"
	flagAzureUATenantID                   = "azure-ua-tenant-id"
	flagAzureUAResourceID                 = "azure-ua-resource-id"
	flagAzureClientID                     = "azure-client-id"
	flagAzureTenantID                     = "azure-tenant-id"
	flagAzureClientSecret                 = "azure-client-secret"
	flagAzureSubscriptionID               = "azure-subscription-id"
	flagContainerRegistryConfiguration    = "container-registry-configuration"
	flagClusterValues                     = "cluster-values"
	flagCertManagerRoute53Region          = "cert-manager-route53-region"
	flagCertManagerRoute53Role            = "cert-manager-route53-role"
	flagCertManagerRoute53AccessKeyID     = "cert-manager-route53-access-key-id"
	flagCertManagerRoute53SecretAccessKey = "cert-manager-route53-secret-access-key" // #nosec G101
)

var (
	taylorBotToken                    string
	deployKeyPassphrase               string
	deployKeyIdentity                 string
	deployKeyKnownHosts               string
	customerDeployKeyPassphrase       string
	customerDeployKeyIdentity         string
	customerDeployKeyKnownHosts       string
	sharedDeployKeyPassphrase         string
	sharedDeployKeyIdentity           string
	sharedDeployKeyKnownHosts         string
	vSphereCredentials                string
	cloudDirectorRefreshToken         string
	azureUAClientID                   string
	azureUATenantID                   string
	azureUAResourceID                 string
	azureClientID                     string
	azureTenantID                     string
	azureClientSecret                 string
	azureSubscriptionID               string
	containerRegistryConfiguration    string
	clusterValues                     string
	certManagerRoute53Region          string
	certManagerRoute53Role            string
	certManagerRoute53AccessKeyID     string
	certManagerRoute53SecretAccessKey string
)

func addFlagsPush() {
	// add general flags
	pushCmd.Flags().StringArrayVarP(&skip, flagSkip, "s", []string{}, fmt.Sprintf("List of repositories to skip. (default: none) Valid values: %s", key.GetValidRepositories()))
	pushCmd.PersistentFlags().StringVarP(&input, flagInput, "i", "", "Input configuration file to use. If not specified, configuration is read from other flags.")
	pushCmd.PersistentFlags().StringVar(&provider, flagProvider, viper.GetString(envProvider), "Provider of the cluster")
	pushCmd.PersistentFlags().StringVar(&baseDomain, flagBaseDomain, viper.GetString(envBaseDomain), "Base domain to use for the cluster")

	// add installations flags
	pushCmd.PersistentFlags().StringVar(&ccrRepository, flagCCRRepository, viper.GetString(envCCRRepository), "CCR repository to use for the cluster")
	pushCmd.PersistentFlags().StringVar(&pipeline, flagPipeline, viper.GetString(envPipeline), "Pipeline to use for the cluster")
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
	pushCmd.PersistentFlags().BoolVar(&clusterIntegratesDefaultApps, flagClusterIntegratesDefaultApps, viper.GetBool(envClusterIntegratesDefaultApps), "Integrate default apps")
	pushCmd.PersistentFlags().BoolVar(&configureContainerRegistries, flagConfigureContainerRegistries, viper.GetBool(envConfigureContainerRegistries), "Configure container registries")
	pushCmd.PersistentFlags().StringVar(&defaultAppsName, flagDefaultAppsName, viper.GetString(envDefaultAppsName), "Name of the default apps")
	pushCmd.PersistentFlags().StringVar(&defaultAppsCatalog, flagDefaultAppsCatalog, viper.GetString(envDefaultAppsCatalog), "Catalog of the default apps")
	pushCmd.PersistentFlags().StringVar(&defaultAppsVersion, flagDefaultAppsVersion, viper.GetString(envDefaultAppsVersion), "Version of the default apps")
	pushCmd.PersistentFlags().BoolVar(&privateCA, flagPrivateCA, viper.GetBool(envPrivateCA), "Use private CA")
	pushCmd.PersistentFlags().BoolVar(&privateMC, flagPrivateMC, viper.GetBool(envPrivateMC), "MC is private")
	pushCmd.PersistentFlags().BoolVar(&certManagerDNSChallenge, flagCertManagerDNSChallenge, viper.GetBool(envCertManagerDNSChallenge), "Use cert-manager DNS01 challenge")
	pushCmd.PersistentFlags().StringVar(&mcCustomCoreDNSConfig, flagMCCustomCoreDNSConfig, viper.GetString(envMCCustomCoreDNSConfig), "Custom CoreDNS configuration")
	pushCmd.PersistentFlags().BoolVar(&mcProxyEnabled, flagMCProxyEnabled, viper.GetBool(envMCProxyEnabled), "Use proxy")
	pushCmd.PersistentFlags().StringVar(&mcHTTPSProxy, flagMCHTTPSProxy, viper.GetString(envMCHTTPSProxy), "HTTPS proxy to use")
	pushCmd.PersistentFlags().StringVar(&mcbBranchSource, flagMCBBranchSource, viper.GetString(envMCBBranchSource), "Branch to use for the mcb repository")
	pushCmd.PersistentFlags().StringVar(&configBranch, flagConfigBranch, viper.GetString(envConfigBranch), "Branch to use for the config repository")
	pushCmd.PersistentFlags().StringVar(&mcAppCollectionBranch, flagMCAppCollectionBranch, viper.GetString(envMCAppCollectionBranch), "Branch to use for the MC app collection repository")
	pushCmd.PersistentFlags().StringVar(&registryDomain, flagRegistryDomain, viper.GetString(envRegistryDomain), "Domain of the registry. Only needed if it's different from the default.")
	pushCmd.PersistentFlags().StringVar(&agePubKey, flagAgePubKey, viper.GetString(envAgePubKey), "Age public key for the cluster")
	err := pushCmd.PersistentFlags().MarkHidden(flagAgePubKey)
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

	pushCmd.PersistentFlags().StringVar(&deployKeyKnownHosts, flagDeployKeyKnownHosts, "", "Deploy key known hosts")

	pushCmd.PersistentFlags().StringVar(&customerDeployKeyPassphrase, flagCustomerDeployKey, "", "Customer deploy key")
	err = pushCmd.PersistentFlags().MarkHidden(flagCustomerDeployKey)
	if err != nil {
		panic(err)
	}
	pushCmd.PersistentFlags().StringVar(&customerDeployKeyIdentity, flagCustomerDeployKeyIdentity, "", "Customer deploy key identity")

	pushCmd.PersistentFlags().StringVar(&customerDeployKeyKnownHosts, flagCustomerDeployKeyKnownHosts, "", "Customer deploy key known hosts")

	pushCmd.PersistentFlags().StringVar(&sharedDeployKeyPassphrase, flagSharedDeployKey, "", "Shared deploy key")
	err = pushCmd.PersistentFlags().MarkHidden(flagSharedDeployKey)
	if err != nil {
		panic(err)
	}
	pushCmd.PersistentFlags().StringVar(&sharedDeployKeyIdentity, flagSharedDeployKeyIdentity, "", "Shared deploy key identity")

	pushCmd.PersistentFlags().StringVar(&sharedDeployKeyKnownHosts, flagSharedDeployKeyKnownHosts, "", "Shared deploy key known hosts")

	pushCmd.PersistentFlags().StringVar(&vSphereCredentials, flagVSphereCredentials, "", "vSphere credentials")
	err = pushCmd.PersistentFlags().MarkHidden(flagVSphereCredentials)
	if err != nil {
		panic(err)
	}
	pushCmd.PersistentFlags().StringVar(&cloudDirectorRefreshToken, flagCloudDirectorRefreshToken, "", "Cloud Director refresh token")
	err = pushCmd.PersistentFlags().MarkHidden(flagCloudDirectorRefreshToken)
	if err != nil {
		panic(err)
	}

	pushCmd.PersistentFlags().StringVar(&azureUAClientID, flagAzureUAClientID, "", "Azure UA client ID")
	pushCmd.PersistentFlags().StringVar(&azureUATenantID, flagAzureUATenantID, "", "Azure UA tenant ID")
	pushCmd.PersistentFlags().StringVar(&azureUAResourceID, flagAzureUAResourceID, "", "Azure UA resource ID")
	pushCmd.PersistentFlags().StringVar(&azureClientID, flagAzureClientID, "", "Azure client ID")
	pushCmd.PersistentFlags().StringVar(&azureTenantID, flagAzureTenantID, "", "Azure tenant ID")
	pushCmd.PersistentFlags().StringVar(&azureSubscriptionID, flagAzureSubscriptionID, "", "Azure subscription ID")
	pushCmd.PersistentFlags().StringVar(&azureClientSecret, flagAzureClientSecret, "", "Azure client secret")
	err = pushCmd.PersistentFlags().MarkHidden(flagAzureClientSecret)
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

	pushCmd.PersistentFlags().StringVar(&certManagerRoute53Role, flagCertManagerRoute53Role, "", "Cert Manager Route53 role")

	pushCmd.PersistentFlags().StringVar(&certManagerRoute53AccessKeyID, flagCertManagerRoute53AccessKeyID, "", "Cert Manager Route53 access key ID")

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
	if mcAppCollectionBranch == "" {
		mcAppCollectionBranch = key.GetDefaultPRBranch(cluster)
	}
	if configBranch == "" {
		configBranch = key.GetDefaultConfigBranch(cluster)
	}
	if customer == "" {
		customer = key.OrganizationGiantSwarm
	}
	if customer == "gs" {
		customer = key.OrganizationGiantSwarm
	}
	if cmcRepository == "" {
		cmcRepository = key.GetCMCName(customer)
	}
	if mcbBranchSource == "" {
		mcbBranchSource = key.MCBMainBranch
	}
	if pipeline == "" {
		pipeline = "testing"
	}
}
