package cli

import (
	"fmt"
	"time"

	"github.com/urfave/cli/v2"
)

// GetCliApp returns the CLI application for the Plex Monitor application.
func GetCliApp() *cli.App {
	return &cli.App{
		Name:     "Plex Monitor",
		Version:  "v1",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{Name: "Noah Howard"},
		},
		Copyright:   "(c) 2023 NoHowTech",
		HelpName:    "pm-cli",
		Usage:       "CLI interface for the Plex Monitor application.",
		UsageText:   "Can be used to configure and manage the Plex Monitor application.",
		HideHelp:    false,
		HideVersion: false,
		CommandNotFound: func(cCtx *cli.Context, command string) {
			fmt.Fprintf(cCtx.App.Writer, "Command %q not found.\n", command)
		},
		OnUsageError: func(cCtx *cli.Context, err error, isSubcommand bool) error {
			if isSubcommand {
				return err
			}

			fmt.Fprintf(cCtx.App.Writer, "WRONG: %#v\n", err)
			return nil
		},
		Action: func(cCtx *cli.Context) error {
			cli.ShowAppHelp(cCtx)

			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "create",
				Aliases: []string{"c"},
				Usage:   "Create a new object in the system",
				Subcommands: []*cli.Command{
					getUserCreateCmd(),
					{
						Name:    "service",
						Aliases: []string{"s"},
						Usage:   "Create a new object in the system",
						Subcommands: []*cli.Command{
							getServiceCreateCmd(),
							getServiceKeyRotateCmd(),
						},
					},
				},
			},
			{
				Name:    "fetch",
				Aliases: []string{"f"},
				Usage:   "Fetches an object from the system",
				Subcommands: []*cli.Command{
					getDumpWireFileCmd(),
				},
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "Lists objects from the system",
				Subcommands: []*cli.Command{
					getListFilesCmd(),
				},
			},
		},
	}
}
