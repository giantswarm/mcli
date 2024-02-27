package github

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/google/go-github/v57/github"
	"github.com/rs/zerolog/log"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

const (
	MaxRedirects = 10
)

type Github struct {
	*github.Client
	Graph *githubv4.Client
}

type Repository struct {
	*Github
	Name         string
	Organization string
	Branch       string
}

type Config struct {
	Token string
}

func New(config Config) *Github {
	log.Debug().Msg("creating github client")
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: config.Token})
	httpClient := oauth2.NewClient(context.Background(), src)

	return &Github{
		Client: github.NewClient(httpClient),
		Graph:  githubv4.NewClient(httpClient),
	}
}

func (r *Repository) Check(ctx context.Context) error {
	// check if Organization exists
	if err := r.CheckOrganization(ctx); err != nil {
		return err
	}

	// check if Repository exists
	if err := r.CheckRepository(ctx); err != nil {
		return err
	}

	// check if Branch exists
	if err := r.CheckBranch(ctx); err != nil {
		return err
	}

	return nil
}

func (r *Repository) getContents(ctx context.Context, path string) (fileContent *github.RepositoryContent, directoryContent []*github.RepositoryContent, resp *github.Response, err error) {
	log.Debug().Msg(fmt.Sprintf("getting contents %s of branch %s of repository %s/%s", path, r.Branch, r.Organization, r.Name))
	return r.Repositories.GetContents(ctx, r.Organization, r.Name, path, &github.RepositoryContentGetOptions{
		Ref: r.Branch,
	})
}

func (r *Repository) GetFile(ctx context.Context, path string) (string, error) {
	log.Debug().Msg(fmt.Sprintf("getting file %s of branch %s of repository %s/%s", path, r.Branch, r.Organization, r.Name))
	file, _, resp, err := r.getContents(ctx, path)
	if err != nil {
		if resp.StatusCode == 404 {
			return "", fmt.Errorf("file %s of branch %s of repository %s/%s does not exist.\n%w\n%w", path, r.Branch, r.Organization, r.Name, err, ErrNotFound)
		} else {
			return "", fmt.Errorf("failed to get file %s of branch %s of repository %s/%s.\n%w", path, r.Branch, r.Organization, r.Name, err)
		}
	}
	return r.GetStringFromFile(file)
}

func (r *Repository) GetDirectory(ctx context.Context, path string) (map[string]string, error) {
	log.Debug().Msg(fmt.Sprintf("getting directory %s of branch %s of repository %s/%s", path, r.Branch, r.Organization, r.Name))
	_, directory, resp, err := r.getContents(ctx, path)
	if err != nil {
		if resp.StatusCode == 404 {
			return nil, fmt.Errorf("directory %s of branch %s of repository %s/%s does not exist.\n%w\n%w", path, r.Branch, r.Organization, r.Name, err, ErrNotFound)
		} else {
			return nil, fmt.Errorf("failed to get directory %s of branch %s of repository %s/%s.\n%w", path, r.Branch, r.Organization, r.Name, err)
		}
	}

	//ensure directory is not nil
	if directory == nil {
		return nil, fmt.Errorf("directory %s of branch %s of repository %s/%s is nil.\n%w", path, r.Branch, r.Organization, r.Name, ErrInvalidFormat)
	}

	// check if directory contains files
	if len(directory) == 0 {
		return nil, fmt.Errorf("directory %s of branch %s of repository %s/%s is empty.\n%w", path, r.Branch, r.Organization, r.Name, ErrNotFound)
	}

	// create a map of files
	files := make(map[string]string)
	for _, file := range directory {
		if file.GetType() == "dir" {
			// if file is a directory, get its contents
			dir, err := r.GetDirectory(ctx, file.GetPath())
			if err != nil {
				return nil, err
			}
			for k, v := range dir {
				files[k] = v
			}
		} else {
			contents, err := r.GetFile(ctx, file.GetPath())
			if err != nil {
				return nil, err
			}
			files[file.GetPath()] = contents
		}
	}
	return files, nil
}

