package kustomization

import (
	"fmt"
	"strings"

	"github.com/giantswarm/mcli/pkg/key"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
	kustomize "sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/resid"
)

const (
	AgeKeyFile            = "age-secret-keys.yaml"
	ClusterAppsFile       = "cluster-app-manifests.yaml"
	DefaultAppsFile       = "default-apps-manifests.yaml"
	KustomizationFile     = "kustomization.yaml"
	TaylorBotFile         = "github-giantswarm-https-credentials.yaml"
	SSHdeployKeyFile      = "giantswarm-clusters-ssh-credentials.yaml"
	CustomerDeployKeyFile = "configs-ssh-credentials.yaml"
	SharedDeployKeyFile   = "shared-configs-ssh-credentials.yaml"
)

const (
	AllowNetPolFile                    = "networkpolicy-egress-to-proxy.yaml"
	SourceControllerFile               = "kustomization-post-build-proxy.yaml"
	RegistryFile                       = "container-registries-configuration-secret.yaml"
	CoreDNSFile                        = "coredns-configmap.yaml"
	DenyNetPolFile                     = "deny-all-policies.yaml"
	CertManagerFile                    = "cert-manager-dns01-secret.yaml"
	IssuerFile                         = "private-cluster-issuer.yaml"
	VsphereCredentialsFile             = "vsphere-cloud-config-secret.yaml"
	CloudDirectorCredentialsFile       = "cloud-director-cloud-config-secret.yaml"
	AzureClusterIdentitySPFile         = "azureclusteridentity-sp.yaml"
	AzureClusterIdentityUAFile         = "azureclusteridentity-ua.yaml"
	AzureSecretClusterIdentityStaticSP = "secret-clusteridentity-static-sp.yaml"
)

func GetSourceControllerPatchPath() string {
	return fmt.Sprintf("https://raw.githubusercontent.com/giantswarm/management-cluster-bases/%s/extras/flux/patch-source-controller-deployment-host-alias.yaml", key.CMCMainBranch)
}
func GetSharedConfigsPatchPath() string {
	return fmt.Sprintf("https://raw.githubusercontent.com/giantswarm/management-cluster-bases/%s/extras/flux/patch-gitrepository-shared-configs-private.yaml", key.CMCMainBranch)
}
func GetSourceControllerSocatSidecarPath() string {
	return fmt.Sprintf("https://raw.githubusercontent.com/giantswarm/management-cluster-bases/%s/extras/flux/patch-source-controller-deployment-socat.yaml", key.CMCMainBranch)
}

type Config struct {
	CertManagerDNSChallenge      bool
	Provider                     string
	PrivateCA                    bool
	ConfigureContainerRegistries bool
	CustomCoreDNS                bool
	DisableDenyAllNetPol         bool
	MCProxy                      bool
}

func GetKustomizationConfig(file string) (Config, error) {
	log.Debug().Msg("Getting Kustomization configuration")
	kustomization, err := getKustomization(file)
	if err != nil {
		return Config{}, err
	}
	return Config{
		CertManagerDNSChallenge:      containsResource(kustomization.Resources, CertManagerFile),
		Provider:                     getProvider(kustomization.Resources),
		PrivateCA:                    containsResource(kustomization.Resources, IssuerFile),
		ConfigureContainerRegistries: containsResource(kustomization.Resources, RegistryFile),
		CustomCoreDNS:                containsResource(kustomization.Resources, CoreDNSFile),
		DisableDenyAllNetPol:         !containsResource(kustomization.Resources, DenyNetPolFile),
		MCProxy:                      containsPatch(kustomization.Patches, SourceControllerFile),
	}, nil
}

func GetKustomizationFile(c Config, file string) (string, error) {
	log.Debug().Msg("Creating Kustomization file")

	if c.MCProxy {
		file = strings.ReplaceAll(file, "patch-gitrepository-mcf.yaml", "patch-gitrepository-mcf-private.yaml")
		file = strings.ReplaceAll(file, "patch-gitrepository-config.yaml", "patch-gitrepository-config-private.yaml")
	}

	k, err := getKustomization(file)
	if err != nil {
		return "", err
	}
	k.Resources = append(k.Resources,
		TaylorBotFile,
		SSHdeployKeyFile,
		CustomerDeployKeyFile,
		SharedDeployKeyFile,
		AgeKeyFile,
	)

	if c.CertManagerDNSChallenge {
		k.Resources = append(k.Resources, CertManagerFile)
	}
	if key.IsProviderVsphere(c.Provider) {
		k.Resources = append(k.Resources, VsphereCredentialsFile)
	} else if key.IsProviderVCD(c.Provider) {
		k.Resources = append(k.Resources, CloudDirectorCredentialsFile)
	} else if key.IsProviderAzure(c.Provider) {
		k.Resources = append(k.Resources,
			AzureClusterIdentitySPFile,
			AzureClusterIdentityUAFile,
			AzureSecretClusterIdentityStaticSP)
	}
	if c.PrivateCA {
		k.Resources = append(k.Resources, IssuerFile)
	}
	if c.ConfigureContainerRegistries {
		k.Resources = append(k.Resources, RegistryFile)
	}
	if c.CustomCoreDNS {
		k.Resources = append(k.Resources, CoreDNSFile)
	}
	if c.DisableDenyAllNetPol {
		k.Resources = removeResource(k.Resources, DenyNetPolFile)
	}
	if c.MCProxy {
		k.Patches = append(k.Patches,
			kustomize.Patch{Path: GetSourceControllerPatchPath()},
			kustomize.Patch{Path: GetSharedConfigsPatchPath()},
			kustomize.Patch{Path: GetSourceControllerSocatSidecarPath(),
				Target: &kustomize.Selector{
					ResId: resid.ResId{
						Gvk: resid.Gvk{
							Kind: "Deployment",
						},
						Name:      "source-controller",
						Namespace: key.FluxNamespace,
					},
				},
			},
			kustomize.Patch{Path: SourceControllerFile},
		)
	}
	data, err := key.GetData(k)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func getKustomization(file string) (kustomize.Kustomization, error) {
	k := kustomize.Kustomization{}
	if err := yaml.Unmarshal([]byte(file), &k); err != nil {
		return k, fmt.Errorf("failed to unmarshal kustomization object.\n%w", err)
	}
	return k, nil
}

func removeResource(resources []string, resource string) []string {
	for i, r := range resources {
		if r == resource {
			return append(resources[:i], resources[i+1:]...)
		}
	}
	return resources
}

func containsResource(resources []string, resource string) bool {
	for _, r := range resources {
		if r == resource {
			return true
		}
	}
	return false
}

func containsPatch(patches []kustomize.Patch, patch string) bool {
	for _, p := range patches {
		if p.Path == patch {
			return true
		}
	}
	return false
}

func getProvider(resources []string) string {
	for _, r := range resources {
		if r == VsphereCredentialsFile {
			return key.ProviderVsphere
		} else if r == CloudDirectorCredentialsFile {
			return key.ProviderVCD
		} else if r == AzureClusterIdentitySPFile {
			return key.ProviderAzure
		}
	}
	return ""
}
