package subcmds

import (
	"fmt"
	"github.com/bgartzi/uhhmm/inventory"
	"github.com/bgartzi/uhhmm/local"
	"github.com/urfave/cli/v2"
	"strconv"
)

func sshHost(hostID int) error {
	// Get inventory
	invFP, err := inventory.InventoryFilePath()
	if err != nil {
		return err
	}
	hosts := inventory.InitSimpleInventory(invFP)
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
		Description: "Receives a host ID (index in the inventory) and opens a ssh session to it.",
		Category:    "Remote SSH sessions",
		Flags:       []cli.Flag{},
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() != 1 {
				errMsg := fmt.Errorf("Expected 1 argument, got %d", ctx.NArg())
				return cli.NewExitError(errMsg, 1)
			}
			hostID, err := strconv.Atoi(ctx.Args().Get(0))
			if err != nil {
				errMsg := fmt.Errorf("Provided host ID is not a valid integer")
				return cli.NewExitError(errMsg, 1)
			}
			err = sshHost(hostID)
			if err != nil {
				return cli.NewExitError(err, 2)
			}
			return nil
		},
	}
}
