package cli

import "github.com/urfave/cli/v2"

func getServiceCreateCmd() *cli.Command {
	return &cli.Command{
		Name:    "new",
		Aliases: []string{"n"},
		Usage:   "Configure a new key in the system. Returns the unique key for the service",
		Action: func(cCtx *cli.Context) error {
			//
			return nil
		},
	}
}

func getServiceKeyRotateCmd() *cli.Command {
	return &cli.Command{
		Name:    "rotate",
		Aliases: []string{"r"},
		Usage:   "Rotates a key for a service",
		Action: func(cCtx *cli.Context) error {
			//
			return nil
		},
	}
}
