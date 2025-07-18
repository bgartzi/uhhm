package subcmds

import (
	"fmt"
	"github.com/bgartzi/uhhm/lib/inventory"
	"github.com/bgartzi/uhhm/lib/local"
	"github.com/urfave/cli/v2"
	"strconv"
)

func sshHost(hostIdentifier *string) error {
	// Get inventory
	invFP, err := inventory.InventoryFilePath()
	if err != nil {
		return err
	}
	hosts := inventory.InitSimpleInventory(invFP)
	hostID, err := strconv.Atoi(*hostIdentifier)
	if err != nil {
		// Assume it could be a nickname
		hostID, err = hosts.HostIDFromNickname(hostIdentifier)
		if err != nil {
			return err
		}
	}
	// Get Host from inventory
	targetHost, err := hosts.GetHost(hostID)
	if err != nil {
		return err
	}
	// Get privateKeyAuthMethod
	privKeyPath, err := local.UHHMPrivateSSHKeyPath()
	if err != nil {
		return err
	}
	// Open ssh session
	return targetHost.Session(privKeyPath)
}

func Sesh() *cli.Command {
	return &cli.Command{
		Name:        "sesh",
		Aliases:     []string{"s"},
		Usage:       "Opens an ssh terminal session into a remote host",
		Description: "Receives a host identifier and opens a ssh session to it. Both indexes in the inventory and host nicknames are considered valid host identifiers.",
		Category:    "Remote SSH sessions",
		Flags:       []cli.Flag{},
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() != 1 {
				errMsg := fmt.Errorf("Expected 1 argument, got %d", ctx.NArg())
				return cli.NewExitError(errMsg, 1)
			}
			hostID := ctx.Args().Get(0)
			err := sshHost(&hostID)
			if err != nil {
				return cli.NewExitError(err, 2)
			}
			return nil
		},
	}
}
