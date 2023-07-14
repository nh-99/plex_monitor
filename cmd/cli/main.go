package main

import (
	"fmt"
	"os"
	"time"

	pmcli "plex_monitor/internal/cli"
	"plex_monitor/internal/database"
	"plex_monitor/internal/database/models"
	"plex_monitor/internal/utils"

	"github.com/urfave/cli/v2"
)

func main() {
	database.InitDB(os.Getenv("DATABASE_URL"), os.Getenv("DATABASE_NAME"))

	app := &cli.App{
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
					{
						Name:    "user",
						Aliases: []string{"u"},
						Usage:   "Configure a new user in the system",
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "email"},
						},
						Action: func(cCtx *cli.Context) error {
							email := cCtx.String("email")
							password := pmcli.GetPassword("Enter a password: ")
							hashBytes, _ := utils.HashString(password)
							s := string(hashBytes)
							_, err := database.DB.Collection("users").InsertOne(database.Ctx, models.User{
								Email:          email,
								HashedPassword: s,
								Activated:      true,
								CreatedAt:      time.Now(),
								UpdatedAt:      time.Now(),
							})
							if err != nil {
								panic(err)
							}
							return nil
						},
					},
					{
						Name:    "service",
						Aliases: []string{"s"},
						Usage:   "Configure a new service in the system",
						Action: func(cCtx *cli.Context) error {
							//
							return nil
						},
					},
				},
			},
		},
	}

	app.Run(os.Args)
}
