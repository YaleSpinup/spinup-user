package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/YaleSpinup/spinup-user/pkg/user"
	"github.com/spf13/cobra"
)

var shell string

var addCmd = &cobra.Command{
	Use:   "add [username]",
	Short: "Add a new user and set SSH authorized keys",
	Long:  `Add a new Linux user and set their SSH authorized keys.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return fmt.Errorf("unexpected number of arguments: %d", len(args))
		}

		if len(args) == 0 {
			// prompt for username if not given
			fmt.Print("Enter username to create: ")
			fmt.Scanf("%s", &username)
		} else {
			username = args[0]
		}

		fmt.Println("Paste one or more SSH public keys for this user (hit Enter when done): ")
		scanner := bufio.NewScanner(os.Stdin)
		for {
			scanner.Scan()
			line := scanner.Text()
			if len(line) == 0 {
				break
			}

			if !user.ValidAuthorizedKey(line) {
				fmt.Printf("\ninvalid public key specified:\n%s\n", line)
				os.Exit(1)
			}

			pubKeys = append(pubKeys, line)
		}
		if err := scanner.Err(); err != nil {
			return err
		}

		if len(pubKeys) == 0 {
			fmt.Print("no public key specified\n")
			os.Exit(1)
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := user.Create(username, shell); err != nil {
			fmt.Printf("\nfailed to add new user: %s\n", err)
			os.Exit(1)
		}

		if err := user.UpdateAuthorizedKeys(username, pubKeys); err != nil {
			fmt.Printf("\nfailed to set authorized_keys: %s\n", err)
			// TODO: delete user
			os.Exit(1)
		}

		fmt.Printf("Added user %s\n", username)
	},
}

func init() {
	addCmd.PersistentFlags().StringVarP(&shell, "shell", "s", "/bin/bash", "login shell for the user")

	rootCmd.AddCommand(addCmd)
}