func (r *Repository) GetStringFromFile(file *github.RepositoryContent) (string, error) {
	//ensure file is not nil
	if file == nil {
		return "", fmt.Errorf("file of branch %s of repository %s/%s is nil.\n%w", r.Branch, r.Organization, r.Name, ErrInvalidFormat)
	}

	// check if file is a directory
	if file.GetType() == "dir" {
		return "", fmt.Errorf("file %s of branch %s of repository %s/%s is a directory.\n%w", file.GetPath(), r.Branch, r.Organization, r.Name, ErrInvalidFormat)
	}

	// check if file is encoded and decode it if necessary
	if enc := file.GetEncoding(); enc != "" {
		if enc == "base64" {
			log.Debug().Msg(fmt.Sprintf("decoding file %s of branch %s of repository %s/%s", file.GetPath(), r.Branch, r.Organization, r.Name))
			decoded, err := base64.StdEncoding.DecodeString(*file.Content)
			if err != nil {
				return "", fmt.Errorf("failed to decode file %s of branch %s of repository %s/%s.\n%w", file.GetPath(), r.Branch, r.Organization, r.Name, err)
			}
			return string(decoded), nil
		}
		return "", fmt.Errorf("file %s of branch %s of repository %s/%s is encoded in %s\n%w", file.GetPath(), r.Branch, r.Organization, r.Name, file.GetEncoding(), ErrInvalidFormat)
	}
	content, err := file.GetContent()
	if err != nil {
		return "", fmt.Errorf("failed to get content of file %s of branch %s of repository %s/%s.\n%w", file.GetPath(), r.Branch, r.Organization, r.Name, err)
	}
	return content, nil
}

func (r *Repository) CreateFile(ctx context.Context, content []byte, path string) error {
	if err := r.Check(ctx); err != nil {
		return err
	}

	// get the SHA in case the file already exists
	fileSHA, err := r.GetFileSHA(ctx, path)
	if err != nil {
		return err
	}

	// get commit message
	var message string
	if fileSHA == "" {
		message = fmt.Sprintf("creating %s", path)
	} else {
		message = fmt.Sprintf("updating %s", path)
	}

	// create the file and the directory structure if necessary
	log.Debug().Msg(fmt.Sprintf("creating file %s of branch %s of repository %s/%s", path, r.Branch, r.Organization, r.Name))
	_, _, err = r.Repositories.CreateFile(ctx, r.Organization, r.Name, path, &github.RepositoryContentFileOptions{
		Message: github.String(message),
		Content: content,
		Branch:  github.String(r.Branch),
		SHA:     github.String(fileSHA),
	})
	if err != nil {
		return fmt.Errorf("failed to create file %s of branch %s of repository %s/%s.\n%w", path, r.Branch, r.Organization, r.Name, err)
	}

	return nil
}

func (r *Repository) CreateDirectory(ctx context.Context, path string, content map[string]string) error {
	// create the files
	for k, v := range content {
		if err := r.CreateFile(ctx, []byte(v), fmt.Sprintf("%s/%s", path, k)); err != nil {
			return fmt.Errorf("failed to create file %s of directory %s of branch %s of repository %s/%s.\n%w", k, path, r.Branch, r.Organization, r.Name, err)
		}
	}
	return nil
}

func (r *Repository) CheckOrganization(ctx context.Context) error {
	// check if Organization exists
	log.Debug().Msg(fmt.Sprintf("checking if organization %s exists", r.Organization))
	_, resp, err := r.Organizations.Get(ctx, r.Organization)
	if err != nil {
		if resp.StatusCode == 404 {
			return fmt.Errorf("organization %s does not exist.\n%w\n%w", r.Organization, err, ErrNotFound)
		} else {
			return fmt.Errorf("failed to get organization %s.\n%w", r.Organization, err)
		}
	}
	return nil
}

func (r *Repository) CheckRepository(ctx context.Context) error {
	// check if Repository exists
	log.Debug().Msg(fmt.Sprintf("checking if repository %s/%s exists", r.Organization, r.Name))
	_, resp, err := r.Repositories.Get(ctx, r.Organization, r.Name)
	if err != nil {
		if resp.StatusCode == 404 {
			return fmt.Errorf("repository %s/%s does not exist.\n%w\n%w", r.Organization, r.Name, err, ErrNotFound)
		} else {
			return fmt.Errorf("failed to get repository %s/%s.\n%w", r.Organization, r.Name, err)
		}
	}
	return nil
}

