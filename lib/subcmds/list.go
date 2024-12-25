package subcmds

import (
	"fmt"
	"github.com/bgartzi/uhhm/lib/display"
	"github.com/bgartzi/uhhm/lib/inventory"
	"github.com/urfave/cli/v2"
)

func listHosts() error {
	// Get inventory
	invFP, err := inventory.InventoryFilePath()
	if err != nil {
		return err
	}
	hostInventory := inventory.InitSimpleInventory(invFP)
	hosts := hostInventory.ListHosts()
	displayer := display.DefaultHostDisplayerConfig()
	err = displayer.Display(hosts)
	return err
}

func List() *cli.Command {
	return &cli.Command{
		Name:        "list",
		Aliases:     []string{"ls"},
		Usage:       "List host inventory data in stdout",
		UsageText:   "uhhm list [command options]",
		Description: "Lists hosts added to the inventory previously. Data is written onto stdout. Data is formatted as a table.",
		Category:    "Host inventory management",
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() != 0 {
				errMsg := fmt.Errorf("list doesn't accept arguments, got %d", ctx.NArg())
				return cli.NewExitError(errMsg, 1)
			}
			err := listHosts()
			if err != nil {
				return cli.NewExitError(err, 2)
			}
			return nil
		},
	}
}
