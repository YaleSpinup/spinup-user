package linuxuser

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
)

type Details struct {
	Admin    bool
	Username string
	Shell    string
	HomeDir  string
	Uid      string
	Gid      string
}

func (d Details) String() string {
	return fmt.Sprintf("Username: %s\nAdmin: %t\nShell: %s\nHomedir: %s\nUID: %s\nGID: %s\n", d.Username, d.Admin, d.Shell, d.HomeDir, d.Uid, d.Gid)
}

// Create uses the useradd executable to create a new Linux user
func Create(username, shell string) error {
	if username == "" || shell == "" {
		return errors.New("username and shell cannot be empty")
	}

	if !contains(shells(), shell) {
		return errors.New("invalid shell, valid options are: " + strings.Join(shells(), " "))
	}

	if _, err := exec.LookPath("useradd"); err != nil {
		return errors.New("didn't find 'useradd' executable: make sure it's installed and in the current PATH")
	}

	homeDir := "/home/" + username
	useraddArgs := []string{"--create-home", "--home-dir", homeDir, "--shell", shell, username}

	useraddCmd := exec.Command("useradd", useraddArgs...)
	useraddCmd.Stderr = os.Stderr

	// log.Printf("running command: %v", useraddCmd)

	if err := useraddCmd.Run(); err != nil {
		return err
	}

	return nil
}

// Delete uses the userdel executable to remove an existing Linux user
func Delete(username string, removeHomedir bool) error {
	if username == "" {
		return errors.New("username cannot be empty")
	}

	if _, err := exec.LookPath("userdel"); err != nil {
		return errors.New("didn't find 'userdel' executable: make sure it's installed and in the current PATH")
	}

	userdelArgs := []string{username}
	if removeHomedir {
		userdelArgs = append([]string{"-r"}, userdelArgs...)
	}

	userdelCmd := exec.Command("userdel", userdelArgs...)
	userdelCmd.Stderr = os.Stderr

	// log.Printf("running command: %v", userdelCmd)

	if err := userdelCmd.Run(); err != nil {
		return err
	}

	return nil
}

// Get returns details about the specified user
func Get(username string) (*Details, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	userMap, err := List()
	if err != nil {
		return nil, err
	}

	user, found := userMap[username]
	if !found {
		return nil, errors.New("unable to find user " + username + " in list of regular users")
	}

	return &user, nil
}

// List gets a map of all "human" Linux users details from /etc/passwd
func List() (map[string]Details, error) {
	userMap := make(map[string]Details)

	file, err := os.Open("/etc/passwd")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for {
		scanner.Scan()
		line := scanner.Text()
		if len(line) == 0 {
			break
		}
		if strings.HasPrefix(line, "#") {
			continue
		}

		flds := strings.Split(line, ":")
		if len(flds) < 5 {
			continue
		}

		// only pick users that have a valid shell in /etc/passwd
		// e.g. tester:x:1000:1000:Full Name:/home/tester:/bin/bash
		if contains(shells(), flds[len(flds)-1]) && flds[0] != "root" {
			admin, _ := HasSudo(flds[0])
			userMap[flds[0]] = Details{
				Admin:    admin,
				Username: flds[0],
				Shell:    flds[6],
				HomeDir:  flds[5],
				Uid:      flds[2],
				Gid:      flds[3],
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return userMap, nil
}

// AuthorizedKeys gets the SSH authorized keys for the given Linux user
func AuthorizedKeys(username string) ([]string, error) {
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	user, err := user.Lookup(username)
	if err != nil {
		return nil, errors.New("unable to get user information for " + username)
	}

	keysFile := fmt.Sprintf("%s/.ssh/authorized_keys", user.HomeDir)

	_, err = os.Stat(keysFile)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	keys, err := ioutil.ReadFile(keysFile)
	if err != nil {
		return nil, err
	}

	return strings.Split(strings.TrimSpace(string(keys)), "\n"), nil
}

// UpdateAuthorizedKeys updates the SSH authorized keys for the given Linux user
func UpdateAuthorizedKeys(username string, keys []string) error {
	if username == "" {
		return errors.New("username cannot be empty")
	}

	user, err := user.Lookup(username)
	if err != nil {
		return errors.New("unable to get user information for " + username)
	}

	uid, err := strconv.Atoi(user.Uid)
	if err != nil {
		return errors.New(fmt.Sprintf("couldn't get uid: %s", err))
	}
	gid, err := strconv.Atoi(user.Gid)
	if err != nil {
		return errors.New(fmt.Sprintf("couldn't get gid: %s", err))
	}

	sshDir := fmt.Sprintf("%s/.ssh", user.HomeDir)

	// create sshDir if it doesn't exist and chown
	if err := os.Mkdir(sshDir, 0700); err != nil && !os.IsExist(err) {
		return errors.New(fmt.Sprintf("couldn't create %s: %s", sshDir, err))
	}
	if err := os.Chown(sshDir, uid, gid); err != nil {
		return errors.New(fmt.Sprintf("couldn't chown %s: %s", sshDir, err))
	}

	// write public keys to authorized_keys and chown
	if err := os.WriteFile(sshDir+"/authorized_keys", []byte(strings.Join(keys, "\n")), 0600); err != nil {
		return errors.New(fmt.Sprintf("couldn't write to authorized_keys: %s", err))
	}
	if err := os.Chown(sshDir+"/authorized_keys", uid, gid); err != nil {
		return errors.New(fmt.Sprintf("couldn't chown authorized_keys: %s", err))
	}

	return nil
}

// ValidAuthorizedKey returns true if the specified key is a valid SSH public key
func ValidAuthorizedKey(key string) bool {
	if _, _, _, _, err := ssh.ParseAuthorizedKey([]byte(key)); err != nil {
		return false
	}

	return true
}

// contains returns true if the list contains the given string
func contains(list []string, s string) bool {
	for _, l := range list {
		if l == s {
			return true
		}
	}

	return false
}

// shells returns a list of valid shells
func shells() []string {
	shells := []string{"/bin/sh", "/bin/bash", "/bin/csh"}

	if shellsBytes, err := ioutil.ReadFile("/etc/shells"); err == nil {
		shells = strings.Split(strings.TrimSpace(string(shellsBytes)), "\n")
	}

	return shells
}
