package cli

import (
	"fmt"
	"plex_monitor/internal/database"

	"github.com/urfave/cli/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getFixCreatedAtTimesCommand() *cli.Command {
	return &cli.Command{
		Name:    "createdattime",
		Aliases: []string{"ct"},
		Usage:   "Add created_at times to records that don't have them",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "collection"},
		},
		Action: func(cCtx *cli.Context) error {
			db := database.DB
			collectionName := cCtx.String("collection")

			if collectionName == "" {
				fmt.Println("Please specify a collection name")
				return nil
			}

			fmt.Printf("Fixing %s\n", collectionName)
			c1, _ := db.Collection(collectionName).Find(database.Ctx, bson.M{"created_at": bson.M{"$exists": false}}, nil)
			for c1.Next(database.Ctx) {
				var result bson.M
				err := c1.Decode(&result)
				if err != nil {
					panic(err)
				}

				// Parse the _id to get the timestamp
				fmt.Printf("Record id as string %s\n", result["_id"].(primitive.ObjectID).Hex())
				objectID, _ := primitive.ObjectIDFromHex(result["_id"].(primitive.ObjectID).Hex())
				timestamp := objectID.Timestamp()

				// Set the created_at time
				update := bson.M{"$set": bson.M{"created_at": timestamp}}
				db.Collection(collectionName).UpdateOne(database.Ctx, bson.M{"_id": result["_id"]}, update)

				fmt.Printf("Fixed record with id %s\n", result["_id"])
			}
			return nil
		},
	}
}
