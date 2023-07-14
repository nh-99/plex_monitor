package firehose

import (
	"context"
	"net/http"
	"plex_monitor/internal/database"

	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Firehose is the endpoint that streams data to the client
func Firehose(w http.ResponseWriter, r *http.Request) {
	opts := options.Find().SetLimit(1000) // Only return 1000 entries
	cursor, err := database.DB.Collection("raw_responses").Find(context.Background(), bson.D{}, opts)
	if err != nil {
		panic(err)
	}
	rawResponse := bson.M{"data": []bson.M{}}
	for cursor.Next(context.Background()) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			panic(err)
		}
		rawResponse["data"] = append(rawResponse["data"].([]bson.M), result)
	}
	render.JSON(w, r, rawResponse) // A chi router helper for serializing and returning json
}
