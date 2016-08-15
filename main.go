package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/colorstring"
	//"github.com/spf13/viper"

	d2dcmd "go-d2d/cmd"
	d2dconf "go-d2d/conf"
	//d2dtypes "go-d2d/types"
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

)

func init() {

	// parse flags
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
}

// Check Prerequisites to run this program
// 1- config file
func checkPrerequisites() {

	// First load the required config file
	d2dconf.LoadConfigFile()

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
