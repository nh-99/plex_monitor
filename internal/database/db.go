package database

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// DB is the global database connection
var DB *mongo.Database

// Ctx is the global context
var Ctx = context.Background()

// InitDB initializes the database connection
func InitDB(dataSourceName string, dbName string) {
	var err error
	client, err := mongo.Connect(Ctx, options.Client().ApplyURI(dataSourceName))
	if err != nil {
		logrus.Panic(err)
	}

	pingErr := client.Ping(Ctx, readpref.PrimaryPreferred())
	if pingErr != nil {
		logrus.Fatal(pingErr)
	}

	DB = client.Database(dbName)

	setupIndexes()
}

// setupIndexes sets up the indexes for the database
func setupIndexes() {
	// Setup index on created_at
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "createdAt", Value: -1}},
	}
	_, err := DB.Collection(WebhookCollectionName).Indexes().CreateOne(Ctx, indexModel)
	if err != nil {
		logrus.Fatal(err)
	}

	// Setup index on serviceName
	indexModel = mongo.IndexModel{
		Keys: bson.D{{Key: "serviceName", Value: 1}},
	}
	_, err = DB.Collection(WebhookCollectionName).Indexes().CreateOne(Ctx, indexModel)
	if err != nil {
		logrus.Fatal(err)
	}
}

// CloseDB closes the database connection
func CloseDB() {
	DB.Client().Disconnect(Ctx)
}
