package services

import (
	"fmt"
	"github.com/google/go-github/github"

	"net/http"
)

var (
	// github informations
	token  string
	orgs   stringSlice
	dryrun bool

)

// stringSlice is a slice of strings
type stringSlice []string

// implement the flag interface for stringSlice
func (s *stringSlice) String() string {
	return fmt.Sprintf("%s", *s)
}
func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

// GetRepositories repositories from GiHub
func GetRepositories(client *github.client, page, perPage int) error {
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

func in(a stringSlice, s string) bool {
	for _, b := range a {
		if b == s {
			return true
		}
	}
	return false
}