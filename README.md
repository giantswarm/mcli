[![CircleCI](https://dl.circleci.com/status-badge/img/gh/giantswarm/`mcli`/tree/main.svg?style=svg&circle-token=6a77fd68819c9f0d95d09d9ba329ba5e5bf169e6)](https://dl.circleci.com/status-badge/redirect/gh/giantswarm/`mcli`/tree/main)

# MCLI

CLI tool for managing Giant Swarm installations.

## What is MCLI?

`mcli` is a golang based CLI tool for managing configuration for Giant Swarm management clusters.
The idea is to provide a tool that allows for easier review, creation and updates of configuration data.
Commands found inside the tool are based on scripts that are part of [mc-bootstrap](github.com/giantswarm/mc-bootstrap) with the goal to replace more and more of them as new commands are added to `mcli`.

## Installation

`mcli` can be installed as Go binary.

```bash
go install github.com/giantswarm/mcli@latest
```

## Requirements

`mcli` uses the same environment variables and mechanisms as `mc-bootstrap`.

### GitHub

Ensure that you have a valid GitHub token set in the `GITHUB_TOKEN` environment variable.

### SOPS

For the time being, `mcli` uses the `sops` binary to encrypt and decrypt secrets in the same way as `mc-bootstrap`.
Ensure that you have [installed `sops`](https://github.com/getsops/sops) and that it is available in your PATH.

Furthermore, ensure that the `SOPS_AGE_KEY` environment variable is set before running a command.
This is always the case when running `mc-bootstrap` to create a new management cluster.
When using the tool to update an existing management cluster, make sure to set the `SOPS_AGE_KEY` environment variable to the correct value.
Please only use secure mechanisms to store and retrieve this key and do not expose it in the command line.
For example, if it is stored in LastPass, it can be set like this:

```bash
export SOPS_AGE_KEY=$(lpass show "Path/To/Secret/$CLUSTER.agekey" --notes)
```

## Usage

```bash
mcli pull -c $MC_NAME > config.yaml
```

Use `--help` to learn about more options.

## Commands

As most configuration is stored in git repositories, commands generally involve pulling from, pushing to, or creating repositories.

### `mcli pull`

Pulls the configuration of a given management cluster and prints it to stdout.
This can be used to review the configuration before making changes or to use as a base for creating a new configuration.

### `mcli push`

Pushes configuration of a management cluster. This can be used to create or update a management cluster.

### `mcli create`

Creates a repository. For the time being, this is only used to create a new cmc repository.

> [!TIP]
> The tool will not print any logs unless it is run in `--verbose` mode.

## Repositories

`mcli` uses a set of repositories to store configuration data.

`push` and `pull` commands will by default access all repositories below one by one to get the complete configuration.
However, subcommands can be used to access only a specific repository.
As of right now, the `create` command can only be called explicitly for the `cmc` repository.

### `cmc`

The customer management cluster repository `giantswarm/$CUSTOMER-management-clusters`.

This repository contains the configuration for all management clusters of a customer.
A new cmc repository is created for each customer using the `mcli create cmc` command.

The `mcli pull cmc` and `mcli push cmc` commands can be used to pull and push the configuration for a specific management cluster to the cmc repository.

> [!IMPORTANT]
> When nothing else is specified via `--customer` or `--cmc-repository` flags, the tool will use the [giantswarm-management-clusters](github.com/giantswarm/giantswarm-management-clusters) repository as the default.

### `installations`

The installations repository [`giantswarm/installations`](github.com/giantswarm/installations).

This repository mostly contains configuration for vintage installations.
However, it also contains basic information about the management cluster.
The `mcli pull installations` and `mcli push installations` commands can be used to pull and push this information to the installations repository.

## Input

The tool can be used either with an input file or flags.
To ease integration into [mc-bootstrap](github.com/giantswarm/mc-bootstrap), flags can be read from its environment variables and secret files.

### Input file

The input file is a YAML file that contains the configuration for a management cluster.
The output of an `mcli pull` command can be used as a base for an input file to the corresponding `mcli push` command.
Please note that the output of `mcli pull` will redact sensitive information like secrets.

The file can be passed to the tool with the `--input` flag.

> [!IMPORTANT]
> When using an input file via `--input`, the tool will ignore other configuration flags.

### Flags

The tool can be used with flags to provide configuration data.
To see all available flags for a command, use the `--help` flag.
Generally, when creating a new management cluster via flags, mandatory flags need to be set.
That should always be the case when running `mc-bootstrap` to create a new management cluster.

For updating a management cluster via flags, only the flags that need to be updated need to be set.

### Environment variables

The tool can read configuration from environment variables.
The available environment variables are the same as those used by [mc-bootstrap](github.com/giantswarm/mc-bootstrap).
Environment variables are used to set the default values for corresponding flags.


## Example usage

Some examples of how the tool can be used.

### Pull management cluster configuration

Pull the configuration of mc called `$CLUSTER` belonging to `$CUSTOMER` and print it to stdout.

```bash
mcli pull --cluster $CLUSTER --customer $CUSTOMER
```
The output will look something like this:

```yaml
installations:
  base: example.gigantic.io
  codename: $CLUSTER
  customer: $CUSTOMER
  cmc_repository: $CUSTOMER-management-clusters
  ccr_repository: $CUSTOMER-configs
  accountEngineer: Elmo
  pipeline: stable
  provider: capa
  aws:
    region: eu-central-1
    hostCluster:
      account: "12345"
      cloudtrailBucket: ""
      adminRoleARN: arn:aws:iam::12345:role/RoleName
      guardDuty: false
    guestCluster:
      account: "12345"
      cloudtrailBucket: ""
      guardDuty: false
cmc:
  agePubKey: age12345
  cluster: $CLUSTER
  clusterApp:
    name: $CLUSTER
    catalog: cluster
    version: 0.65.0
    values: |
        global:
            metadata:
            description: "example MC"
            name: "example"
            organization: "giantswarm"
        [...]
    appName: cluster-aws
  defaultApps:
    name: $CLUSTER-default-apps
    catalog: cluster
    version: v0.49.0
    values: |
      clusterName: $CLUSTER
      organization: giantswarm
      managementCluster: $CLUSTER
    appName: default-apps-aws
  clusterIntegratesDefaultApps: false
  mcAppsPreventDeletion: true
  privateCA: false
  privateMC: false
  clusterNamespace: org-giantswarm
  provider:
    name: capa
  taylorBotToken: REDACTED
  sshDeployKey:
    key: REDACTED
    identity: REDACTED
    knownHosts: |
      github.com ecdsa-sha2...
  customerDeployKey:
    key: REDACTED
    identity: REDACTED
    knownHosts: |
      github.com ecdsa-sha2...
  sharedDeployKey:
    key: REDACTED
    identity: REDACTED
    knownHosts: |
      github.com ecdsa-sha2...
  certManagerDNSChallenge:
    enabled: false
  configureContainerRegistries:
    enabled: true
    values: REDACTED
  customCoreDNS:
    enabled: false
  disableDenyAllNetPol: false
  mcProxy:
    enabled: false
  baseDomain: awstest.gigantic.io
  gitOps:
    cmcRepository: $CUSTOMER-management-clusters
    cmcBranch: $CLUSTER_auto_branch
    mcbBranchSource: main
    configBranch: $CLUSTER_auto_config
    mcAppCollectionBranch: $CLUSTER_auto_branch
```

Alternatively, the `cmc` and `installations` repositories can be accessed separately.

```bash
mcli pull cmc --cluster $CLUSTER --customer $CUSTOMER
```

```bash
mcli pull installations --cluster $CLUSTER
```
If needed, the REDACTED secrets can be displayed in the output by setting the `--display-secrets` flag.
> [!WARNING]
> It's important to be careful with the `--display-secrets` flag as it will print sensitive information in the output. Only use it when necessary.

### Create management cluster repository

Create a new management cluster repository for `$CUSTOMER`

```bash
mcli create cmc --customer $CUSTOMER -v
```

> [!TIP]
> The `-v` flag is used to print logs to stdout.
>

The result will be a new repository in the `giantswarm` organization called `$CUSTOMER-management-clusters`.
Also, a pull request in the `giantswarm/github` repository will be created to add the new repository.

### Push management cluster configuration

In order to push the configuration of a new management cluster, a couple of secrets need to be created first.
Usually, this will be done in the process of running `mc-bootstrap` to create a new management cluster.

- To create entirely new configuration for testing or otherwise, a new age key pair can be generated via the `age-keygen` command.

- More information on the various secrets and values can be found in the [mc-bootstrap repository](github.com/giantswarm/mc-bootstrap).
- The output of the `mcli pull` command can be used as a base for the input file.

_Below examples assume that the cluster either already exists and we are updating it or that the needed secrets are already created and present in the input file or secret folder._

#### With an input file

Push the configuration of mc called `$CLUSTER` belonging to `$CUSTOMER` to repositories.

```bash
mcli push --cluster $CLUSTER --customer $CUSTOMER --input config.yaml
```

A succesful push should print the resulting configuration to stdout in the same way as the `pull` command.
It should also create a new branch in the affected repositories.

Alternatively, the `cmc` and `installations` repositories can be accessed separately.

```bash
mcli push cmc --cluster $CLUSTER --customer $CUSTOMER --input cmc-config.yaml
```

```bash
mcli push installations --cluster $CLUSTER --customer --input installations-config.yaml
```

#### With flags

> [!WARNING]
> We recommend the use of an input file when working with mcli outside of mc-bootstrap since it is more reliable, less error-prone and easier to manage.

However, flags or environment variables can be used to push configuration as well. Please refer to the [Flags and environment variables](#flags-and-environment-variables) section for more information on available flags and environment variables.

Push the configuration of mc called `$CLUSTER` to the installations repository.

```bash
mcli push installations -c $CLUSTER --customer giantswarm --provider capa --base-domain example.gigantic.io --team bigmac --aws-region eu-central-1 --aws-account-id 12345
```

Update the configuration of mc called `$CLUSTER` belonging to `$CUSTOMER` in the cmc repository with new values.
Here, cert-manager-dns-challenge is enabled.

```bash
mcli push cmc --cluster $CLUSTER --customer $CUSTOMER --cert-manager-dns-challenge --secret-folder secret/folder
```

> [!NOTE]
> The `--secret-folder` flag is used here to specify the folder where the secrets are stored. 
> When using the tool with flags the expectation is that the secrets are stored in a folder and the tool will read them from there.
> More information on the secret folder can be found in the [mc-bootstrap repository](github.com/giantswarm/mc-bootstrap).

## Reference

### Flags and environment variables

| Command | Flag | Environment variable | Description | Note |
| --- | --- | --- | --- | --- |
| all | `--cluster`, `-c` | `INSTALLATION` | The name of the management cluster. |
|  | `--customer` | `CUSTOMER` | The name of the customer. | Defaults to "giantswarm"
|  | `--cmc-branch` | `CMC_BRANCH` | The branch of the cmc repository to use. | Defaults to "main" in pull case and auto naming in push case
|  | `--cmc-repository` | `CMC_REPOSITORY` | The name of the cmc repository to use. | Defaults to "giantswarm-management-clusters"
|  | `--input` | | The path to the input file. |
|  | `--verbose`, `-v` | | Print logs to stdout. |
|  | `--display-secrets` |  | Display secrets in the output. |
|  | `--github-token` | `GITHUB_TOKEN` | The GitHub token to use. |
|  | `--skip` | | Repositories to skip. |
|  | `--installations-branch` | `INSTALLATIONS_BRANCH` | The branch of the installations repository to use. | Defaults to "master" in pull case and auto naming in push case
|  |  |  |  |
| `push` | `--provider` | `PROVIDER` | The provider of the management cluster. |
|  | `--base-domain` | `BASE_DOMAIN` | The base domain of the management cluster. |
| `push installations` | `--team` | `TEAM_NAME` | The team name of the management cluster. |
|  | `--aws-region` | `AWS_REGION` | The AWS region of the management cluster. |
|  | `--aws-account-id` | `INSTALLATION_AWS_ACCOUNT` | The AWS account ID of the management cluster. |
|  | `--ccr-repository` | `CCR_REPOSITORY` | The name of the ccr repository to use. |
|  | `--pipeline` | `MC_PIPELINE` | The pipeline to use for the installation. Defaults to "testing" |
| `push cmc` | `--mc-apps-prevent-deletion` | `MC_APPS_PREVENT_DELETION` | Prevent deletion of mc apps. |
|  | `--cluster-app-name` | `CLUSTER_APP_NAME` | The name of the cluster app. |
|  | `--cluster-app-catalog` | `CLUSTER_APP_CATALOG` | The catalog of the cluster app. |
|  | `--cluster-app-version` | `CLUSTER_APP_VERSION` | The version of the cluster app. |
|  | `--cluster-integrates-default-apps` | `CLUSTER_INTEGRATES_DEFAULT_APPS` | Default apps are integrated into the cluster app. |
|  | `--cluster-namespace` | `CLUSTER_NAMESPACE` | The namespace of the management cluster. |
|  | `--configure-container-registries` | `CONFIGURE_CONTAINER_REGISTRIES` | Configure container registries. |
|  | `--default-apps-name` | `DEFAULT_APPS_APP_NAME` | The name of the default apps. |
|  | `--default-apps-catalog` | `DEFAULT_APPS_APP_CATALOG` | The catalog of the default apps. |
|  | `--default-apps-version` | `DEFAULT_APPS_APP_VERSION` | The version of the default apps. |
|  | `--private-ca` | `PRIVATE_CA` | Use a private CA. |
|  | `--private-mc` | `MC_PRIVATE` | The management cluster is private. |
|  | `--cert-manager-dns-challenge` | `CERT_MANAGER_DNS01_CHALLENGE` | Use cert-manager DNS challenge. |
|  | `--mc-custom-coredns-config` | `MC_CUSTOM_COREDNS_CONFIG` | Use custom CoreDNS config. |
|  | `--mc-proxy-enabled` | `MC_PROXY_ENABLED` | Use mc proxy. |
|  | `--mc-https-proxy` | `MC_HTTPS_PROXY` | Use mc https proxy. |
|  | `--age-pub-key` | `AGE_PUBKEY` | The age public key. |
|  | `--mcb-branch-source` | `MCB_BRANCH_SOURCE` | The source branch of the mcb repository to use. | Defaults to "main"
|  | `--config-branch` | `CONFIG_BRANCH` | The branch of the config repository to use. | Defaults to auto naming
|  | `--mc-app-collection-branch` | `MC_APP_COLLECTION_BRANCH` | The branch of the mc-app-collection repository to use. | Defaults to auto naming
|  | `--registry-domain` | `REGISTRY_DOMAIN` | Custom registry domain. |
|  | `--taylor-bot-token` |  | The Taylor bot token. | overrides value from secret files
|  | `--deploy-key` |  | The deploy key passphrase. |  overrides value from secret files
|  | `--deploy-key-identity` |  | The deploy key identity. | overrides value from secret files
|  | `--deploy-key-known-hosts` |  | The deploy key known hosts. | overrides value from secret files
|  | `--customer-deploy-key` |  | The customer deploy key passphrase. | overrides value from secret files
|  | `--customer-deploy-key-identity` |  | The customer deploy key identity. | overrides value from secret files
|  | `--customer-deploy-key-known-hosts` |  | The customer deploy key known hosts. | overrides value from secret files
|  | `--shared-deploy-key` |  | The shared deploy key passphrase. | overrides value from secret files
|  | `--shared-deploy-key-identity` |  | The shared deploy key identity. | overrides value from secret files
|  | `--shared-deploy-key-known-hosts` |  | The shared deploy key known hosts. | overrides value from secret files
|  | `--vsphere-credentials` |  | The vSphere credentials. | overrides value from secret files
|  | `--cloud-director-refresh-token` |  | The cloud director refresh token. | overrides value from secret files
|  | `--azure-ua-client-id` |  | The Azure UA client ID. | overrides value from secret files
|  | `--azure-ua-tenant-id` |  | The Azure UA tenant ID. | overrides value from secret files
|  | `--azure-ua-resource-id` |  | The Azure UA resource ID. | overrides value from secret files
|  | `--azure-client-id` | | The Azure client ID. | overrides value from secret files
|  | `--azure-tenant-id` |  | The Azure tenant ID. | overrides value from secret files
|  | `--azure-client-secret` |  | The Azure client secret. | overrides value from secret files
|  | `--azure-subscription-id` |  | The Azure subscription ID. | overrides value from secret files
|  | `--container-registry-configuration` |  | The container registry configuration. | overrides value from secret files
|  | `--cluster-values` |  | The cluster values. | overrides value from secret files
|  | `--cert-manager-route53-region` |  | The cert-manager Route53 region. | overrides value from secret files
|  | `--cert-manager-route53-role` |  | The cert-manager Route53 role. | overrides value from secret files
|  | `--cert-manager-route53-access-key-id` |  | The cert-manager Route53 access key ID. | overrides value from secret files
|  | `--cert-manager-route53-secret-access-key` |  | The cert-manager Route53 secret access key. | overrides value from secret files
|  |  |  |  |
