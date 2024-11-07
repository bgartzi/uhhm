package subcmds

import (
	"fmt"
	"github.com/bgartzi/uhhmm/lib/host"
	"github.com/bgartzi/uhhmm/lib/inventory"
	"github.com/bgartzi/uhhmm/lib/local"
	"github.com/urfave/cli/v2"
	"sort"
	"strconv"
	"strings"
)

func retrieveLabels(rawLabels string) []string {
	isComma := func(c rune) bool {
		return c == ','
	}
	labels := strings.FieldsFunc(rawLabels, isComma)
	sort.Strings(labels)
	return labels
}

func addHost(addedHost host.Host) error {
	// Get inventory
	invFP, err := inventory.InventoryFilePath()
	if err != nil {
		return err
	}
	// Get known hosts file
	knownHosts, err := local.InitKnownHosts()
	if err != nil {
		return fmt.Errorf(
			"Couldn't initialize uhhm known_hosts file successfully: %s", err)
	}
	hosts := inventory.InitSimpleInventory(invFP)
	// Look for host in the inventory or add it temporarily
	id, err := hosts.AddHost(addedHost)
	if err != nil {
		return err
	}
	// Get Host fingerprint/public key
	fingerprint, err := addedHost.PublicKey()
	if err != nil {
		return fmt.Errorf("Couldn't get fingerprint for host '%s'", addedHost.Address)
	}
	// Update known hosts file
	err = knownHosts.AddHost(addedHost, fingerprint)
	if err != nil {
		return fmt.Errorf("Couldn't add host to known hosts file: %s", err)
	}
	// If a UHHM rsa key pair doesn't exist, create one now
	uhhmPrivateKeyPath, err := local.UHHMPrivateSSHKeyPath()
	err = local.SafeCreateRSAKeyPair(uhhmPrivateKeyPath)
	if err != nil {
		return fmt.Errorf("UHHM RSA key pair doesn't exist and it couldn't be created.")
	}
	// Transfer UHHM rsa public key to host
	err = local.CopyPublicKeyTo(addedHost, uhhmPrivateKeyPath)
	if err != nil {
		errMsg := fmt.Errorf(
			"Couldn't copy local UHHM public ssh key to the remote host: %s.\n",
			err)
		err = knownHosts.RemoveHost(addedHost)
		if err != nil {
			errMsg = fmt.Errorf(
				"%s While on it, host %s couldn't be removed from known hosts file",
				errMsg, addedHost.Address)
		}
		return errMsg
	}
	err = hosts.Write()
	if err != nil {
		return err
	}
	fmt.Printf("Host %s added. Global inventory position:\n%d\n", addedHost.Address, id)
	return nil
}

func Add() *cli.Command {
	var hostPort string
	var hostInfo string
	var hostUser string
	var rawLabels string
	return &cli.Command{
		Name:        "add",
		Aliases:     []string{"a"},
		Usage:       "Adds host to the inventory",
		UsageText:   "uhhm add [command options] <address> <nickname>",
		Description: "Gets host address, nickname (and optional extra information) and updates the inventory.\nIt also adds the remote host entry to uhhm's known_hosts file.\nFinally it transfers uhhm's public rsa key to the remote host.",
		Category:    "Host inventory management",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "port",
				Aliases:     []string{"p"},
				Usage:       "Remote host's ssh `PORT`.",
				Value:       "22",
				Destination: &hostPort,
				Action: func(ctx *cli.Context, s string) error {
					_, err := strconv.Atoi(s)
					if err != nil {
						return fmt.Errorf("Port %s is not an integer", s)
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:        "info",
				Aliases:     []string{"i"},
				Usage:       "Optional extended text `INFO` of the remote host.",
				Value:       "",
				Destination: &hostInfo,
			},
			&cli.StringFlag{
				Name:        "user",
				Aliases:     []string{"u"},
				Usage:       "Remote host `USER`.",
				Value:       "root",
				Destination: &hostUser,
			},
			&cli.StringFlag{
				Name:        "labels",
				Aliases:     []string{"l"},
				Usage:       "Comma-sepparated host labels: `LABEL1,LABEL2,...,LABELN`",
				Value:       "",
				Destination: &rawLabels,
			},
		},
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() != 2 {
				errMsg := fmt.Errorf("Expected 2 arguments, got %d", ctx.NArg())
				return cli.NewExitError(errMsg, 1)
			}
			err := addHost(
				host.Host{
					Address:  ctx.Args().Get(0),
					Port:     hostPort,
					NickName: ctx.Args().Get(1),
					Info:     hostInfo,
					User:     hostUser,
					Labels:   retrieveLabels(rawLabels),
				},
			)
			if err != nil {
				return cli.NewExitError(err, 2)
			}
			return nil
		},
	}
}
