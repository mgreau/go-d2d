package cfg

import (
	"os"
    "fmt"

	"github.com/spf13/viper"

	"github.com/mgreau/go-d2d/cmd"
)

const (
	APP_CONFIG_FILE      string = "go-d2d"
	APP_CONFIG_FILE_TYPE string = "properties"
	APP_CONFIG_FILE_PATH string = "/"
)

// LoadConfigFile load the needed configuration file
func LoadConfigFile() {
	// Init the conf app
	// Looking for ~/.go-d2d file
	viper.SetConfigName(APP_CONFIG_FILE)
	viper.SetConfigType(APP_CONFIG_FILE_TYPE)
	viper.AddConfigPath(os.Getenv("HOME") + APP_CONFIG_FILE_PATH)
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		cmd.PrintError("Fatal error while trying to load config file:", err)
	}
}

// DisplayUserConfiguration Show the loaded configuration 
func DisplayUserConfiguration() {
     fmt.Println("key | value: ")
    for _, k := range viper.AllKeys() {
        fmt.Println(""+k + " = "+ viper.GetString(k))
    }
}

// GetUserToken get the token from the config file
func GetUserToken() string {
	return viper.GetString("githubToken")
}
