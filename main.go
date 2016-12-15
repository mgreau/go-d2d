package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/colorstring"

	"github.com/mgreau/go-d2d/cmd"
	"github.com/mgreau/go-d2d/cfg"
	"github.com/mgreau/go-d2d/services"
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
	command     string
	version bool
	// to display help commands
	help bool

	// default values
	defaultCommand = cmd.CmdInfo

	// github informations
	token  string
	orgs   services.StringSlice
	dryrun bool

)

func init() {

	// parse flags
	flag.StringVar(&token, "token", "", "GitHub API token")
	//flag.Var(&orgs, "orgs", "organizations to include")
	flag.BoolVar(&dryrun, "dry-run", false, "do not change branch settings just print the changes that would occur")

	flag.BoolVar(&version, "version", false, "print version and exit")
	flag.BoolVar(&version, "v", false, "print version and exit (shorthand)")
	flag.BoolVar(&help, "help", false, "print help and exit")
	flag.BoolVar(&help, "h", false, "print help and exit (shorthand)")
	flag.StringVar(&command, "command", defaultCommand, "Command to run in each container")
	flag.StringVar(&command, "c", defaultCommand, "Command to run in each container (shorthand)")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(colorstring.Color("[blue]"+BANNER), VERSION))
		flag.PrintDefaults()
	}

	flag.Parse()
	// Display help if no flag
	if flag.NArg() > 0 {
		fmt.Println("nb args:", flag.Args())
		cmd.UsageAndExit("Pass the command name.", 1)
	}

	// if flags are ok, then check for pre-requisites
	checkPrerequisites()

	services.CreateClient(token)
}

// Check Prerequisites to run this program
// 1- config file
func checkPrerequisites() {

	// First load the required config file
	cfg.LoadConfigFile()

	token = cfg.GetUserToken()
	if token == "" {
		cmd.UsageAndExit("GitHub token cannot be empty.", 1)
	}

}

func main() {

	if version {
		fmt.Printf("go-d2d version %s, build %s", VERSION, "eDDdsD2EDSD")
		return
	}

	if help {
		cmd.UsageAndExit("", 0)
		return
	}

	log.Printf(colorstring.Color("[yellow]" + "[DEBUG] Running command %s"), command)
	switch command {
    case cmd.CmdInfo:
        fmt.Println("TODO; default informations")    
	default:
        fmt.Println(colorstring.Color("[blue]"+BANNER))
		fmt.Printf(colorstring.Color("[red] The command %s is not available. Please try -help to know all commands.\n"), command)
	}
}

