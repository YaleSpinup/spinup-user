package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/YaleSpinup/spinup-user/pkg/user"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [username]",
	Short: "List existing users and their SSH keys",
	Long:  `Lists all existing Linux users on this host, or a specific user and their SSH authorized keys.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return fmt.Errorf("unexpected number of arguments: %d", len(args))
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			// get a list of all human users
			users, err := user.List()
			if err != nil {
				fmt.Printf("error listing users: %s\n", err)
				os.Exit(1)
			}

			for _, user := range users {
				fmt.Println(user)
			}
		} else {
			// get information about a specific user
			u, err := user.Get(args[0])
			if err != nil {
				fmt.Printf("\nfailed to get user: %s\n", err)
				os.Exit(1)
			}

			fmt.Printf("Username: %s\nHomedir: %s\nUID: %s\nGID: %s\n", u.Username, u.HomeDir, u.Uid, u.Gid)

			keys, err := user.AuthorizedKeys(args[0])
			if err != nil {
				fmt.Printf("\nfailed to get authorized_keys: %s\n", err)
				os.Exit(1)
			}

			fmt.Printf("Authorized keys:\n%s\n", strings.Join(keys, "\n"))
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
