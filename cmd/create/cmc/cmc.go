package createcmc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/giantswarm/mcli/pkg/github"
	"github.com/giantswarm/mcli/pkg/key"
	"github.com/giantswarm/mcli/pkg/repositories"
	"github.com/rs/zerolog/log"
)

const (
	kustomizationPostBuild = "bases/patches/kustomization-post-build.yaml"
	makeFile               = "Makefile.custom.mk"
	ownershipFile          = "repositories/team-honeybadger.yaml"
	customerKey            = "CUSTOMER_CODENAME"
)

type Config struct {
	Github        *github.Github
	CMCRepository string
	Customer      string
	CMCBranch     string
}

func (c *Config) Run(ctx context.Context) (*github.Repository, error) {

	// Create customer repository
	cmc, err := c.createCMC(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create CMC repository.\n%w", err)
	}

	// Add custom changes on top of repository template
	if err := c.customizeMC(ctx, cmc); err != nil {
		return nil, fmt.Errorf("failed to customize CMC repository.\n%w", err)
	}

	// Setup branch protection rules
	if err := c.setupBranchProtection(ctx, cmc); err != nil {
		return nil, fmt.Errorf("failed to setup branch protection rules.\n%w", err)
	}

	if err := c.createOwnershipPR(ctx); err != nil {
		return nil, fmt.Errorf("failed to create ownership PR.\n%w", err)
	}

	return cmc, nil
}

func (c *Config) createCMC(ctx context.Context) (*github.Repository, error) {
	cmcRepository := github.Repository{
		Github:       c.Github,
		Name:         c.CMCRepository,
		Organization: key.OrganizationGiantSwarm,
		Branch:       c.CMCBranch,
	}
	if err := cmcRepository.CheckRepository(ctx); err != nil {
		if !github.IsNotFound(err) {
			return nil, fmt.Errorf("failed to check CMC repository %s.\n%w", c.CMCRepository, err)
		} else {
			// CMC repository does not exist, create it
			description := fmt.Sprintf("Management Clusters configuration for %s", c.Customer)
			log.Debug().Msgf("creating CMC repository %s", c.CMCRepository)
			if err := cmcRepository.CreatePrivateRepo(ctx, description, key.CMCTemplateRepository); err != nil {
				return nil, fmt.Errorf("failed to create CMC repository %s.\n%w", c.CMCRepository, err)
			}
			log.Debug().Msgf("waiting for CMC repository %s to be ready", c.CMCRepository)
			time.Sleep(3 * time.Second)
		}
	} else {
		// CMC repository already exists, nothing to do
		log.Debug().Msgf("CMC repository %s already exists", c.CMCRepository)
	}

	if err := cmcRepository.AddCollaborator(ctx, key.Employees, "admin"); err != nil {
		return nil, fmt.Errorf("failed to add collaborator %s to CMC repository %s.\n%w", key.Employees, c.CMCRepository, err)
	}
	if err := cmcRepository.AddCollaborator(ctx, key.Bots, "push"); err != nil {
		return nil, fmt.Errorf("failed to add collaborator %s to CMC repository %s.\n%w", key.Bots, c.CMCRepository, err)
	}
	return &cmcRepository, nil
}

