package subcmds

import (
	"fmt"
	"github.com/bgartzi/uhhm/lib/inventory"
	"github.com/bgartzi/uhhm/lib/local"
	"github.com/urfave/cli/v2"
	"strconv"
)

func delHost(hostID int) error {
	// Get inventory file path
	invFP, err := inventory.InventoryFilePath()
	if err != nil {
		return err
	}
	hosts := inventory.InitSimpleInventory(invFP)
	// Remove host from inventory; It doesn't overwrite the inventory file yet
	removedHost, err := hosts.DelHost(hostID)
	if err != nil {
		return err
	}
	// Get known hosts file
	knownHosts, err := local.InitKnownHosts()
	if err != nil {
		return fmt.Errorf(
			"Couldn't initialize uhhm known_hosts file successfully: %s", err)
	}
	// Remove host from known hosts file
	err = knownHosts.RemoveHost(removedHost)
	if err != nil {
		return fmt.Errorf(
			"Error removing host %d(%s) from known hosts file: %s", hostID, removedHost.Address, err)
	}
	// Now that all operations were successful, update the inventory file
	err = hosts.Write()
	if err != nil {
		// FIXME: Should we try restoring the knownHosts file to its previous state?
		return err
	}
	return nil
}

func Delete() *cli.Command {
	return &cli.Command{
		Name:        "del",
		Aliases:     []string{"d"},
		Usage:       "Deletes a host from the inventory",
		Description: "Receives a host ID (index in the inventory) and removes it from there.\nIt also removes host references made in the local known_hosts file.\nNote that local uhhm's public rsa key is not removed from the remote host.",
		Category:    "Host inventory management",
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
			err = delHost(hostID)
			if err != nil {
				return cli.NewExitError(err, 2)
			}
			return nil
		},
	}

}
