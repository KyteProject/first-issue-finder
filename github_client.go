package main

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gofri/go-github-ratelimit/github_ratelimit"
	"github.com/google/go-github/v53/github"
)

func createGitHubClient() (*github.Client, error) {
	//TODO: check for token and use that if supplied
	// ts := oauth2.StaticTokenSource(
	// 	&oauth2.Token{AccessToken: token},
	// )
	// tc := oauth2.NewClient(ctx, ts)

	rateLimiter, err := github_ratelimit.NewRateLimitWaiterClient(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create rate limiter: %w", err)
	}

	client := github.NewClient(rateLimiter)
	return client, nil
}

func checkGitHubIssues() tea.Msg {
	client, err := createGitHubClient()
	if err != nil {
		return errMsg{err: fmt.Errorf("failed to create GitHub client: %w", err)}
	}

	issues, res, err := client.Search.Issues(
		context.Background(),
		"language:go is:issue is:public is:open no:assignee label:beginner,easy,first-timers-only,good-first-bug,good-first-issue,starter created:>2023-01-01",
		&github.SearchOptions{
			Sort:      "reactions",
			Order:     "dec",
			TextMatch: true,
		},
	)
	if err != nil {
		return errMsg{err: fmt.Errorf("failed to search issues: %w", err)}
	}
	if res.StatusCode != 200 {
		return errMsg{err: fmt.Errorf("unexpected status code: %d", res.StatusCode)}
	}

	return fetchIssuesMsg(issues)
}

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

type fetchIssuesMsg *github.IssuesSearchResult