func (r *Repository) CheckBranch(ctx context.Context) error {
	// check if Branch exists
	log.Debug().Msg(fmt.Sprintf("checking if branch %s of repository %s/%s exists", r.Branch, r.Organization, r.Name))
	_, resp, err := r.Repositories.GetBranch(ctx, r.Organization, r.Name, r.Branch, MaxRedirects)
	if err != nil {
		if resp.StatusCode == 404 {
			return fmt.Errorf("branch %s of repository %s/%s does not exist.\n%w\n%w", r.Branch, r.Organization, r.Name, err, ErrNotFound)
		} else {
			return fmt.Errorf("failed to get branch %s of repository %s/%s.\n%w", r.Branch, r.Organization, r.Name, err)
		}
	}
	return nil
}

func (r *Repository) CreateBranch(ctx context.Context, mainbranch string) error {
	// get main branch sha
	log.Debug().Msg(fmt.Sprintf("getting sha of %s branch of repository %s/%s", mainbranch, r.Organization, r.Name))
	branch, _, err := r.Repositories.GetBranch(ctx, r.Organization, r.Name, mainbranch, MaxRedirects)
	if err != nil {
		return fmt.Errorf("failed to get sha of %s branch of repository %s/%s.\n%w", mainbranch, r.Organization, r.Name, err)
	}

	// create Branch called r.Branch from main
	log.Debug().Msg(fmt.Sprintf("creating branch %s of repository %s/%s", r.Branch, r.Organization, r.Name))
	reference := &github.Reference{
		Ref: github.String(fmt.Sprintf("refs/heads/%s", r.Branch)),
		Object: &github.GitObject{
			SHA: branch.Commit.SHA,
		},
	}
	_, _, err = r.Git.CreateRef(ctx, r.Organization, r.Name, reference)
	if err != nil {
		return fmt.Errorf("failed to create branch %s of repository %s/%s.\n%w", r.Branch, r.Organization, r.Name, err)
	}
	return nil
}

func (r *Repository) GetFileSHA(ctx context.Context, path string) (string, error) {
	var sha string
	{
		log.Debug().Msg(fmt.Sprintf("getting sha of file %s of branch %s of repository %s/%s", path, r.Branch, r.Organization, r.Name))
		file, _, resp, err := r.Repositories.GetContents(ctx, r.Organization, r.Name, path, &github.RepositoryContentGetOptions{
			Ref: r.Branch,
		})
		if err != nil {
			if resp.StatusCode == 404 {
				// if the file doesn't exist, sha is empty
				sha = ""
			} else {
				return "", fmt.Errorf("failed to get sha of file %s of branch %s of repository %s/%s.\n%w", path, r.Branch, r.Organization, r.Name, err)
			}
		} else {
			// if the file exists, get its sha
			sha = *file.SHA
		}
	}
	return sha, nil
}

func (r *Repository) CreatePrivateRepo(ctx context.Context, description string, template string) error {
	// check if Organization exists
	if err := r.CheckOrganization(ctx); err != nil {
		return err
	}

	// create a private repository from the template
	log.Debug().Msg(fmt.Sprintf("creating repository %s/%s", r.Organization, r.Name))
	_, _, err := r.Repositories.CreateFromTemplate(ctx, r.Organization, template, &github.TemplateRepoRequest{
		Owner:       github.String(r.Organization),
		Name:        github.String(r.Name),
		Description: github.String(description),
		Private:     github.Bool(true),
	})
	if err != nil {
		return fmt.Errorf("failed to create repository %s/%s.\n%w", r.Organization, r.Name, err)
	}
	return nil
}

