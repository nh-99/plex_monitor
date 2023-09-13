package cli

import (
	"fmt"
	"plex_monitor/internal/database/models"
	"time"

	"github.com/urfave/cli/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func getServiceCreateCmd() *cli.Command {
	return &cli.Command{
		Name:    "new",
		Aliases: []string{"n"},
		Usage:   "Configure a new key in the system. Returns the unique key for the service",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "name",
				Usage:    "The name of the service",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "key",
				Usage:    "The key for the service",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "host",
				Usage:    "The host for the service",
				Required: true,
			},
		},
		Action: func(cCtx *cli.Context) error {
			data := models.ServiceData{
				ServiceName: cCtx.String("name"),
				Config: bson.M{
					"key":  cCtx.String("key"),
					"host": cCtx.String("host"),
				},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			err := models.CreateService(data)
			if err != nil {
				return err
			}

			// Log a success message
			fmt.Println("Successfully created service")

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
