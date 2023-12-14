package key

import "fmt"

const (
	OrganizationGiantSwarm  = "giantswarm"
	RepositoryInstallations = "installations"
	InstallationsMainBranch = "master"
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
	}
}

func GetInstallationsPath(cluster string) string {
	return fmt.Sprintf("%s/cluster.yaml", cluster)
}

func GetDefaultInstallationsBranch(cluster string) string {
	return fmt.Sprintf("%s_auto_branch", cluster)
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
