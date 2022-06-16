package linuxuser

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

// We make "admin" users by giving them full sudo privileges to run any command with no password by default.
// We create a separate file for each user in the sudoersDir which has the default sudo rule. This way we can
// later simply check if a file exists for a user to know if they are admin. It also allows users to manually
// tweak the sudo rule in a file, if needed, without affecting the functionality.
// Obviously, this won't work seamlessly with any pre-existing sudo users but the different ways to set
// sudo permissions are too many so we choose to ignore those and only manage users created by this CLI.
// However, we implement an additional check for sudo privileges if a file doesn't exist for a user by calling
// sudo -l and checking the output for the default sudo rule (it has to be an exact match).

// location of sudoers includedir
const sudoersDir = "/etc/sudoers.d"

// by default we allow admin users to sudo as anyone and run all commands with no password
const sudoersPrivs = "(ALL) NOPASSWD: ALL"

// HasSudo returns true if the Linux user has sudo privileges
func HasSudo(username string) (bool, error) {
	if username == "" {
		return false, errors.New("username cannot be empty")
	}

	if _, err := exec.LookPath("sudo"); err != nil {
		return false, errors.New("didn't find 'sudo' executable: make sure it's installed and in the current PATH")
	}

	_, err := user.Lookup(username)
	if err != nil {
		return false, errors.New("unable to get user information for " + username)
	}

	sudoersFile := fmt.Sprintf("%s/%s", sudoersDir, username)

	// we just check that a file exists for that user in the sudoers dir
	sudoer, err := os.Stat(sudoersFile)
	if os.IsNotExist(err) {
		// but if the file doesn't exist we'll do a real check with the sudo command
		return sudoUserCheck(username, sudoersPrivs)
	} else if err != nil {
		return false, err
	}

	return !sudoer.IsDir(), nil
}

// RemoveSudo removes sudo privileges (if set) for a user without checking that user actually exists
func RemoveSudo(username string) error {
	if username == "" {
		return errors.New("username cannot be empty")
	}

	sudoersFile := fmt.Sprintf("%s/%s", sudoersDir, username)

	_, err := os.Stat(sudoersFile)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	// remove sudoers file
	if err := os.Remove(sudoersFile); err != nil {
		return errors.New(fmt.Sprintf("couldn't remove file %s: %s", sudoersFile, err))
	}

	return nil
}

// UpdateSudo enables or disables sudo privileges for the given Linux user
func UpdateSudo(username string, admin bool) error {
	if username == "" {
		return errors.New("username cannot be empty")
	}

	if _, err := exec.LookPath("sudo"); err != nil {
		return errors.New("didn't find 'sudo' executable: make sure it's installed and in the current PATH")
	}

	_, err := user.Lookup(username)
	if err != nil {
		return errors.New("unable to get user information for " + username)
	}

	// give the user sudo privileges on all hosts
	sudoersCmd := fmt.Sprintf("%s ALL=%s\n", username, sudoersPrivs)
	sudoersFile := fmt.Sprintf("%s/%s", sudoersDir, username)

	if admin {
		// write sudoers file and chown
		if err := os.WriteFile(sudoersFile, []byte(sudoersCmd), 0440); err != nil {
			return errors.New(fmt.Sprintf("couldn't write to %s: %s", sudoersFile, err))
		}
	} else {
		// remove sudoers file
		if err := os.Remove(sudoersFile); err != nil {
			return errors.New(fmt.Sprintf("couldn't remove file %s: %s", sudoersFile, err))
		}
	}

	return nil
}

// sudoUserCheck checks if the user has the specified sudo privileges by running "sudo -l"
func sudoUserCheck(username, privs string) (bool, error) {
	out, err := exec.Command("sudo", "-l", "-U", username).Output()
	if err != nil {
		return false, err
	}

	return strings.Contains(string(out), privs), nil
}
