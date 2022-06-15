package user

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
)

const sudoersDir = "/etc/sudoers.d"

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
		return false, nil
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

	// use full sudo privileges without password
	sudoersCmd := fmt.Sprintf("%s ALL=(ALL:ALL) NOPASSWD: ALL\n", username)
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
