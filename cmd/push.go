/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/giantswarm/mcli/cmd/push"
	pushcmc "github.com/giantswarm/mcli/cmd/push/cmc"
	pushinstallations "github.com/giantswarm/mcli/cmd/push/installations"
	"github.com/giantswarm/mcli/pkg/github"
	"github.com/giantswarm/mcli/pkg/managementcluster/cmc"
	"github.com/giantswarm/mcli/pkg/managementcluster/installations"
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Pushes configuration of a Management Cluster",
	Long: `Pushes configuration of a Management Cluster to all
relevant git repositories. For example:

mcli push --cluster=gigmac --input=cluster.yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		defaultPush()
		err := validateRoot(cmd, args)
		if err != nil {
			return err
		}
		err = validatePush(cmd, args)
		if err != nil {
			return err
		}
		ctx := context.Background()
		c := push.Config{
			Cluster:             cluster,
			GithubToken:         githubToken,
			InstallationsBranch: installationsBranch,
			Skip:                skip,
			Input:               input,
			CMCBranch:           cmcBranch,
			CMCRepository:       cmcRepository,
			Provider:            provider,
			InstallationsFlags: pushinstallations.InstallationsFlags{
				BaseDomain: baseDomain,
				Team:       team,
				Customer:   customer,
				AWS: pushinstallations.AWSFlags{
					Region:                 awsRegion,
					InstallationAWSAccount: awsAccountID,
				},
			},
			CMCFlags: pushcmc.CMCFlags{
				SecretFolder:                 secretFolder,
				AgePubKey:                    agePubKey,
				AgeKey:                       ageKey,
				MCAppsPreventDeletion:        mcAppsPreventDeletion,
				ClusterAppName:               clusterAppName,
				ClusterAppCatalog:            clusterAppCatalog,
				ClusterAppVersion:            clusterAppVersion,
				ClusterNamespace:             clusterNamespace,
				ConfigureContainerRegistries: configureContainerRegistries,
				DefaultAppsName:              defaultAppsName,
				DefaultAppsCatalog:           defaultAppsCatalog,
				DefaultAppsVersion:           defaultAppsVersion,
				PrivateCA:                    privateCA,
				CertManagerDNSChallenge:      certManagerDNSChallenge,
				MCCustomCoreDNSConfig:        mcCustomCoreDNSConfig,
				MCProxyEnabled:               mcProxyEnabled,
				MCHTTPSProxy:                 mcHTTPSProxy,
				TaylorBotToken:               taylorBotToken,
				Secrets: pushcmc.SecretFlags{
					SSHDeployKey: pushcmc.DeployKey{
						Passphrase: deployKeyPassphrase,
						Identity:   deployKeyIdentity,
						KnownHosts: deployKeyKnownHosts,
					},
					CustomerDeployKey: pushcmc.DeployKey{
						Passphrase: customerDeployKeyPassphrase,
						Identity:   customerDeployKeyIdentity,
						KnownHosts: customerDeployKeyKnownHosts,
					},
					SharedDeployKey: pushcmc.DeployKey{
						Passphrase: sharedDeployKeyPassphrase,
						Identity:   sharedDeployKeyIdentity,
						KnownHosts: sharedDeployKeyKnownHosts,
					},
					VSphereCredentials:                 vSphereCredentials,
					CloudDirectorCredentials:           cloudDirectorCredentials,
					AzureClusterIdentityUA:             azureClusterIdentityUA,
					AzureClusterIdentitySP:             azureClusterIdentitySP,
					AzureSecretClusterIdentityStaticSP: azureSecretClusterIdentityStaticSP,
					ContainerRegistryConfiguration:     containerRegistryConfiguration,
					ClusterValues:                      clusterValues,
					CertManagerRoute53Region:           certManagerRoute53Region,
					CertManagerRoute53Role:             certManagerRoute53Role,
					CertManagerRoute53AccessKeyID:      certManagerRoute53AccessKeyID,
					CertManagerRoute53SecretAccessKey:  certManagerRoute53SecretAccessKey,
				},
			},
		}
		err = push.Run(c, ctx)
		if err != nil {
			return err
		}
		return nil
	},
}

// pushInstallationsCmd represents the push installations command
var pushInstallationsCmd = &cobra.Command{
	Use:   "installations",
	Short: "Pushes configuration of a Management Cluster installations repository entry",
	Long: `Pushes configuration of a Management Cluster installations repository entry to
installations repository. For example:

