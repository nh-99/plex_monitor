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

func InitDB(dataSourceName string) {
	var err error
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dataSourceName))
	if err != nil {
		logrus.Panic(err)
	}

	pingErr := client.Ping(context.TODO(), readpref.PrimaryPreferred())
	if pingErr != nil {
		logrus.Fatal(pingErr)
	}

	DB = client.Database("plex_monitor")
}
