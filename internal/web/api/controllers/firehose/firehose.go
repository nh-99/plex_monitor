package firehose

import (
	"context"
	"net/http"
	"plex_monitor/internal/database"

	"github.com/go-chi/render"
	"go.mongodb.org/mongo-driver/bson"
)

// Firehose is the endpoint that streams data to the client
func Firehose(w http.ResponseWriter, r *http.Request) {
	cursor, err := database.DB.Collection("raw_responses").Find(context.TODO(), bson.D{})
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