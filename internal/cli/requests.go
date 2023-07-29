package cli

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"plex_monitor/internal/database/models"
	"strings"

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
			&cli.StringFlag{Name: "filter", Required: false, Usage: "Partial matches of filename"},
		},
		Action: func(cCtx *cli.Context) error {
			filename := cCtx.String("filename")
			filter := cCtx.String("filter")

			if filename == "" && filter == "" {
				return cli.Exit("Filename or filter must be set", 1)
			}

			filenamesToParse := []string{}
			os.Mkdir("./output", 0755) // Create the output directory if it doesn't exist

			if filename != "" {
				filenamesToParse = append(filenamesToParse, filename)
			}

			if filter != "" {
				// Get the file from the database
				files, err := models.ListFilesInBucket(models.RawRequestWiresBucket, bson.M{"filename": bson.M{"$regex": filter}})
				if err != nil {
					return cli.Exit(err, 1)
				}

				for _, file := range files {
					fmt.Println("Found file:", file["filename"])
					filenamesToParse = append(filenamesToParse, file["filename"].(string))
				}
			}

			for _, fName := range filenamesToParse {
				fmt.Println("Writing file:", fName)
				// Get the file from the database
				fileBytes, err := models.GetFileFromBucket(models.RawRequestWiresBucket, fName)
				if err != nil {
					return cli.Exit(err, 1)
				}

				// Write the file to the filesystem
				f, err := os.Create("./output/" + fName)
				if err != nil {
					return err
				}
				defer f.Close()

				// Write the file to the filesystem
				_, err = f.Write(fileBytes)
				if err != nil {
					return err
				}
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

func getReplayWireFileCmd() *cli.Command {
	readRequest := func(raw, scheme string) (*http.Request, error) {
		r, err := http.ReadRequest(bufio.NewReader(strings.NewReader(raw)))
		if err != nil {
			return nil, err
		}
		r.RequestURI, r.URL.Scheme, r.URL.Host = "", scheme, r.Host
		return r, nil
	}
	return &cli.Command{
		Name:    "replay",
		Aliases: []string{"rp"},
		Usage:   "Replays a wire file against the local environment",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "filename"},
			&cli.StringFlag{Name: "filter", Required: false, Usage: "Partial matches of filename"},
		},
		Action: func(cCtx *cli.Context) error {
			filename := cCtx.String("filename")
			filter := cCtx.String("filter")

			if filename == "" && filter == "" {
				return cli.Exit("Filename or filter must be set", 1)
			}

			if filename != "" && filter != "" {
				return cli.Exit("Filename and filter cannot both be set", 1)
			}

			files := []string{}

			if filename != "" {
				files = append(files, filename)
			}

			if filter != "" {
				// Walk the filesystem and get all the files that match the filter
				err := filepath.Walk("./output", func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}

					if strings.Contains(path, filter) {
						files = append(files, path)
					}

					return nil
				})

				if err != nil {
					return cli.Exit(err, 1)
				}
			}

			for _, fName := range files {
				fmt.Println("Replaying file:", fName)
				// Read the file from the filesystem and put the string into a variable
				raw, err := os.ReadFile("./output/" + fName)
				if err != nil {
					return cli.Exit(err, 1)
				}
				rawStr := string(raw)

				fmt.Println("  Replaying request with data:\n", rawStr)

				// Read the request
				req, err := readRequest(rawStr, "http")
				if err != nil {
					return cli.Exit(err, 1)
				}

				// Execute the request
				new(http.Client).Do(req)
			}

			return nil
		},
	}
}
