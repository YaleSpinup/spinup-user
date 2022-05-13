# spinup-user

A simple CLI for managing Linux users

## Usage

```
$ spinup-user help
A command line utility for easily managing Linux users and their SSH keys

Usage:
  spinup-user [command]

Available Commands:
  add         Add a new user and set SSH authorized keys
  help        Help about any command
  list        List existing users and their SSH keys
  remove      Remove an existing user
  version     Show the current version

Flags:
  -h, --help   help for spinup-user

Use "spinup-user [command] --help" for more information about a command.
```

### Adding a user

Note that password authentication is not supported and you have to specify at least one public SSH key (for authorized_keys).
By default the `/bin/bash` shell is used but you can overide it with `--shell`

```
$ sudo spinup-user add alice
Paste one or more SSH public keys for this user (hit Enter when done):
ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIHSp/eBwht3KW6Kf6TQ+GTmubWYiaFfxf0BIKYq+4mDO

Added user alice
```

### Listing users

Will list all "human" users on the system (excluding root)

```
$ sudo spinup-user list
alice
bob
```

To get details about a specific user

```
$ sudo spinup-user list alice
Username: alice
Homedir: /home/alice
UID: 1001
GID: 1001
Authorized keys:
ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIHSp/eBwht3KW6Kf6TQ+GTmubWYiaFfxf0BIKYq+4mDO
```

### Removing a user

This will also remove the user's home directory, but you can add `-k` if you want to keep it

```
$ sudo spinup-user remove alice
Removed user alice
```

## Author

Tenyo Grozev <tenyo.grozev@yale.edu>

## License

GNU Affero General Public License v3.0 (GNU AGPLv3)
Copyright (c) 2021 Yale University
