package main

import (
	"fmt"

	"github.com/cloudappforce/cloudctl/cmd"
	"github.com/cloudappforce/cloudctl/internal/constants"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {

	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	viper.SetEnvPrefix("CLOUDAPPFORCE")

	viper.BindEnv("ACCESS_TOKEN")
	viper.BindEnv("HOST")

	// TODO: move this to being dynamically populated via the REST API and JWT access token
	viper.BindEnv("REPO")
	viper.BindEnv("REPO_USERNAME")
	viper.BindEnv("REPO_PASSWORD")

	viper.BindEnv("SCHEME")
	viper.SetDefault("SCHEME", "https")

	viper.SetDefault("HOST", constants.CloudAppForceHost)

	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	if viper.GetString("ACCESS_TOKEN") == "" {
		err := fmt.Errorf("missing access token")
		fmt.Println(err)
		return
	}

	var rootCmd = &cobra.Command{
		Use:   "cloudctl",
		Short: "the cli interface for cloudappforce.com",
		Long:  `the cli interface for cloudappforce.com`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here

		},
	}
	rootCmd.AddCommand(cmd.Deploy())
	rootCmd.AddCommand(cmd.Database())
	rootCmd.Execute()

}
