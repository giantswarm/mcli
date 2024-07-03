package kustomization

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	kustomize "sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/resid"

	"github.com/giantswarm/mcli/pkg/key"
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
	CertManagerConfigMapFile           = "cert-manager-configmap.yaml"
	IssuerFile                         = "private-cluster-issuer.yaml"
	VsphereCredentialsFile             = "vsphere-cloud-config-secret.yaml" // #nosec G101
	CloudDirectorCredentialsFile       = "cloud-director-cloud-config-secret.yaml"
	AzureClusterIdentitySPFile         = "azureclusteridentity-sp.yaml"
	AzureClusterIdentityUAFile         = "azureclusteridentity-ua.yaml"
	AzureSecretClusterIdentityStaticSP = "secret-clusteridentity-static-sp.yaml"
	ExternalDNSFile                    = "external-dns-configmap.yaml"
)

func GetSourceControllerPatchPath(branch string) string {
	return fmt.Sprintf("https://raw.githubusercontent.com/giantswarm/management-cluster-bases/%s/extras/flux/patch-source-controller-deployment-host-alias.yaml", branch)
}
func GetSourceControllerSocatSidecarPath(branch string) string {
	return fmt.Sprintf("https://raw.githubusercontent.com/giantswarm/management-cluster-bases/%s/extras/flux/patch-source-controller-deployment-socat.yaml", branch)
}

type Config struct {
	CertManagerDNSChallenge      bool
	Provider                     string
	PrivateCA                    bool
	PrivateMC                    bool
	ConfigureContainerRegistries bool
	CustomCoreDNS                bool
	DisableDenyAllNetPol         bool
	MCProxy                      bool
	IntegratedDefaultAppsValues  bool
	MCBBranchSource              string
}

func GetKustomizationConfig(file string) (Config, error) {
	log.Debug().Msg("Getting Kustomization configuration")
	kustomization, err := getKustomization(file)
	if err != nil {
		return Config{}, err
	}
	c := Config{
		CertManagerDNSChallenge:      containsResource(kustomization.Resources, CertManagerFile),
		Provider:                     getProvider(kustomization.Resources),
		PrivateCA:                    containsResource(kustomization.Resources, IssuerFile),
		ConfigureContainerRegistries: containsResource(kustomization.Resources, RegistryFile),
		CustomCoreDNS:                containsResource(kustomization.Resources, CoreDNSFile),
		DisableDenyAllNetPol:         !containsResource(kustomization.Resources, DenyNetPolFile),
		MCProxy:                      containsPatch(kustomization.Patches, SourceControllerFile),
		IntegratedDefaultAppsValues:  !containsResource(kustomization.Resources, DefaultAppsFile),
	}
	c.PrivateMC = key.IsProviderAzure(c.Provider) && c.IntegratedDefaultAppsValues && !containsResource(kustomization.Resources, ExternalDNSFile)

	return c, nil
}

func GetKustomizationFile(c Config, file string) (string, error) {
	log.Debug().Msg("Creating Kustomization file")

	k, err := getKustomization(file)
	if err != nil {
		return "", err
	}
	k.Resources = appendResource(k.Resources, TaylorBotFile)
	k.Resources = appendResource(k.Resources, SSHdeployKeyFile)
	k.Resources = appendResource(k.Resources, CustomerDeployKeyFile)
	k.Resources = appendResource(k.Resources, SharedDeployKeyFile)
	k.Resources = appendResource(k.Resources, AgeKeyFile)

	if c.CertManagerDNSChallenge {
		k.Resources = appendResource(k.Resources, CertManagerFile)
	}
	if key.IsProviderVsphere(c.Provider) {
		k.Resources = appendResource(k.Resources, VsphereCredentialsFile)
	} else if key.IsProviderVCD(c.Provider) {
		k.Resources = appendResource(k.Resources, CloudDirectorCredentialsFile)
	} else if key.IsProviderAzure(c.Provider) {
		k.Resources = appendResource(k.Resources, AzureClusterIdentitySPFile)
		k.Resources = appendResource(k.Resources, AzureClusterIdentityUAFile)
		k.Resources = appendResource(k.Resources, AzureSecretClusterIdentityStaticSP)
		if !c.PrivateMC && c.IntegratedDefaultAppsValues {
			k.Resources = appendResource(k.Resources, ExternalDNSFile)
		}
	}
	if c.PrivateCA {
		k.Resources = appendResource(k.Resources, IssuerFile)
		if c.IntegratedDefaultAppsValues {
			k.Resources = appendResource(k.Resources, CertManagerConfigMapFile)
		}
	}
	if c.ConfigureContainerRegistries {
		k.Resources = appendResource(k.Resources, RegistryFile)
	}
	if c.CustomCoreDNS {
		k.Resources = appendResource(k.Resources, CoreDNSFile)
	}
	if c.DisableDenyAllNetPol {
		k.Resources = removeResource(k.Resources, DenyNetPolFile)
	}
	if c.MCProxy {
		k.Patches = appendPatch(k.Patches, kustomize.Patch{Path: GetSourceControllerPatchPath(c.MCBBranchSource)})
		k.Patches = appendPatch(k.Patches, kustomize.Patch{Path: GetSourceControllerSocatSidecarPath(c.MCBBranchSource),
			Target: &kustomize.Selector{
				ResId: resid.ResId{
					Gvk: resid.Gvk{
						Kind: "Deployment",
					},
					Name:      "source-controller",
					Namespace: key.FluxNamespace,
				},
			},
		})
		k.Patches = appendPatch(k.Patches, kustomize.Patch{Path: SourceControllerFile})
	}
	if c.IntegratedDefaultAppsValues {
		k.Resources = removeResource(k.Resources, DefaultAppsFile)
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

func appendResource(resources []string, resource string) []string {
	for _, r := range resources {
		if r == resource {
			return resources
		}
	}
	return append(resources, resource)
}

func appendPatch(patches []kustomize.Patch, patch kustomize.Patch) []kustomize.Patch {
	for _, p := range patches {
		if p.Path == patch.Path {
			return patches
		}
	}
	return append(patches, patch)
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
