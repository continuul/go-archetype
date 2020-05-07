package cmd

import (
	"fmt"
	"github.com/continuul/go-archetype/cmd/generate"
	"github.com/continuul/go-archetype/cmd/list"
	"github.com/continuul/go-archetype/cmd/version"
	"github.com/spf13/cobra"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-archetype",
	Short: "A project templating toolkit",
	Long: `
Archetype is a project templating toolkit. An archetype is defined as an
original pattern or model from which all other things of the same kind are made.
Archetype is a system that provides a consistent means of generating projects of
virtually any kind.

Using archetypes provides a great way to enable developers quickly in a way
consistent with best practices employed by your project or organization.  We use
archetypes to try and get our users up and running as quickly as possible by providing a
sample project while introducing new users to current best practices. In a matter of
seconds, a new user can have a working project to use as a jumping board for their
work.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(
		version.Cmd,
		generate.Cmd,
		list.Cmd,
	)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-archetype.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".go-archetype" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".go-archetype")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
