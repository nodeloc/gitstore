package services

import (
	"context"
	"fmt"

	"github.com/google/go-github/v57/github"
	"github.com/nodeloc/git-store/internal/config"
	"golang.org/x/oauth2"
)

type GitHubService struct {
	config *config.Config
	client *github.Client
}

func NewGitHubService(cfg *config.Config) *GitHubService {
	// Create OAuth2 token source from Personal Access Token
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.GitHubAdminToken},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

	return &GitHubService{
		config: cfg,
		client: client,
	}
}

// AddRepositoryCollaborator adds a user as a collaborator to a repository
func (s *GitHubService) AddRepositoryCollaborator(ctx context.Context, owner, repo, username, permission string) error {
	opts := &github.RepositoryAddCollaboratorOptions{
		Permission: permission, // "pull", "push", "admin", "maintain", "triage"
	}

	_, _, err := s.client.Repositories.AddCollaborator(ctx, owner, repo, username, opts)
	if err != nil {
		return fmt.Errorf("failed to add collaborator: %w", err)
	}

	return nil
}

// RemoveRepositoryCollaborator removes a user as a collaborator from a repository
func (s *GitHubService) RemoveRepositoryCollaborator(ctx context.Context, owner, repo, username string) error {
	_, err := s.client.Repositories.RemoveCollaborator(ctx, owner, repo, username)
	if err != nil {
		return fmt.Errorf("failed to remove collaborator: %w", err)
	}

	return nil
}

// ListRepositories lists all repositories for the authenticated user or organization
func (s *GitHubService) ListRepositories(ctx context.Context, username string) ([]*github.Repository, error) {
	var allRepos []*github.Repository
	opt := &github.RepositoryListOptions{
		Affiliation: "owner", // owner, collaborator, organization_member
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	for {
		// Pass empty string to get authenticated user's repos (including private)
		repos, resp, err := s.client.Repositories.List(ctx, "", opt)
		if err != nil {
			return nil, fmt.Errorf("failed to list repositories: %w", err)
		}

		allRepos = append(allRepos, repos...)

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allRepos, nil
}

// GetRepository gets repository information
func (s *GitHubService) GetRepository(ctx context.Context, owner, repo string) (*github.Repository, error) {
	repository, _, err := s.client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository: %w", err)
	}

	return repository, nil
}
