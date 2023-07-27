package cli

import (
	"plex_monitor/internal/database"
	"plex_monitor/internal/database/models"
	"plex_monitor/internal/utils"
	"time"

	"github.com/urfave/cli/v2"
)

func getUserCreateCmd() *cli.Command {
	return &cli.Command{
		Name:    "user",
		Aliases: []string{"u"},
		Usage:   "Configure a new user in the system",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "email"},
		},
		Action: func(cCtx *cli.Context) error {
			email := cCtx.String("email")
			password := getPassword("Enter a password: ")
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
	}
}
