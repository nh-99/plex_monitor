package cli

import (
	"fmt"
	"os"
	"plex_monitor/internal/database/models"

	"github.com/urfave/cli/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func getDumpWireFileCmd() *cli.Command {
	return &cli.Command{
		Name:    "request",
		Aliases: []string{"r"},
		Usage:   "Get a wire request from the database",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "filename"},
		},
		Action: func(cCtx *cli.Context) error {
			filename := cCtx.String("filename")
			if filename == "" {
				return cli.Exit("Filename cannot be empty", 1)
			}

			// Get the file from the database
			fileBytes, err := models.GetFileFromBucket(models.RawRequestWiresBucket, filename)
			if err != nil {
				return cli.Exit(err, 1)
			}

			// Write the file to the filesystem
			f, err := os.Create(filename)
			if err != nil {
				return err
			}
			defer f.Close()

			// Write the file to the filesystem
			_, err = f.Write(fileBytes)
			if err != nil {
				return err
			}

			return nil
		},
	}
}

func getListFilesCmd() *cli.Command {
	return &cli.Command{
		Name:    "requests",
		Aliases: []string{"r"},
		Usage:   "Lists the request filenames in the database",
		Action: func(cCtx *cli.Context) error {
			// Get the file from the database
			files, err := models.ListFilesInBucket(models.RawRequestWiresBucket, bson.M{})
			if err != nil {
				return cli.Exit(err, 1)
			}

			for _, file := range files {
				fmt.Printf("%s\n", file["filename"].(string))
			}

			return nil
		},
	}
}
