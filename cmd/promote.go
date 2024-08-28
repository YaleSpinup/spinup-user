package cmd

import (
	"fmt"
	"os"

	"github.com/YaleSpinup/spinup-user/pkg/linuxuser"
	"github.com/spf13/cobra"
)

var promoteCmd = &cobra.Command{
	Use:   "promote [username]",
	Short: "Promote an existing user to admin status",
	Long:  `Promote an existing Linux user to admin status by granting them sudo privileges.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("requires exactly one argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]

		// Check if the user exists
		_, err := linuxuser.Get(username)
		if err != nil {
			fmt.Printf("Error: User %s does not exist\n", username)
			os.Exit(1)
		}

		// Check if the user is already an admin
		isAdmin, err := linuxuser.HasSudo(username)
		if err != nil {
			fmt.Printf("Error checking admin status: %s\n", err)
			os.Exit(1)
		}
		if isAdmin {
			fmt.Printf("User %s is already an admin\n", username)
			os.Exit(0)
		}

		// Promote the user to admin
		err = linuxuser.UpdateSudo(username, true)
		if err != nil {
			fmt.Printf("Error promoting user: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully promoted %s to admin status\n", username)
	},
}

func init() {
	rootCmd.AddCommand(promoteCmd)
}
