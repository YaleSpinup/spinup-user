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
  promote     Promote an existing user to admin status
  remove      Remove an existing user
  version     Show the current version

Flags:
  -h, --help   help for spinup-user

Use "spinup-user [command] --help" for more information about a command.
```

### Adding a user

Note that password authentication is not supported and by default you have to specify at least one public SSH key (for authorized_keys). You can skip setting an SSH key with the `--no-ssh` flag.

By default the `/bin/bash` shell is used but you can overide it with `--shell`

```
$ sudo spinup-user add alice
Paste one or more SSH public keys for this user (hit Enter when done):
ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIHSp/eBwht3KW6Kf6TQ+GTmubWYiaFfxf0BIKYq+4mDO

Added user alice
```

To create an admin user with full sudo privileges, just use the `-a` flag.

```
$ sudo spinup-user add helm -a
Paste one or more SSH public keys for this user (hit Enter when done):
ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIHSp/eBwht3KW6Kf6TQ+GTmubWYiaFfxf0BIKYq+4mDO

Added admin user helm
```

### Listing users

List all "human" users on the system

```
$ sudo spinup-user list
alice
bob
helm (admin)
```

Get details about a specific user

```
$ sudo spinup-user list alice
Username: alice
Admin: false
Homedir: /home/alice
UID: 1001
GID: 1001

Authorized keys:
ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIHSp/eBwht3KW6Kf6TQ+GTmubWYiaFfxf0BIKYq+4mDO
```

### Removing a user

When removing a user we also remove the user's home directory, but you can add `-k` if you want to keep it

```
$ sudo spinup-user remove alice
Removed user alice
```

### Promoting a user to admin status

To promote an existing user to admin status (granting them sudo privileges):

```
$ sudo spinup-user promote bob
Successfully promoted bob to admin status
```

## Author

Tenyo Grozev <tenyo.grozev@yale.edu>

## License

GNU Affero General Public License v3.0 (GNU AGPLv3)
Copyright (c) 2022 Yale University

