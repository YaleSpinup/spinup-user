package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var username string
var pubKeys []string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "spinup-user",
	Short: "Linux user and SSH key management utility",
	Long:  `A command line utility for easily managing Linux users and their SSH keys.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
