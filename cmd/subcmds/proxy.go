package subcmds

import (
	"fmt"
	"github.com/urfave/cli/v2"
)

func Proxy() *cli.Command {
	return &cli.Command{
		Name:        "proxy",
		Usage:       "????",
		UsageText:   "uhhm proxy",
		Description: "Do you really know what are you trying to achieve?",
		Category:    "Others",
		Action: func(ctx *cli.Context) error {
			fmt.Println("Really? A proxy... Could you explain to me what a proxy really is?")
			return nil
		},
	}
}
