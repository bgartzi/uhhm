# UHHM: Ultra Humble Host Manager

Simple tool to manage little pools of remote hosts. It helps manage, tag, and use remote hosts by storing their addresses, users, ssh ports and some more metadata into a simple local file. Once hosts are registered into the inventory (the local file), the tool can be used to ssh into any of the hosts without needing to input the host address nor the login password, or carry out other sort of related tasks over them.

## Why UHHM?
I used to write down the addresses to the hosts I owned at the time into a plain text file. Every single time I wanted to open a ssh session to a host I had to copy the address and paste it into a larger ssh command. Even, I had to write down the access password to such host.

I was really tired of doing that.

## Why is it written in GO?
Honestly, the project was originally implemented by a bunch of bash scripts. I wanted to learn a bit about GO, so rewriting the project seemed to be a good idea. No actual tehcnical reasons behind the decision.

If you prefer the original bash scripts (lol) ask me. I can hand you a little tar file including them. Note that my bash skills are below my golang skills so it is up to you to take care of them.

# Installation
To install the tool, clone, build and install it. After cloning,
```{bash}
$ make
$ make install
$ make clean  # If you wish
```
Notes on the installation:
- It will install the binary into `/usr/local`. If you would like to alter the installation path, set the `PREFIX` environment variable accordingly during `make install`.
- It will copy a file into `/etc/bash_completion.d/uhhm` by default for bash autocompletion support. `sudo` is needed for that.

## Uninstall
The tool can be also uninstalled by a make recipe:
```{bash}
$ make uninstall
```
It will remove the binary and the autocompletion dependencies for you.

# Usage
The whole idea is based on the idea of a "host inventory". It is a struct in memory that holds the list of registered hosts and their metadata. This data is serialized into a file that is stored in the local computer. For more information about the host inventory file, please check below in the [Working Principles](#working-principles) section.

## Add hosts
The first step is to add a host into the inventory. `uhhm add` can be used for that.
```{bash}
$ uhhm add foo-bar-001.my-hw-host-pool.com foobar001
```
That's the required data to add a host: the address (it can be an IP address too) and its nickname. Other configuration can be also added, as the user for login, or the default ssh port of the host, for example. Hosts also can carry metadata, such as "information" and "labels". For more information on those options, just `uhhm add --help`.

The command will probably ask for the user's login password to copy a public ssh key to the target host, so password won't be needed for login ever again.
The command will output the global inventory position of the newly added host. For example, if it was the third host you added to your inventory, it will return `2`.

## Open a ssh session
Now we have some hosts in our inventory. Let's open a ssh session to one of them. The `uhhm sesh` command can be used for that:
```{bash}
$ uhhm sesh 0
```
It will open a ssh session into the first host of the inventory (global position 0). We should also allow users to use nicknames instead, but we're still working to support that feature, sorry :(.

## List the hosts on your inventory
The hosts added in the local inventory can be checked by listing them:
```{bash}
$ uhhm ls
```
It will write a table with registered hosts and their metadata into the stdout.

## Delete a host
You're done working with a host. It's time to remove it from the inventory list. Just:
```{bash}
$ uhhm del <global position>
```
Where global position is the integer/index of the host you'd like to remove from your local inventory.


# Working principles

## Directories and files
### UHHM_HOME
`UHHM_HOME` is the path to the directory in which uhhm-related files will be stored into. To ease the workflow related to uhhm, it's recommended `UHHM_HOME` is set into a user-accessible location.

By default, `UHHM_HOME` will be `~/.uhhm/`, unless `UHHM_HOME` is an environment variable with a value.

Note that if modified, `UHHM_HOME` must exist whenever a `uhhm` command is used. That's why it should be added to the `.bashrc` file if it is desired to modify the default `UHHM_HOME` directory path.

### Host inventory file
The host inventory file holds the golang inventory structure serialized into a file. It is stored in `$UHHM_HOME/.HOST_INVENTORY`.

### ssh keys
uhhm uses it's own ssh keys. Those can be found under `$UHHM_HOME`. They will be created whenever a running `add` command can't find the keys where they're expected to be. Right now only RSA keys are supported. Contributions are welcome though :)

### known hosts
uhhm also adds and removes entries from the known_hosts file, usually stored in `.ssh/known_hosts`.
