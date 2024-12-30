package main

import (
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
		Add(),
		Delete(),
		Sesh(),
		Proxy(),
		List(),
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
