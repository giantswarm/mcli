package key

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
	OrganizationGiantSwarm  = "giantswarm"
	RepositoryInstallations = "installations"
	RepositoryCMC           = "cmc"
	RepositoryGithub        = "github"
	RepositoryMCBootstrap   = "mc-bootstrap"
	InstallationsMainBranch = "master"
	CMCMainBranch           = "main"
	Employees               = "employees"
	Bots                    = "bots"
	CMCTemplateRepository   = "template-management-clusters"
	CMCEntryTemplatePath    = "scripts/setup-cmc-branch/management-cluster-template"
	FluxNamespace           = "flux-giantswarm"
)

const (
	ClusterValuesFile                  = "cluster-values.yaml"
	CommonSecretsFile                  = "common.secrets"
	VsphereCredentialsFile             = "vsphere-cloud-config-secret.yaml"
	CloudDirectorCredentialsFile       = "cloud-director-cloud-config-secret.yaml"
	AzureClusterIdentitySPFile         = "azureclusteridentity-sp.yaml"
	AzureClusterIdentityUAFile         = "azureclusteridentity-ua.yaml"
	AzureSecretClusterIdentityStaticSP = "secret-clusteridentity-static-sp.yaml"
)

const (
	ProviderAWS     = "capa"
	ProviderAzure   = "capz"
	ProviderVCD     = "cloud-director"
	ProviderVsphere = "vsphere"
	ProviderGCP     = "gcp"
)

func GetValidProviders() []string {
	return []string{
		ProviderAWS,
		ProviderAzure,
		ProviderVCD,
		ProviderVsphere,
		ProviderGCP,
	}
}

func GetValidRepositories() []string {
	return []string{
		RepositoryInstallations,
		RepositoryCMC,
	}
}

func CMCTemplate() string {
	return fmt.Sprintf("%s/%s", OrganizationGiantSwarm, CMCTemplateRepository)
}

func GetInstallationsPath(cluster string) string {
	return fmt.Sprintf("%s/cluster.yaml", cluster)
}

func GetCMCPath(cluster string) string {
	return fmt.Sprintf("management-clusters/%s", cluster)
}

func GetCMCName(customer string) string {
	return fmt.Sprintf("%s-management-clusters", customer)
}

func GetDefaultPRBranch(cluster string) string {
	return fmt.Sprintf("%s_auto_branch", cluster)
}

func GetOwnershipBranch(customer string) string {
	return fmt.Sprintf("add-%s-mc-to-honeybadger-%s", customer, GetRandom())
}

func GetClusterSecretFile(cluster string) string {
	return fmt.Sprintf("%s.secrets", cluster)
}

func GetContainerRegistriesFile(cluster string) string {
	return fmt.Sprintf("%s-container-registries-configuration.yaml", cluster)
}

func GetCertManagerSecretName(cluster string) string {
	return fmt.Sprintf("%s--cert-manager-user-secrets", cluster)
}

func GetDeployKey(cluster string) string {
	return fmt.Sprintf("%s-key", cluster)
}

func GetKnownHosts(cluster string) string {
	return fmt.Sprintf("%s-known-hosts", cluster)
}

func GetPassphrase(cluster string) string {
	return fmt.Sprintf("%s-passphrase", cluster)
}

func GetAgeKey(cluster string) string {
	return fmt.Sprintf("%s.agekey", cluster)
}

func GetRandom() string {
	random, err := rand.Int(rand.Reader, big.NewInt(32767))
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%d", random)
}

func IsValidRepository(repository string) bool {
	for _, validRepository := range GetValidRepositories() {
		if repository == validRepository {
			return true
		}
	}
	return false
}

func IsValidProvider(provider string) bool {
	for _, validProvider := range GetValidProviders() {
		if provider == validProvider {
			return true
		}
	}
	return false
}

func Skip(name string, skip []string) bool {
	for _, s := range skip {
		if s == name {
			return true
		}
	}
	return false
}
