package database

import (
	"context"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var DB *mongo.Database
var Ctx = context.Background()

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

func CloseDB() {
	DB.Client().Disconnect(Ctx)
}
