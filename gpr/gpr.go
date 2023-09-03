package gpr

import (
	"context"
	"fmt"

	"github.com/gcraciun/simple-prboard/config"

	"github.com/google/go-github/v53/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type GithubClient struct {
	RestClient  *github.Client
	GraphClient *githubv4.Client
}

type RepoInfo struct {
	Name      string
	URL       string
	SshUrl    string
	IsPrivate bool
	OpenPRs   int
	ClosedPRs int
}

type RepoList map[string][]RepoInfo

// CreateGithubClient creates the Github clients REST and GraphQL used by other functions afterward
func CreateGithubClient(token string) *GithubClient {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &GithubClient{
		RestClient:  github.NewClient(tc),
		GraphClient: githubv4.NewClient(tc),
	}

}

// GetReposData creates the variable being passed to html/template ExecuteTemplate
// It calls GetPullRequestInfo for each repo and combines the data from all repos
func (client *GithubClient) GetReposData(ctx context.Context, owner string, config config.Config) (*RepoList, error) {
	repoList := make(RepoList)

	for category, names := range config.Repos {
		var repos []RepoInfo
		for _, name := range names {
			currentRepoInfo, err := client.GetPullRequestInfo(ctx, owner, name)
			if err != nil {
				return nil, err
			}
			repos = append(repos, *currentRepoInfo)
		}
		repoList[category] = repos
	}
	return &repoList, nil
}

// GetPullRequestInfo method on *GithubClient requests the PR information of a repo using the GraphQL client
func (client *GithubClient) GetPullRequestInfo(ctx context.Context, owner, repo string) (*RepoInfo, error) {
	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"repo":  githubv4.String(repo),
	}

	var query struct {
		Repository struct {
			Name      string
			Url       string
			SshUrl    string //githubv4.String
			IsPrivate bool   //githubv4.Boolean
			Open      struct {
				TotalCount int
			} `graphql:"Open: pullRequests(states: [OPEN])"`
			Closed struct {
				TotalCount int
			} `graphql:"Closed: pullRequests(states: [CLOSED, MERGED])"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}

	if err := client.GraphClient.Query(ctx, &query, variables); err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}

	return &RepoInfo{
		Name:      query.Repository.Name,
		URL:       query.Repository.Url,
		SshUrl:    query.Repository.SshUrl,
		IsPrivate: query.Repository.IsPrivate,
		OpenPRs:   query.Repository.Open.TotalCount,
		ClosedPRs: query.Repository.Closed.TotalCount,
	}, nil
}
