package models

import (
	"plex_monitor/internal/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// ActivityStream is the data structure for the activity stream.
type ActivityStream struct {
	ID               string    `bson:"_id"`
	CreatedAt        time.Time `bson:"createdAt"`
	FormattedTimeAgo string    `bson:"formattedTimeAgo"`
	ServiceName      string    `bson:"serviceName"`
	Summary          string    `bson:"summary"`
}

// GetWebhookDataActivityCount returns the number of webhook data "activity" entries.
func GetWebhookDataActivityCount() (int64, error) {
	collection := database.DB.Collection(database.WebhookCollectionName)
	return collection.CountDocuments(database.Ctx, bson.M{"serviceName": "plex"})
}

// GetWebhookDataAsActivityStream returns the webhook data as an activity stream.
func GetWebhookDataAsActivityStream(offset int64, limit int64) ([]ActivityStream, error) {
	collection := database.DB.Collection(database.WebhookCollectionName)

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"serviceName": "plex",
			},
		},
		{
			"$addFields": bson.M{
				"timeDiff": bson.M{
					"$subtract": []interface{}{time.Now(), "$createdAt"},
				},
			},
		},
		{
			"$addFields": bson.M{
				"formattedTimeAgo": bson.M{
					"$cond": []interface{}{
						bson.M{"$lt": []interface{}{"$timeDiff", 60000}}, // Less than 1 minute
						bson.M{"$concat": []interface{}{bson.M{"$toString": bson.M{"$floor": bson.M{"$divide": []interface{}{"$timeDiff", 1000}}}}, "s ago"}},
						bson.M{
							"$cond": []interface{}{
								bson.M{"$lt": []interface{}{"$timeDiff", 3600000}}, // Less than 1 hour
								bson.M{"$concat": []interface{}{bson.M{"$toString": bson.M{"$floor": bson.M{"$divide": []interface{}{"$timeDiff", 60000}}}}, "m ago"}},
								bson.M{
									"$cond": []interface{}{
										bson.M{"$lt": []interface{}{"$timeDiff", 86400000}}, // Less than 1 day
										bson.M{"$concat": []interface{}{bson.M{"$toString": bson.M{"$floor": bson.M{"$divide": []interface{}{"$timeDiff", 3600000}}}}, "h ago"}},
										bson.M{
											"$cond": []interface{}{
												bson.M{"$lt": []interface{}{"$timeDiff", 604800000}}, // Less than 1 week
												bson.M{"$concat": []interface{}{bson.M{"$toString": bson.M{"$floor": bson.M{"$divide": []interface{}{"$timeDiff", 86400000}}}}, "d ago"}},
												bson.M{
													"$cond": []interface{}{
														bson.M{"$lt": []interface{}{"$timeDiff", 31536000000}}, // Less than 1 year
														bson.M{"$concat": []interface{}{bson.M{"$toString": bson.M{"$floor": bson.M{"$divide": []interface{}{"$timeDiff", 604800000}}}}, "w ago"}},
														bson.M{"$concat": []interface{}{bson.M{"$toString": bson.M{"$floor": bson.M{"$divide": []interface{}{"$timeDiff", 31536000000}}}}, "y ago"}},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			"$addFields": bson.M{
				"summary": bson.M{
					"$cond": []interface{}{
						bson.M{"$eq": []interface{}{"$serviceName", "plex"}},
						bson.M{"$concat": []interface{}{"$Account.title", " performed ", "$event", " on Plex with ", "$Player.title"}},
						"$serviceName",
					},
				},
			},
		},
		{
			"$sort": bson.M{
				"createdAt": -1,
			},
		},
		{
			"$skip": offset,
		},
		{
			"$limit": limit,
		},
		{
			"$project": bson.M{
				"_id":              1,
				"createdAt":        1,
				"formattedTimeAgo": 1,
				"serviceName":      1,
				"summary":          1,
			},
		},
	}

	cursor, err := collection.Aggregate(database.Ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(database.Ctx)

	var results []ActivityStream
	err = cursor.All(database.Ctx, &results)

	return results, err
}
