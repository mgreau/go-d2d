package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/colorstring"

	"golang.org/x/oauth2"
	"github.com/google/go-github/github"

	d2dcmd "go-d2d/cmd"
	d2dconf "go-d2d/conf"
	"net/http"
)

const (
	// BANNER is what is printed for help/info output
	BANNER = ` 

                           
        Improve my Day to Day tasks with automate tools.
        Version: %s
        
`
	// VERSION is the binary version.
	VERSION = "v0.1.0"
)

var (
	cmd     string
	version bool
	// to display help commands
	help bool

	// default values
	defaultCommand = d2dcmd.CmdInfo

	// github informations
	token  string
	orgs   stringSlice
	dryrun bool

)

func init() {

	// parse flags
	flag.StringVar(&token, "token", "", "GitHub API token")
	flag.Var(&orgs, "orgs", "organizations to include")
	flag.BoolVar(&dryrun, "dry-run", false, "do not change branch settings just print the changes that would occur")

	flag.BoolVar(&version, "version", false, "print version and exit")
	flag.BoolVar(&version, "v", false, "print version and exit (shorthand)")
	flag.BoolVar(&help, "help", false, "print help and exit")
	flag.BoolVar(&help, "h", false, "print help and exit (shorthand)")
	flag.StringVar(&cmd, "command", defaultCommand, "Command to run in each container")
	flag.StringVar(&cmd, "c", defaultCommand, "Command to run in each container (shorthand)")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(colorstring.Color("[blue]"+BANNER), VERSION))
		flag.PrintDefaults()
	}

	flag.Parse()
	// Display help if no flag
	if flag.NArg() > 0 {
		fmt.Println("nb args:", flag.Args())
		d2dcmd.UsageAndExit("Pass the command name.", 1)
	}

	// if flags are ok, then check for pre-requisites
	checkPrerequisites()

	// Create the http client.
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	// Create the github client.
	client := github.NewClient(tc)

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
	if err := getRepositories(client, page, perPage); err != nil {
		fmt.Printf("go-d2d fatal")
	}

}

// Check Prerequisites to run this program
// 1- config file
func checkPrerequisites() {

	// First load the required config file
	d2dconf.LoadConfigFile()

	if token == "" {
		d2dcmd.UsageAndExit("GitHub token cannot be empty.", 1)
	}

}

func main() {

	if version {
		fmt.Printf("go-d2d version %s, build %s", VERSION, "eDDdsD2EDSD")
		return
	}

	if help {
		d2dcmd.UsageAndExit("", 0)
		return
	}

	log.Printf(colorstring.Color("[yellow]" + "[DEBUG] Running command %s"), cmd)
	switch cmd {
    case d2dcmd.CmdInfo:
        fmt.Println("TODO; default informations")    
	default:
        fmt.Println(colorstring.Color("[blue]"+BANNER))
		fmt.Printf(colorstring.Color("[red] The command %s is not available. Please try -help to know all commands.\n"), cmd)
	}
}

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


func getRepositories(client *github.Client, page, perPage int) error {
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
	return getRepositories(client, page, perPage)
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