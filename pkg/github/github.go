package github

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/giantswarm/mcli/pkg/key"
	"github.com/google/go-github/v57/github"
	"github.com/rs/zerolog/log"
)

const (
	MaxRedirects = 10
)

type Github struct {
	*github.Client
}

type Repository struct {
	*Github
	Name         string
	Organization string
	Branch       string
	Path         string
}

type Config struct {
	Token string
}

func New(config Config) *Github {
	log.Debug().Msg("creating github client")

	return &Github{
		Client: github.NewClient(nil).WithAuthToken(config.Token),
	}
}

func (r *Repository) GetFile(ctx context.Context) (string, error) {
	// check if Organization exists
	if err := r.CheckOrganization(ctx); err != nil {
		return "", err
	}

	// check if Repository exists
	if err := r.CheckRepository(ctx); err != nil {
		return "", err
	}

	// check if Branch exists
	if err := r.CheckBranch(ctx); err != nil {
		return "", err
	}

	// get file
	log.Debug().Msg(fmt.Sprintf("getting file %s of branch %s of repository %s/%s", r.Path, r.Branch, r.Organization, r.Name))
	file, _, resp, err := r.Repositories.GetContents(ctx, r.Organization, r.Name, r.Path, &github.RepositoryContentGetOptions{
		Ref: r.Branch,
	})
	if err != nil {
		if resp.StatusCode == 404 {
			return "", fmt.Errorf("file %s of branch %s of repository %s/%s does not exist.\n%w\n%w", r.Path, r.Branch, r.Organization, r.Name, err, ErrNotFound)
		} else {
			return "", fmt.Errorf("failed to get file %s of branch %s of repository %s/%s.\n%w", r.Path, r.Branch, r.Organization, r.Name, err)
		}
	}

	// check if file is a directory
	if file.GetType() == "dir" {
		return "", fmt.Errorf("file %s of branch %s of repository %s/%s is a directory.\n%w", r.Path, r.Branch, r.Organization, r.Name, ErrInvalidFormat)
	}

	// check if file is encoded and decode it if necessary
	if enc := file.GetEncoding(); enc != "" {
		if enc == "base64" {
			log.Debug().Msg(fmt.Sprintf("decoding file %s of branch %s of repository %s/%s", r.Path, r.Branch, r.Organization, r.Name))
			decoded, err := base64.StdEncoding.DecodeString(*file.Content)
			if err != nil {
				return "", fmt.Errorf("failed to decode file %s of branch %s of repository %s/%s.\n%w", r.Path, r.Branch, r.Organization, r.Name, err)
			}
			return string(decoded), nil
		}
		return "", fmt.Errorf("file %s of branch %s of repository %s/%s is encoded in %s\n%w", r.Path, r.Branch, r.Organization, r.Name, file.GetEncoding(), ErrInvalidFormat)
	}

	return *file.Content, nil
}

func (r *Repository) CreateFile(ctx context.Context, content []byte) error {
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

	// get the SHA in case the file already exists
	fileSHA, err := r.GetFileSHA(ctx)
	if err != nil {
		return err
	}

	// get commit message
	var message string
	if fileSHA == "" {
		message = fmt.Sprintf("creating installation %s", r.Path)
	} else {
		message = fmt.Sprintf("updating installation %s", r.Path)
	}

	// create the file and the directory structure if necessary
	log.Debug().Msg(fmt.Sprintf("creating file %s of branch %s of repository %s/%s", r.Path, r.Branch, r.Organization, r.Name))
	_, _, err = r.Repositories.CreateFile(ctx, r.Organization, r.Name, r.Path, &github.RepositoryContentFileOptions{
		Message: github.String(message),
		Content: content,
		Branch:  github.String(r.Branch),
		SHA:     github.String(fileSHA),
	})
	if err != nil {
		return fmt.Errorf("failed to create file %s of branch %s of repository %s/%s.\n%w", r.Path, r.Branch, r.Organization, r.Name, err)
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

func (r *Repository) CreateBranch(ctx context.Context) error {
	// get master branch sha
	log.Debug().Msg(fmt.Sprintf("getting sha of %s branch of repository %s/%s", key.InstallationsMainBranch, r.Organization, r.Name))
	branch, _, err := r.Repositories.GetBranch(ctx, r.Organization, r.Name, key.InstallationsMainBranch, MaxRedirects)
	if err != nil {
		return fmt.Errorf("failed to get sha of %s branch of repository %s/%s.\n%w", key.InstallationsMainBranch, r.Organization, r.Name, err)
	}

	// create Branch called r.Branch from master
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

func (r *Repository) GetFileSHA(ctx context.Context) (string, error) {
	var sha string
	{
		log.Debug().Msg(fmt.Sprintf("getting sha of file %s of branch %s of repository %s/%s", r.Path, r.Branch, r.Organization, r.Name))
		file, _, resp, err := r.Repositories.GetContents(ctx, r.Organization, r.Name, r.Path, &github.RepositoryContentGetOptions{
			Ref: r.Branch,
		})
		if err != nil {
			if resp.StatusCode == 404 {
				// if the file doesn't exist, sha is empty
				sha = ""
			} else {
				return "", fmt.Errorf("failed to get sha of file %s of branch %s of repository %s/%s.\n%w", r.Path, r.Branch, r.Organization, r.Name, err)
			}
		} else {
			// if the file exists, get its sha
			sha = *file.SHA
		}
	}
	return sha, nil
}