mcli push installations --cluster=gigmac --input=cluster.yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		defaultPush()
		err := validateRoot(cmd, args)
		if err != nil {
			return err
		}
		err = validatePush(cmd, args)
		if err != nil {
			return err
		}
		ctx := context.Background()
		client := github.New(github.Config{
			Token: githubToken,
		})
		i := pushinstallations.Config{
			Cluster:             cluster,
			Github:              client,
			InstallationsBranch: installationsBranch,
			CMCRepository:       cmcRepository,
			Provider:            provider,
			Flags: pushinstallations.InstallationsFlags{
				BaseDomain: baseDomain,
				Team:       team,
				Customer:   customer,
				AWS: pushinstallations.AWSFlags{
					Region:                 awsRegion,
					InstallationAWSAccount: awsAccountID,
				},
			},
		}
		if input != "" {
			i.Input, err = installations.GetInstallationsFromFile(input)
			if err != nil {
				return fmt.Errorf("failed to get new installations object from input file.\n%w", err)
			}
		}
		installations, err := i.Run(ctx)
		if err != nil {
			return fmt.Errorf("failed to push installations.\n%w", err)
		}
		return installations.Print()
	},
}

// pushCMCCmd represents the push CMC command
var pushCMCCmd = &cobra.Command{
	Use:   "cmc",
	Short: "Pushes configuration of a Management Cluster CMC repository entry",
	Long: `Pushes configuration of a Management Cluster CMC repository entry to
CMC repository. For example:

mcli push cmc --cluster=gigmac --input=cluster.yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		defaultPush()
		err := validateRoot(cmd, args)
		if err != nil {
			return err
		}
		err = validatePush(cmd, args)
		if err != nil {
			return err
		}
		ctx := context.Background()
		client := github.New(github.Config{
			Token: githubToken,
		})
		c := pushcmc.Config{
			Cluster:       cluster,
			Github:        client,
			CMCRepository: cmcRepository,
			CMCBranch:     cmcBranch,
			Provider:      provider,
			Flags: pushcmc.CMCFlags{
				SecretFolder:                 secretFolder,
				AgePubKey:                    agePubKey,
				AgeKey:                       ageKey,
				MCAppsPreventDeletion:        mcAppsPreventDeletion,
				ClusterAppName:               clusterAppName,
				ClusterAppCatalog:            clusterAppCatalog,
				ClusterAppVersion:            clusterAppVersion,
				ClusterNamespace:             clusterNamespace,
				ConfigureContainerRegistries: configureContainerRegistries,
				DefaultAppsName:              defaultAppsName,
				DefaultAppsCatalog:           defaultAppsCatalog,
				DefaultAppsVersion:           defaultAppsVersion,
				PrivateCA:                    privateCA,
				CertManagerDNSChallenge:      certManagerDNSChallenge,
				MCCustomCoreDNSConfig:        mcCustomCoreDNSConfig,
				MCProxyEnabled:               mcProxyEnabled,
				MCHTTPSProxy:                 mcHTTPSProxy,
				TaylorBotToken:               taylorBotToken,
				Secrets: pushcmc.SecretFlags{
					SSHDeployKey: pushcmc.DeployKey{
						Passphrase: deployKeyPassphrase,
						Identity:   deployKeyIdentity,
						KnownHosts: deployKeyKnownHosts,
					},
					CustomerDeployKey: pushcmc.DeployKey{
						Passphrase: customerDeployKeyPassphrase,
						Identity:   customerDeployKeyIdentity,
						KnownHosts: customerDeployKeyKnownHosts,
					},
					SharedDeployKey: pushcmc.DeployKey{
						Passphrase: sharedDeployKeyPassphrase,
						Identity:   sharedDeployKeyIdentity,
						KnownHosts: sharedDeployKeyKnownHosts,
					},
					VSphereCredentials:                 vSphereCredentials,
					CloudDirectorCredentials:           cloudDirectorCredentials,
					AzureClusterIdentityUA:             azureClusterIdentityUA,
					AzureClusterIdentitySP:             azureClusterIdentitySP,
					AzureSecretClusterIdentityStaticSP: azureSecretClusterIdentityStaticSP,
					ContainerRegistryConfiguration:     containerRegistryConfiguration,
					ClusterValues:                      clusterValues,
					CertManagerRoute53Region:           certManagerRoute53Region,
					CertManagerRoute53Role:             certManagerRoute53Role,
					CertManagerRoute53AccessKeyID:      certManagerRoute53AccessKeyID,
					CertManagerRoute53SecretAccessKey:  certManagerRoute53SecretAccessKey,
				},
			},
		}
		if input != "" {
			c.Input, err = cmc.GetCMCFromFile(input)
			if err != nil {
				return fmt.Errorf("failed to get new CMC object from input file.\n%w", err)
			}
		}
		cmc, err := c.Run(ctx)
		if err != nil {
			return fmt.Errorf("failed to push CMC.\n%w", err)
		}
		return cmc.Print()
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.AddCommand(pushInstallationsCmd)
	pushCmd.AddCommand(pushCMCCmd)
	addFlagsPush()
}
