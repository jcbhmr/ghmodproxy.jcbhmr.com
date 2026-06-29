package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/google/go-github/v88/github"
)

func NewClient() (*github.Client, error) {
	tokenString := os.Getenv("GITHUB_TOKEN")
	if tokenString == "" {
		tokenString = os.Getenv("GH_TOKEN")
	}

	token, err := ParseGitHubToken(tokenString)
	if err != nil {
		return nil, err
	}

	return github.NewClient(github.WithAuthToken(string(token)))
}

type GitHubToken string

var githubClassicTokenRegexp = regexp.MustCompile(`^ghp_[a-zA-Z0-9]{36}$`)
var githubFineGrainedTokenRegexp = regexp.MustCompile(`^github_pat_[a-zA-Z0-9]{22}_[a-zA-Z0-9]{59}$`)
var githubActionsTokenRegexp = regexp.MustCompile(`^ghs_[a-zA-Z0-9]{36}$`)

func ParseGitHubToken(token string) (GitHubToken, error) {
	if !(githubClassicTokenRegexp.MatchString(token) || githubFineGrainedTokenRegexp.MatchString(token) || githubActionsTokenRegexp.MatchString(token)) {
		return GitHubToken(""), fmt.Errorf("%q does not match %s, %s, or %s: %w", token, githubClassicTokenRegexp, githubFineGrainedTokenRegexp, githubActionsTokenRegexp, strconv.ErrSyntax)
	}
	return GitHubToken(token), nil
}

type Repo struct {
	Owner string
	Repo  string
}

var ownerRegexp = regexp.MustCompile(`^[a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9-]{0,37}[a-zA-Z0-9]$`)
var repoRegexp = regexp.MustCompile(`^[a-zA-Z0-9._-]{1,100}$`)

func ParseRepo(owner string, repo string) (Repo, error) {
	if !ownerRegexp.MatchString(owner) {
		return Repo{}, fmt.Errorf("%q does not match %s: %w", owner, ownerRegexp, strconv.ErrSyntax)
	}
	if !repoRegexp.MatchString(repo) {
		return Repo{}, fmt.Errorf("%q does not match %s: %w", repo, repoRegexp, strconv.ErrSyntax)
	}
	return Repo{Owner: owner, Repo: repo}, nil
}