func (r *Repository) AddCollaborator(ctx context.Context, slug string, permission string) error {
	collaborator := fmt.Sprintf("@%s/%s", r.Organization, slug)

	// check if collaborator is already a collaborator
	log.Debug().Msg(fmt.Sprintf("checking if %s is already a collaborator of repository %s/%s", collaborator, r.Organization, r.Name))
	result, resp, err := r.Teams.IsTeamRepoBySlug(ctx, r.Organization, slug, r.Organization, r.Name)
	if err != nil {
		if resp.StatusCode != 404 {
			return fmt.Errorf("failed to check if %s is already a collaborator of repository %s/%s.\n%w", collaborator, r.Organization, r.Name, err)
		}
	}
	if result != nil {
		if result.Permissions[permission] {
			return nil
		}
	}

	log.Debug().Msg(fmt.Sprintf("adding %s as collaborator to repository %s/%s", collaborator, r.Organization, r.Name))
	_, err = r.Teams.AddTeamRepoBySlug(ctx, r.Organization, slug, r.Organization, r.Name, &github.TeamAddTeamRepoOptions{
		Permission: permission,
	})
	if err != nil {
		return fmt.Errorf("failed to add %s as collaborator to repository %s/%s.\n%w", collaborator, r.Organization, r.Name, err)
	}
	return nil
}

func (r *Repository) SetupBranchProtection(ctx context.Context, branch string) error {
	log.Debug().Msg(fmt.Sprintf("protecting branch %s of repository %s/%s", branch, r.Organization, r.Name))

	var query struct {
		Repository struct {
			ID githubv4.ID
		} `graphql:"repository(owner: $owner, name: $name)"`
	}
	variables := map[string]interface{}{
		"owner": githubv4.String(r.Organization),
		"name":  githubv4.String(r.Name),
	}
	err := r.Github.Graph.Query(ctx, &query, variables)
	if err != nil {
		return fmt.Errorf("failed to get id of repository %s/%s.\n%w", r.Organization, r.Name, err)
	}
	id := query.Repository.ID

	// check if branch protection rule already exists
	var check struct {
		Repository struct {
			BranchProtectionRules struct {
				Nodes []struct {
					Pattern githubv4.String
				}
			} `graphql:"branchProtectionRules(first: 100)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}
	err = r.Github.Graph.Query(ctx, &check, variables)
	if err != nil {
		return fmt.Errorf("failed to check if branch protection rule for branch %s of repository %s/%s already exists.\n%w", branch, r.Organization, r.Name, err)
	}
	for _, rule := range check.Repository.BranchProtectionRules.Nodes {
		if rule.Pattern == githubv4.String(branch) {
			log.Debug().Msg(fmt.Sprintf("branch protection rule for branch %s of repository %s/%s already exists", branch, r.Organization, r.Name))
			return nil
		}
	}

	var mutation struct {
		CreateBranchProtectionRule struct {
			BranchProtectionRule struct {
				ID githubv4.ID
			}
		} `graphql:"createBranchProtectionRule(input: $input)"`
	}
	input := githubv4.CreateBranchProtectionRuleInput{
		RepositoryID:                 id,
		Pattern:                      githubv4.String(branch),
		RequiresApprovingReviews:     githubv4.NewBoolean(githubv4.Boolean(true)),
		RequiredApprovingReviewCount: githubv4.NewInt(githubv4.Int(1)),
		RequiresStatusChecks:         githubv4.NewBoolean(githubv4.Boolean(true)),
		RequiresStrictStatusChecks:   githubv4.NewBoolean(githubv4.Boolean(true)),
		IsAdminEnforced:              githubv4.NewBoolean(githubv4.Boolean(true)),
	}
	err = r.Github.Graph.Mutate(ctx, &mutation, input, nil)
	if err != nil {
		return fmt.Errorf("failed to protect branch %s of repository %s/%s.\n%w", branch, r.Organization, r.Name, err)
	}
	return nil
}

func (r *Repository) CreatePullRequest(ctx context.Context, title string, base string) error {
	log.Debug().Msg(fmt.Sprintf("creating pull request to merge %s into %s of repository %s/%s", r.Branch, base, r.Organization, r.Name))
	_, _, err := r.PullRequests.Create(ctx, r.Organization, r.Name, &github.NewPullRequest{
		Title: github.String(title),
		Body:  github.String("Pull request was generated by: 'mcli'."),
		Head:  github.String(r.Branch),
		Base:  github.String(base),
		Draft: github.Bool(true),
	})
	if err != nil {
		return fmt.Errorf("failed to create pull request to merge %s into %s of repository %s/%s.\n%w", r.Branch, base, r.Organization, r.Name, err)
	}
	return nil
}
