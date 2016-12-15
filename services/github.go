package services

import (
	"fmt"
	github "github.com/google/go-github/github"

		"golang.org/x/oauth2"

	"net/http"
	"github.com/mitchellh/cli"
)

var (
	// github informations
	orgs   StringSlice
	dryrun bool
)

// CreateClient create the github/http client
func CreateClient(token string) *github.Client {

	// Create the http client.
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	// Create the github client.
	client := github.NewClient(tc)

	return client
}

// ListRepos list repositories
func ListRepos(token string) {

	// Create the github client.
	client := CreateClient(token)

	// Get the current user
	user, _, err := client.Users.Get("")
	if err != nil {
		fmt.Printf("go-d2d error")
	}
	username := *user.Login
	// add the current user to orgs
	fmt.Println(username)
	orgs = append(orgs, username)

	page := 1
	perPage := 20
	if err := GetRepositories(client, page, perPage); err != nil {
		fmt.Printf("go-d2d fatal")
	}

}

// CreateRepository create a repo
func CreateRepository(repoName *string, token string){

	// Create the github client.
	client := CreateClient(token)


	// Get the current user
	repo := github.Repository {
		Name: repoName,
	}
	rep, resp, err := client.Repositories.Create("mgreau", *repo)
	if err != nil {
		fmt.Printf("go-d2d error to create a repository")
	}

}


// GetRepositories repositories from GiHub
func GetRepositories(client *github.Client, page, perPage int) error {
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{
			Page:    page,
			PerPage: perPage,
		},
	}
	repos, resp, err := client.Repositories.List("", opt)
	if err != nil {
		return err
	}

	for _, repo := range repos {
		if err := handleRepo(client, repo); err != nil {
			fmt.Println("print warm message") 
		}
	}

	// Return early if we are on the last page.
	if page == resp.LastPage || resp.NextPage == 0 {
		return nil
	}

	page = resp.NextPage
	return GetRepositories(client, page, perPage)
}

// handleRepo will return nil error if the user does not have access to something.
func handleRepo(client *github.Client, repo *github.Repository) error {
	opt := &github.ListOptions{
		PerPage: 100,
	}

	branches, resp, err := client.Repositories.ListBranches(*repo.Owner.Login, *repo.Name, opt)
	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusForbidden {
		return nil
	}
	if err != nil {
		return err
	}

	for _, branch := range branches {
		if *branch.Name == "master" && in(orgs, *repo.Owner.Login) {
			// return early if it is already protected
			if *branch.Protected {
				fmt.Printf("[OK] %s:%s is already protected\n", *repo.FullName, *branch.Name)
				return nil
			}

			fmt.Printf("[UPDATE] %s:%s will be changed to protected\n", *repo.FullName, *branch.Name)
			if dryrun {
				// return early
				return nil
			}

		}
	}

	return nil
}

func in(a StringSlice, s string) bool {
	for _, b := range a {
		if b == s {
			return true
		}
	}
	return false
}


// StringSlice is a slice of strings
type StringSlice []string

// implement the flag interface for stringSlice
func (s *StringSlice) String() string {
	return fmt.Sprintf("%s", *s)
}
func (s *StringSlice) set(value string) error {
	*s = append(*s, value)
	return nil
}