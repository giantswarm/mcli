[![CircleCI](https://dl.circleci.com/status-badge/img/gh/giantswarm/`mcli`/tree/main.svg?style=svg&circle-token=6a77fd68819c9f0d95d09d9ba329ba5e5bf169e6)](https://dl.circleci.com/status-badge/redirect/gh/giantswarm/`mcli`/tree/main)

# MCLI

CLI tool for managing Giant Swarm installations.

## What is MCLI?

`mcli` is a golang based CLI tool for managing configuration for Giant Swarm management clusters.
The idea is to provide a tool that allows for easier review, creation and updates of configuration data.
Commands found inside the tool are based on scripts that are part of [mc-bootstrap](github.com/giantswarm/mc-bootstrap) with the goal to replace more and more of them as new commands are added to `mcli`.

## Installation

`mcli` can be installed as Go binary.

```nohighlight
go install github.com/giantswarm/mcli@latest
```

## Usage

```nohighlight
$ mcli pull -c $MC_NAME > config.yaml
```

Use `--help` to learn about more options.

## Commands

As most configuration is stored in git repositories, commands generally involve pulling from, pushing to, or creating repositories.
If no flags or input configuration is provided, the tool will use the same environment variables as mc-bootstrap.

- `mcli pull` - Pulls the configuration of a given management cluster and prints it to stdout.
This can be used to review the configuration before making changes or to use as a base for creating a new configuration.
- `mcli push` - Pushes configuration of a management cluster.
- `mcli create` - Creates a repository. For the time being, this is only needed to create a new customer repository.

The tool will not print any logs unless it is run in `--verbose` mode.