func (c *Config) customizeMC(ctx context.Context, cmcRepository *github.Repository) error {
	// Add custom changes on top of repository template
	err := cmcRepository.CheckBranch(ctx)
	if err != nil {
		if github.IsNotFound(err) {
			log.Debug().Msgf("CMC branch %s not found, creating it", c.CMCBranch)
			err = cmcRepository.CreateBranch(ctx, key.CMCMainBranch)
			if err != nil {
				return fmt.Errorf("failed to create CMC branch %s.\n%w", c.CMCBranch, err)
			}
		} else {
			return fmt.Errorf("failed to check CMC branch %s.\n%w", c.CMCBranch, err)
		}
	}
	// Get kustomization and makefile from repository
	kustomization, err := cmcRepository.GetFile(ctx, kustomizationPostBuild)
	if err != nil {
		return fmt.Errorf("failed to get kustomization file %s.\n%w", kustomizationPostBuild, err)
	}
	makefile, err := cmcRepository.GetFile(ctx, makeFile)
	if err != nil {
		return fmt.Errorf("failed to get makefile %s.\n%w", makeFile, err)
	}

	// update the kustomization if customer key is present
	if strings.Contains(kustomization, customerKey) {
		log.Debug().Msgf("Updating %s with customer codename %s", kustomizationPostBuild, c.Customer)
		kustomization = strings.ReplaceAll(kustomization, customerKey, c.Customer)
		if err := cmcRepository.CreateFile(ctx, []byte(kustomization), kustomizationPostBuild); err != nil {
			return fmt.Errorf("failed to update kustomization file %s.\n%w", kustomizationPostBuild, err)
		}
	}

	// update the makefile if customer key is present
	if strings.Contains(makefile, customerKey) {
		log.Debug().Msgf("Updating %s with customer codename %s", makeFile, c.Customer)
		makefile = strings.ReplaceAll(makefile, customerKey, c.Customer)
		if err := cmcRepository.CreateFile(ctx, []byte(makefile), makeFile); err != nil {
			return fmt.Errorf("failed to update makefile %s.\n%w", makeFile, err)
		}
	}
	return nil
}

func (c *Config) setupBranchProtection(ctx context.Context, cmcRepository *github.Repository) error {
	// Setup branch protection rules
	err := cmcRepository.SetupBranchProtection(ctx, key.CMCMainBranch)
	if err != nil {
		return fmt.Errorf("failed to setup branch protection rules for CMC branch %s.\n%w", c.CMCBranch, err)
	}
	return nil
}

func (c *Config) createOwnershipPR(ctx context.Context) error {
	githubRepository := github.Repository{
		Github:       c.Github,
		Name:         key.RepositoryGithub,
		Organization: key.OrganizationGiantSwarm,
		Branch:       key.GetOwnershipBranch(c.Customer),
	}
	err := githubRepository.CheckBranch(ctx)
	if err != nil {
		if github.IsNotFound(err) {
			log.Debug().Msgf("Ownership branch %s not found, creating it", githubRepository.Branch)
			err = githubRepository.CreateBranch(ctx, key.CMCMainBranch)
			if err != nil {
				return fmt.Errorf("failed to create ownership branch %s.\n%w", githubRepository.Branch, err)
			}
		} else {
			return fmt.Errorf("failed to check ownership branch %s.\n%w", githubRepository.Branch, err)
		}
	}
	// get ownership file from repository
	log.Debug().Msgf("getting ownership file %s from repository %s", ownershipFile, key.RepositoryGithub)
	file, err := githubRepository.GetFile(ctx, ownershipFile)
	if err != nil {
		return fmt.Errorf("failed to get ownership file %s.\n%w", ownershipFile, err)
	}

	// check if an entry for the cmc repository already exists
	if strings.Contains(file, c.CMCRepository) {
		log.Debug().Msgf("CMC repository %s already exists in ownership file", c.CMCRepository)
		return nil
	}
	repos, err := repositories.GetRepos([]byte(file))
	if err != nil {
		return fmt.Errorf("failed to get repositories from ownership file.\n%w", err)
	}
	repository := repositories.Repo{
		Name:          c.CMCRepository,
		ComponentType: "configuration",
		Gen: repositories.Gen{
			Flavours: []string{"generic"},
			Language: "generic",
		},
		Replace: map[string]bool{
			"dependabotRemove": true},
	}
	repos = append(repos, repository)
	repos = repositories.SortReposAlphabetically(repos)
	data, err := repositories.GetData(repos)
	if err != nil {
		return fmt.Errorf("failed to marshal repositories.\n%w", err)
	}

	log.Debug().Msgf("updating ownership file %s with CMC repository %s", ownershipFile, c.CMCRepository)
	if err := githubRepository.CreateFile(ctx, data, ownershipFile); err != nil {
		return fmt.Errorf("failed to update ownership file %s.\n%w", ownershipFile, err)
	}

	log.Debug().Msgf("creating PR for ownership file %s", ownershipFile)
	// create PR
	if err := githubRepository.CreatePullRequest(ctx, fmt.Sprintf("Add %s to to team honeybadger", c.CMCRepository), key.CMCMainBranch); err != nil {
		return fmt.Errorf("failed to create PR for ownership file.\n%w", err)
	}
	return nil
}
