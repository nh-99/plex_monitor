package models

import (
	"context"
	"fmt"
	"plex_monitor/internal/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	// ServiceCollectionName is the name of the collection that stores the services
	ServiceCollectionName = "services"
)

// ServiceType is the type of the service.
type ServiceType string

const (
	// ServiceTypePlex is the type of the Plex service.
	ServiceTypePlex ServiceType = "plex"
	// ServiceTypeOmbi is the type of the Ombi service.
	ServiceTypeOmbi ServiceType = "ombi"
)

// ServiceData is the struct that represents the service data that is stored in the database.
type ServiceData struct {
	ID          string      `bson:"_id,omitempty"`
	ServiceName ServiceType `bson:"service_name"`
	Config      bson.M      `bson:"config"`
	CreatedAt   time.Time   `bson:"created_at"`
	CreatedBy   string      `bson:"created_by,omitempty"`
	UpdatedAt   time.Time   `bson:"updated_at"`
	UpdatedBy   string      `bson:"updated_by,omitempty"`
}

// PlexConfig is the struct that represents the Plex config.
type PlexConfig struct {
	Host      string        `bson:"host"`
	Key       string        `bson:"key"`
	Libraries []PlexLibrary `bson:"libraries"`
}

// PlexLibrary is the struct that represents the Plex library.
type PlexLibrary struct {
	ID   int    `bson:"id"`
	Name string `bson:"name"`
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
func GetServiceByName(name ServiceType) (ServiceData, error) {
	var service ServiceData

	err := database.DB.Collection(ServiceCollectionName).FindOne(database.Ctx, bson.M{"service_name": name}).Decode(&service)
	if err != nil {
		return ServiceData{}, fmt.Errorf("failed to get service by name: %w", err)
	}

	return service, nil
}

// GetConfigAsPlexConfig returns the config as a Plex config.
func (s ServiceData) GetConfigAsPlexConfig() (*PlexConfig, error) {
	var config PlexConfig

	// Convert the bson.M to a byte array
	b, err := bson.Marshal(s.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	// Convert the byte array to a PlexConfig
	err = bson.Unmarshal(b, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config into plex config: %w", err)
	}

	return &config, nil
}

// GetLibraryByID returns a library by ID.
func (p PlexConfig) GetLibraryByID(id int) (PlexLibrary, error) {
	for _, library := range p.Libraries {
		if library.ID == id {
			return library, nil
		}
	}

	return PlexLibrary{}, fmt.Errorf("library with id %d not found", id)
}
