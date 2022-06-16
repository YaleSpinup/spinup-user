package cmd

import (
	"fmt"
	"os"

	"github.com/YaleSpinup/spinup-user/pkg/linuxuser"
	"github.com/spf13/cobra"
)

var keepHomedir bool

var removeCmd = &cobra.Command{
	Use:   "remove [username]",
	Short: "Remove an existing user",
	Long:  `Remove an existing Linux user and their home directory.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return fmt.Errorf("unexpected number of arguments: %d", len(args))
		}

		if len(args) == 0 {
			// prompt for username if not given
			fmt.Print("Enter username to remove: ")
			fmt.Scanf("%s", &username)
		} else {
			username = args[0]
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := linuxuser.Delete(username, !keepHomedir); err != nil {
			fmt.Printf("\nfailed to remove user: %s\n", err)
			os.Exit(1)
		}

		if err := linuxuser.RemoveSudo(username); err != nil {
			fmt.Printf("\nfailed to clean up sudoers: %s\n", err)
		}

		fmt.Printf("Removed user %s\n", username)
	},
}

func init() {
	removeCmd.PersistentFlags().BoolVarP(&keepHomedir, "keep-homedir", "k", false, "keep the user home directory")

	rootCmd.AddCommand(removeCmd)
}
