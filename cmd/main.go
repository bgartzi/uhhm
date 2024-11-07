package main

import (
	"github.com/bgartzi/uhhmm/cmd/subcmds"
	"github.com/urfave/cli/v2"
	"os"
)

func Authors() []*cli.Author {
	return []*cli.Author{
		{
			Name:  "Beñat Gartzia",
			Email: "bgartzia@redhat.com",
		},
	}
}

func getSubCommands() []*cli.Command {
	return []*cli.Command{
		subcmds.Add(),
		subcmds.Delete(),
		subcmds.Sesh(),
		subcmds.Proxy(),
		subcmds.List(),
	}
}

func AppInit() *cli.App {
	return &cli.App{
		Name:                 "uhhm",
		HelpName:             "uhhm",
		Authors:              Authors(),
		Usage:                "Ultra Humble Host Manager",
		UsageText:            "Dummy way of managing your humble (almost inexistent) host inventory.",
		Commands:             getSubCommands(),
		EnableBashCompletion: true,
	}
}

func main() {
	AppInit().Run(os.Args)
}
