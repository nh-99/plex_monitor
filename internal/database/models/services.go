package models

import (
	"context"
	"plex_monitor/internal/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	// ServiceCollectionName is the name of the collection that stores the services
	ServiceCollectionName = "services"
)

// ServiceData is the struct that represents the service data that is stored in the database.
type ServiceData struct {
	ID          string    `bson:"_id,omitempty"`
	ServiceName string    `bson:"service_name"`
	Config      bson.M    `bson:"config"`
	CreatedAt   time.Time `bson:"created_at"`
	CreatedBy   string    `bson:"created_by,omitempty"`
	UpdatedAt   time.Time `bson:"updated_at"`
	UpdatedBy   string    `bson:"updated_by,omitempty"`
}

// GetAllServices returns all services.
func GetAllServices() ([]ServiceData, error) {
	var services []ServiceData

	cursor, err := database.DB.Collection(ServiceCollectionName).Find(database.Ctx, bson.M{"deleted_at": bson.M{"$exists": false}})
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.Background(), &services)
	if err != nil {
		return nil, err
	}

	return services, nil
}

// CreateService adds a new service to the database.
func CreateService(data ServiceData) (err error) {
	_, err = database.DB.Collection(ServiceCollectionName).InsertOne(database.Ctx, data)
	return err
}

// GetServiceByName returns a service by name.
func GetServiceByName(name string) (ServiceData, error) {
	var service ServiceData

	err := database.DB.Collection(ServiceCollectionName).FindOne(database.Ctx, bson.M{"service_name": name}).Decode(&service)
	if err != nil {
		return ServiceData{}, err
	}

	return service, nil
}
