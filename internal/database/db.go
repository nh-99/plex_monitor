package database

import (
	"context"

	"github.com/sirupsen/logrus"
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
}

// CloseDB closes the database connection
func CloseDB() {
	DB.Client().Disconnect(Ctx)
}
