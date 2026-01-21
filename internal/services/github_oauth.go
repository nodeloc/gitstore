package services

import (
	"context"
	"fmt"

	"github.com/google/go-github/v57/github"
	"github.com/nodeloc/git-store/internal/config"
	"golang.org/x/oauth2"
	githubOAuth "golang.org/x/oauth2/github"
)

type GitHubOAuthService struct {
	config *config.Config
	oauth  *oauth2.Config
}

func NewGitHubOAuthService(cfg *config.Config) *GitHubOAuthService {
	oauthConfig := &oauth2.Config{
		ClientID:     cfg.GitHubClientID,
		ClientSecret: cfg.GitHubClientSecret,
		Endpoint:     githubOAuth.Endpoint,
		RedirectURL:  cfg.GitHubRedirectURL,
		Scopes:       []string{"user:email", "read:org"},
	}

	return &GitHubOAuthService{
		config: cfg,
		oauth:  oauthConfig,
	}
}

func (s *GitHubOAuthService) GetAuthURL(state string) string {
	return s.oauth.AuthCodeURL(state, oauth2.AccessTypeOnline)
}

func (s *GitHubOAuthService) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := s.oauth.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	return token, nil
}

func (s *GitHubOAuthService) GetUserInfo(ctx context.Context, accessToken string) (*github.User, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	return user, nil
}

func (s *GitHubOAuthService) GetUserEmails(ctx context.Context, accessToken string) ([]*github.UserEmail, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	emails, _, err := client.Users.ListEmails(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get user emails: %w", err)
	}

	return emails, nil
}

func (s *GitHubOAuthService) GetUserOrganizations(ctx context.Context, accessToken string) ([]*github.Organization, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	orgs, _, err := client.Organizations.List(ctx, "", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get user organizations: %w", err)
	}

	return orgs, nil
}
