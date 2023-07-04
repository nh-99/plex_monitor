package models

import (
	"context"
	"plex_monitor/internal/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type ServiceData struct {
	ID          string    `bson:"_id"`
	ServiceName string    `bson:"service_name"`
	Config      bson.M    `bson:"config"`
	CreatedAt   time.Time `bson:"created_at"`
	CreatedBy   string    `bson:"created_by"`
	UpdatedAt   time.Time `bson:"updated_at"`
	UpdatedBy   string    `bson:"updated_by"`
}

func GetAllServices() ([]ServiceData, error) {
	var services []ServiceData

	cursor, err := database.DB.Collection("services").Find(context.Background(), bson.M{"deleted_at": bson.M{"$exists": false}})
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.Background(), &services)
	if err != nil {
		return nil, err
	}

	return services, nil
}
